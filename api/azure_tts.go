package api

import (
	"net/http"
)

type AzureTTSClient struct {
	HttpClient          *http.Client
	AccessToken         string
	SubscriptionKey     string
	TokenRefreshDoneCh  chan bool
	TokenRefreshURL     string
	VoiceServiceListURL string
	TextToSpeechURL     string
	SpeechToTextURL     string
}
