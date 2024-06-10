package model

type voiceType int

type VoiceListResponse struct {
	Name            string    `json:"Name"`
	LocalName       string    `json:"LocalName"`
	ShortName       string    `json:"ShortName"`
	Gender          string    `json:"Gender"`
	Locale          string    `json:"Locale"`
	LocaleName      string    `json:"LocaleName"`
	SampleRateHertz string    `json:"SampleRateHertz"`
	VoiceType       voiceType `json:"VoiceType"`
	StyleList       []string  `json:"StyleList"`
	Status          string    `json:"Status"`
	WordsPerMinute  string    `json:"WordsPerMinute"`
}
