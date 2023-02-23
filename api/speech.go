package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/barkingdog-ai/azure-tts/model"
	"io/ioutil"
)

type SpeechInterface interface {
	TextToSpeech(ctx context.Context, req model.TextToSpeechRequest) (*model.TextToSpeechResponse, error)
}

func (az *AzureTTSClient) TextToSpeech(ctx context.Context, request model.TextToSpeechRequest) ([]byte, error) {

	v := voiceXML(
		request.SpeechText,
		request.VoiceName,
		request.Locale,
		request.Gender)

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
