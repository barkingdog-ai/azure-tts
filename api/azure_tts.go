package api

import (
	"net/http"
)

type AzureTTSClient struct {
	HTTPClient          *http.Client
	AccessToken         string
	SubscriptionKey     string
	TokenRefreshDoneCh  chan bool
	TokenRefreshURL     string
	VoiceServiceListURL string
	TextToSpeechURL     string
	SpeechToTextURL     string
}
