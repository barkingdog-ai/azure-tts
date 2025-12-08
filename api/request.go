package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/barkingdog-ai/azure-tts/model"
)

// TTSApiXMLPayload templates the payload required for API.
// See: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#sample-request
const ttsAPIXMLPayload = `<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'><prosody rate='%s' pitch='%s'>%s</prosody></voice></speak>`

func (az *AzureTTSClient) newTokenRequest(ctx context.Context, method, path string, payload any) (*http.Request, error) {
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "0")

	return req, nil
}

func (az *AzureTTSClient) newTTSRequest(ctx context.Context,
	method, path string, payload io.Reader, audioOutput model.AudioOutput) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Microsoft-OutputFormat", audioOutput.String())
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", az.AccessToken))
	req.Header.Set("User-Agent", "azuretts")

	return req, nil
}

func createFilePayload(request model.SpeechToTextReq) (io.Reader, error) {
	if request.Reader != nil {
		return request.Reader, nil
	}

	f, err := os.Open(request.FilePath)
	if err != nil {
		return nil, fmt.Errorf("opening audio file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return nil, fmt.Errorf("reading audio file: %w", err)
	}

	return &buf, nil
}

func (az *AzureTTSClient) newSTTRequest(ctx context.Context, method, path string,
	payload io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	req.Header.Set("Content-Type", "audio/wav; codecs=audio/pcm; samplerate=16000")
	return req, nil
}

func (az *AzureTTSClient) newRequest(ctx context.Context, method, path string, payload any) (*http.Request, error) {
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", az.AccessToken))
	req.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	return req, nil
}

func (az *AzureTTSClient) performReq(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx,
		"POST", "https://eastasia.api.cognitive.microsoft.com/sts/v1.0/issueToken",
		bytes.NewBuffer([]byte("")))
	if err != nil {
		return nil, err
	}
	// Add required headers
	req.Header.Add("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	req.Header.Add("Content-Length", "0")
	request = req
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkForSuccess(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (az *AzureTTSClient) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := az.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkForSuccess(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func errorMessage(resp *http.Response) string {
	// list of acceptable resp status codes
	// see: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#http-status-codes-1
	switch resp.StatusCode {
	case http.StatusBadRequest:
		const message = "- A required parameter is missing, empty, or null. " +
			"Or, the value passed to either a required or optional parameter is invalid. " +
			"A common issue is a header that is too long."

		return fmt.Sprintf("%d "+message, resp.StatusCode)
	case http.StatusUnauthorized:
		const message = "- The request is not authorized. Check to make sure your subscription key or token is valid and in the correct region"
		return fmt.Sprintf("%d"+message, resp.StatusCode)
	case http.StatusRequestEntityTooLarge:
		return fmt.Sprintf("%d - The SSML input is longer than 1024 characters", resp.StatusCode)
	case http.StatusUnsupportedMediaType:
		const message = "- It's possible that the wrong Content-Type was provided. Content-Type should be set to application/ssml+xml"
		return fmt.Sprintf("%d "+message, resp.StatusCode)
	case http.StatusTooManyRequests:
		return fmt.Sprintf("%d - You have exceeded the quota or rate of requests allowed for your subscription", resp.StatusCode)
	case http.StatusBadGateway:
		return fmt.Sprintf("%d - Network or server-side issue. May also indicate invalid headers", resp.StatusCode)
	}
	return ""
}

// returns an error if this resp includes an error.
func checkForSuccess(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	var result model.APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		message := errorMessage(resp)
		// if we can't decode the json error then create an unexpected error
		APIError := model.APIError{
			StatusCode: resp.StatusCode,
			Type:       "Unexpected",
			Message:    message + string(data),
		}

		return APIError
	}
	result.Error.StatusCode = resp.StatusCode
	return result.Error
}

func jsonBodyReader(body any) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}
	return bytes.NewBuffer(raw), nil
}

func getResponseObject(rsp *http.Response, v any) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json resp: %w", err)
	}
	return nil
}

// voiceXML renders the XML payload for the TTS api.
// For API reference see https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#sample-request
func voiceXML(speechText, description string, locale model.Locale, gender model.Gender, rate, pitch string, style *model.TTSStyle) string {
	processedText := speechText

	// 先处理IP地址
	reIP := regexp.MustCompile(`\b\d+\.\d+\.\d+\.\d+\b`)
	processedText = reIP.ReplaceAllStringFunc(processedText, func(match string) string {
		return fmt.Sprintf("<say-as interpret-as=\"characters\">%s</say-as>", match)
	})

	// 处理Markdown格式的URL，保留描述文本
	reMarkdownURL := regexp.MustCompile(`\[(.*?)\]\((https?:\/\/[^\s\)]+)\)`)
	processedText = reMarkdownURL.ReplaceAllString(processedText, "$1")

	// 处理普通URL，将URL转换为可读格式
	reURL := regexp.MustCompile(`https?:\/\/[^\s]+`)
	// 用于检测是否为纯IP形式的域名
	reIPDomain := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)

	processedText = reURL.ReplaceAllStringFunc(processedText, func(match string) string {
		// 检查这个URL是否是markdown链接的一部分
		markdownPattern := fmt.Sprintf(`\[.*?\]\(%s\)`, regexp.QuoteMeta(match))
		if regexp.MustCompile(markdownPattern).MatchString(speechText) {
			return match // 如果是markdown链接的一部分，保持原样
		}
		// 将URL转换为可读格式，只处理域名部分
		url := strings.TrimPrefix(match, "http://")
		url = strings.TrimPrefix(url, "https://")
		// 只取域名部分（到第一个/之前）
		if idx := strings.Index(url, "/"); idx != -1 {
			url = url[:idx]
		}

		// 如果域名是纯IP形式，则添加say-as标签
		if reIPDomain.MatchString(url) {
			return fmt.Sprintf("<say-as interpret-as=\"characters\">%s</say-as>", url)
		}
		// 如果域名包含字母，直接返回
		return url
	})

	// 处理日期和其他数字序列
	reDate := regexp.MustCompile(`\b(\d{4}|\d{2})\.\d{1,2}\.\d{1,2}\b`)
	reNumbers := regexp.MustCompile(`\b\d+\.\d+(\.\d+)*\b`)
	processedText = reNumbers.ReplaceAllStringFunc(processedText, func(match string) string {
		// 如果已经包含在say-as标签中，跳过处理
		if strings.Contains(match, "<say-as") {
			return match
		}

		// 检查是否是IP地址格式，如果是则跳过（因为已经处理过）
		if reIP.MatchString(match) {
			return match
		}

		// 检查是否是有效日期格式
		if reDate.MatchString(match) {
			parts := strings.Split(match, ".")
			if len(parts) == 3 {
				month, _ := strconv.Atoi(parts[1])
				day, _ := strconv.Atoi(parts[2])
				if month >= 1 && month <= 12 && day >= 1 && day <= 31 {
					return match
				}
			}
		}
		return fmt.Sprintf("<say-as interpret-as=\"characters\">%s</say-as>", match)
	})

	if style != nil {
		// 使用帶風格的 SSML 模板
		// 如果 rate 和 pitch 為 0%，則不包含 prosody 標籤
		if rate == "0%" && pitch == "0%" {
			styledPayload := `<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="https://www.w3.org/2001/mstts" xml:lang="%s"><voice name="%s"><mstts:express-as style="%s" styledegree="%s">%s</mstts:express-as></voice></speak>`
			return fmt.Sprintf(styledPayload, locale, description, style.Style, style.StyleDegree, processedText)
		} else {
			styledPayload := `<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="https://www.w3.org/2001/mstts" xml:lang="%s"><voice name="%s"><mstts:express-as style="%s" styledegree="%s"><prosody rate="%s" pitch="%s">%s</prosody></mstts:express-as></voice></speak>`
			return fmt.Sprintf(styledPayload, locale, description, style.Style, style.StyleDegree, rate, pitch, processedText)
		}
	}

	// 如果 rate 和 pitch 為 0%，則不包含 prosody 標籤
	if rate == "0%" && pitch == "0%" {
		basicPayload := `<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'>%s</voice></speak>`
		return fmt.Sprintf(basicPayload, locale, locale, gender, description, processedText)
	}

	return fmt.Sprintf(ttsAPIXMLPayload, locale, locale, gender, description, rate, pitch, processedText)
}
