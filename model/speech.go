package model

import (
	"io"
)

type TextToSpeechRequest struct {
	SpeechText  string
	Locale      Locale
	Gender      Gender
	VoiceName   string
	AudioOutput AudioOutput
	Rate        string
	Pitch       string
	Homophones  []Homophones
}

type Homophones struct {
	TargetText  string `json:"target_text"`
	ReplaceText string `json:"replace_text"`
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
