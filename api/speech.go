package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/barkingdog-ai/azure-tts/model"
	"github.com/barkingdog-ai/azure-tts/utils"
)

type SpeechInterface interface {
	TextToSpeech(ctx context.Context, req model.TextToSpeechRequest) ([]byte, error)
	SpeechToText(ctx context.Context, req model.SpeechToTextReq) (*model.SpeechToTextResp, error)
}

func (az *AzureTTSClient) TextToSpeech(ctx context.Context,
	request model.TextToSpeechRequest,
) ([]byte, error) {
	respData := make([]byte, 0)
	rate, _ := utils.ConvertStringToFloat32(request.Rate)
	pitch, _ := utils.ConvertStringToFloat32(request.Pitch)
	v := voiceXML(
		request.SpeechText,
		request.VoiceName,
		request.Locale,
		request.Gender,
		utils.ConvertFloat32ToString((rate-1)*100)+"%",
		utils.ConvertFloat32ToString((pitch-1)*50)+"%")

	req, err := az.newTTSRequest(ctx, "POST", az.TextToSpeechURL, bytes.NewBufferString(v), model.Audio16khz32kbitrateMonoMp3)
	if err != nil {
		return respData, fmt.Errorf("tts request error %v", err)
	}

	resp, err := az.performRequest(req)
	if err != nil {
		return respData, fmt.Errorf("perform request error %v", err)
	}

	respData, err = io.ReadAll(resp.Body)
	if err != nil {
		return respData, fmt.Errorf("perform request error %v", err)
	}

	return respData, nil
}

func (az *AzureTTSClient) SpeechToText(ctx context.Context,
	request model.SpeechToTextReq,
) (*model.SpeechToTextResp, error) {
	file, err := os.Open(request.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	url := fmt.Sprintf("%s?language=%s", az.SpeechToTextURL, request.Language)
	req, err := az.newSTTRequest(ctx, "POST", url, buf)
	if err != nil {
		return nil, fmt.Errorf("STT request error: %v", err)
	}

	resp, err := az.performRequest(req)
	if err != nil {
		return nil, fmt.Errorf("perform request error %v", err)
	}

	output := new(model.SpeechToTextResp)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}
