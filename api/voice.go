package api

import (
	"context"
	"github.com/barkingdog-ai/azure-tts/model"
)

type VoiceInterface interface {
	VoiceList(ctx context.Context) (*[]model.VoiceListResponse, error)
}

func (az *AzureTTSClient) VoiceList(ctx context.Context) (*[]model.VoiceListResponse, error) {
	return az.VoiceListRequest(ctx)
}

func (az *AzureTTSClient) VoiceListRequest(ctx context.Context) (*[]model.VoiceListResponse, error) {
	req, err := az.newRequest(ctx, "GET", az.VoiceServiceListURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := az.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new([]model.VoiceListResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}
