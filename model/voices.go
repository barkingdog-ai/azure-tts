package model

type voiceType int

const (
	voiceStandard voiceType = iota // Standard
	voiceNeural                    // Neural
)

type VoiceListResponse struct {
	Name            string    `json:"Name"`
	LocalName       string    `json:"LocalName"`
	ShortName       string    `json:"ShortName"`
	Gender          Gender    `json:"Gender"`
	Locale          Locale    `json:"Locale"`
	LocaleName      string    `json:"LocaleName"`
	SampleRateHertz string    `json:"SampleRateHertz"`
	VoiceType       voiceType `json:"VoiceType"`
	StyleList       []string  `json:"StyleList"`
	Status          string    `json:"Status"`
	WordsPerMinute  string    `json:"WordsPerMinute"`
}

// supportedVoices represents the key used within the `localeToGender` map.
type supportedVoices struct {
	Gender Gender
	Locale Locale
}
