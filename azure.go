package azuretts

import (
	"context"
	"fmt"
	"net/http"
	"time"

	API "github.com/barkingdog-ai/azure-tts/api"
	"github.com/barkingdog-ai/azure-tts/model"
)

const (
	// voiceListAPI is the source for supported voice list to region mapping
	// See: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#regions-and-endpoints
	voiceListAPI = "https://%s.tts.speech.microsoft.com/cognitiveservices/voices/list"
	// The following are V1 endpoints for Cognitiveservices endpoints.
	textToSpeechAPI = "https://%s.tts.speech.microsoft.com/cognitiveservices/v1"
	speechToTextAPI = "https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1"
	refreshAPI      = "https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken"
)

// synthesizeActionTimeout is the amount of time the http client will wait for a response during Synthesize request.
const synthesizeActionTimeout = time.Second * 30

const (
	defaultTimeoutSeconds = 30
)

type Interface interface {
	API.SpeechInterface
	API.VoiceInterface
	API.TokenInterface
}

func NewClient(subscriptionKey string, region model.Region, options ...API.ClientOption) (*API.AzureTTSClient, error) {
	httpClient := &http.Client{
		Timeout: defaultTimeoutSeconds * time.Second,
	}

	az := &API.AzureTTSClient{
		SubscriptionKey: subscriptionKey,
		HTTPClient:      httpClient,
	}
	az.TextToSpeechURL = fmt.Sprintf(textToSpeechAPI, region)
	az.SpeechToTextURL = fmt.Sprintf(speechToTextAPI, region)
	az.TokenRefreshURL = fmt.Sprintf(refreshAPI, region)
	az.VoiceServiceListURL = fmt.Sprintf(voiceListAPI, region)

	// api requires that the token is refreshed every 10 mintutes.
	// We will do this task in the background every ~9 minutes.
	ctx, cancel := context.WithTimeout(context.Background(), synthesizeActionTimeout)
	defer cancel()
	err := az.RefreshToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial token, %v", err)
	}

	az.TokenRefreshDoneCh = startRefresher(ctx, az)

	for _, o := range options {
		if err := o(az); err != nil {
			return nil, err
		}
	}
	return az, nil
}

// startRefresher updates the authentication token on at a 9 minute interval. A channel is returned
// if the caller wishes to cancel the channel.
func startRefresher(ctx context.Context, az *API.AzureTTSClient) chan bool {
	done := make(chan bool, 1)
	go func() {
		const refreshInterval = time.Minute * 9
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := az.RefreshToken(ctx)
				if err != nil {
					_ = fmt.Sprintf("failed to refresh token, %v", err)
				}
			case <-done:
				return
			}
		}
	}()
	return done
}
