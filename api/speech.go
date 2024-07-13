package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/barkingdog-ai/azure-tts/model"
	"github.com/barkingdog-ai/azure-tts/utils"
)

type SpeechInterface interface {
	TextToSpeech(ctx context.Context, req model.TextToSpeechRequest) ([]byte, error)
	SpeechToText(ctx context.Context, req model.SpeechToTextReq) (*model.SpeechToTextResp, error)
	CorrectHomophones(req model.TextToSpeechRequest) (resp model.TextToSpeechRequest)
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

	request = az.CorrectHomophones(request)

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
	url := fmt.Sprintf("%s?language=%s", az.SpeechToTextURL, request.Language)

	payload, err := createFilePayload(request)
	if err != nil {
		return nil, err
	}

	req, err := az.newSTTRequest(ctx, "POST", url, payload)
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

func (az *AzureTTSClient) CorrectHomophones(req model.TextToSpeechRequest) (resp model.TextToSpeechRequest) {
	for _, homophone := range req.Homophones {
		req.SpeechText = strings.ReplaceAll(req.SpeechText, homophone.TargetText, homophone.ReplaceText)
	}
	return req
}
