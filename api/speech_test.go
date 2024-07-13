package api_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	tts "github.com/barkingdog-ai/azure-tts"
	"github.com/barkingdog-ai/azure-tts/model"
	"github.com/joho/godotenv"
)

func TestTextToSpeech(t *testing.T) {
	var apiKey string
	if apiKey = os.Getenv("AZURE_API_KEY"); apiKey == "" {
		err := godotenv.Load("../.envrc")
		if err != nil {
			t.Fatalf("Failed to load .envrc file: %v", err)
		}
	}

	az, err := tts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		t.Fatalf("failed to create new client, received %v", err)
	}

	req := model.TextToSpeechRequest{
		SpeechText:  "你好123",
		Locale:      model.LocaleZhTW,
		Gender:      model.GenderFemale,
		VoiceName:   "zh-TW-HsiaoChenNeural",
		AudioOutput: model.Audio16khz32kbitrateMonoMp3,
		Rate:        "1",
		Pitch:       "1",
	}

	ctx := context.Background()
	resp, err := az.TextToSpeech(ctx, req)
	if err != nil {
		t.Fatalf("TextToSpeech failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Response: %+v", resp)
}

func TestSpeechToText(t *testing.T) {
	var apiKey string
	if apiKey = os.Getenv("AZURE_API_KEY"); apiKey == "" {
		err := godotenv.Load("../.envrc")
		if err != nil {
			t.Fatalf("Failed to load .envrc file: %v", err)
		}
	}

	az, err := tts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		t.Fatalf("failed to create new client, received %v", err)
	}
	file, err := os.Open("../data/test.mp4")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	req := model.SpeechToTextReq{
		Reader:   buf,
		Language: "zh-TW",
	}

	ctx := context.Background()
	resp, err := az.SpeechToText(ctx, req)
	if err != nil {
		t.Fatalf("SpeechToText failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Response: %+v", resp)
}
