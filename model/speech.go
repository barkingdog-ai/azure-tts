package model

import "mime/multipart"

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
	File     *multipart.FileHeader
	Language string
}

type SpeechToTextResp struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	Offset            int    `json:"Offset"`
	Duration          int    `json:"Duration"`
	DisplayText       string `json:"DisplayText"`
}
