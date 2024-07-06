package api_test

import (
	"context"
	"os"
	"testing"

	tts "github.com/barkingdog-ai/azure-tts"
	"github.com/barkingdog-ai/azure-tts/model"
	"github.com/joho/godotenv"
)

func TestSpeechToText(t *testing.T) {
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		err := godotenv.Load("../.envrc")
		if err != nil {
			t.Fatalf("Failed to load .envrc file: %v", err)
		}
		apiKey = os.Getenv("AZURE_KEY")
		if apiKey == "" {
			t.Fatalf("Please set your AZURE_KEY environment variable")
		}
	}

	az, err := tts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		t.Fatalf("failed to create new client, received %v", err)
	}

	req := model.SpeechToTextReq{
		FilePath: "../data/test.mp4",
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
