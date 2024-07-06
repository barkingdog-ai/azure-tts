package azuretts_test

import (
	"context"
	"os"
	"testing"

	tts "github.com/barkingdog-ai/azure-tts"
	"github.com/barkingdog-ai/azure-tts/model"
)

const (
	testSpeechText  = "你好,我是AI客服嘻嘻,您需要什麼協助嗎?"
	testLocale      = model.LocaleZhTW
	testGender      = model.GenderFemale
	testVoiceName   = "zh-TW-HsiaoChenNeural"
	testAudioOutput = model.AudioRAW8Bit8kHzMonoMulaw
	testRate        = "1.15"
	testPitch       = "1"
)

func TestNewClientAndTextToSpeech(t *testing.T) {
	apiKey := os.Getenv("AZURE_API_KEY")
	if apiKey == "" {
		t.Fatal("AZURE_API_KEY environment variable is required")
	}

	az, err := tts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		t.Fatalf("failed to create new client, received %v", err)
	}
	defer close(az.TokenRefreshDoneCh)

	req := model.TextToSpeechRequest{
		SpeechText:  testSpeechText,
		Locale:      testLocale,
		Gender:      testGender,
		VoiceName:   testVoiceName,
		AudioOutput: testAudioOutput,
		Rate:        testRate,
		Pitch:       testPitch,
	}

	ctx := context.Background()

	t.Run("SynthesizeSpeech", func(t *testing.T) {
		audioData, err := az.TextToSpeech(ctx, req)
		if err != nil {
			t.Fatalf("unable to synthesize, received: %v", err)
		}

		outputFile := "output.raw"
		err = os.WriteFile(outputFile, audioData, 0644)
		if err != nil {
			t.Fatalf("failed to write audio data to file: %v", err)
		}

		t.Logf("Audio synthesized and written to %s", outputFile)
		// defer os.Remove(outputFile)

		fileInfo, err := os.Stat(outputFile)
		if err != nil {
			t.Fatalf("failed to stat output file: %v", err)
		}

		if fileInfo.Size() == 0 {
			t.Fatalf("output file is empty")
		}
	})
}
