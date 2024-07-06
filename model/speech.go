package model

import "io"

// TextToSpeechRequest is a request for the TextToSpeech API
type TextToSpeechRequest struct {
	SpeechText  string
	Locale      Locale
	Gender      Gender
	VoiceName   string
	AudioOutput AudioOutput
	Rate        string
	Pitch       string
}

type TextToSpeechResponse struct {
	Audio []byte
}

type SpeechToTextReq struct {
	Reader   io.Reader
	FilePath string
	Language string
}

type SpeechToTextResp struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	Offset            int    `json:"Offset"`
	Duration          int    `json:"Duration"`
	DisplayText       string `json:"DisplayText"`
}
