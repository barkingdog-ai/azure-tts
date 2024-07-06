package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/barkingdog-ai/azure-tts/model"
)

// TTSApiXMLPayload templates the payload required for API.
// See: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#sample-request
const ttsApiXMLPayload = "<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'><prosody rate='%s' pitch='%s'>%s</prosody></voice></speak>"

func (az *AzureTTSClient) newTokenRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
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

func (az *AzureTTSClient) newTTSRequest(ctx context.Context, method, path string, payload io.Reader, audioOutput model.AudioOutput) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Microsoft-OutputFormat", fmt.Sprint(audioOutput))
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", az.AccessToken))
	req.Header.Set("User-Agent", "azuretts")

	return req, nil
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

func (az *AzureTTSClient) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
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

func (az AzureTTSClient) performReq(req2 *http.Request) (*http.Response, error) {
	// Create a new HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://eastasia.api.cognitive.microsoft.com/sts/v1.0/issueToken", bytes.NewBuffer([]byte("")))
	if err != nil {
		return nil, err
	}
	// Add required headers
	req.Header.Add("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	req.Header.Add("Content-Length", "0")
	req2 = req // copy --> TODO: fix this
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
		return fmt.Sprintf("%d - A required parameter is missing, empty, or null. Or, the value passed to either a required or optional parameter is invalid. A common issue is a header that is too long", resp.StatusCode)
	case http.StatusUnauthorized:
		return fmt.Sprintf("%d - The request is not authorized. Check to make sure your subscription key or token is valid and in the correct region", resp.StatusCode)
	case http.StatusRequestEntityTooLarge:
		return fmt.Sprintf("%d - The SSML input is longer than 1024 characters", resp.StatusCode)
	case http.StatusUnsupportedMediaType:
		return fmt.Sprintf("%d - It's possible that the wrong Content-Type was provided. Content-Type should be set to application/ssml+xml", resp.StatusCode)
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

func jsonBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}
	return bytes.NewBuffer(raw), nil
}

func getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json resp: %w", err)
	}
	return nil
}

// voiceXML renders the XML payload for the TTS api.
// For API reference see https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#sample-request
func voiceXML(speechText, description string, locale model.Locale, gender model.Gender, rate string, pitch string) string {
	return fmt.Sprintf(ttsApiXMLPayload, locale, locale, gender, description, rate, pitch, speechText)
}
