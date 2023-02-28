package model

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
