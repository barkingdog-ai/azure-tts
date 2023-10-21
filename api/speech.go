package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/barkingdog-ai/azure-tts/model"
	"github.com/barkingdog-ai/azure-tts/utils"
	"io/ioutil"
)

type SpeechInterface interface {
	TextToSpeech(ctx context.Context, req model.TextToSpeechRequest) ([]byte, error)
}

func (az *AzureTTSClient) TextToSpeech(ctx context.Context, request model.TextToSpeechRequest) ([]byte, error) {

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
		return nil, fmt.Errorf("tts request error %v", err)
	}

	resp, err := az.performRequest(req)
	if err != nil {
		return nil, fmt.Errorf("perform request error %v", err)
	}

	return ioutil.ReadAll(resp.Body)

}
