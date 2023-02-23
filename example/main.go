package main

import (
	"context"
	"fmt"
	tts "github.com/barkingdog-ai/azure-tts"
	"github.com/barkingdog-ai/azure-tts/model"
	"io/ioutil"
	"os"
)

func main() {
	// create a key for "Cognitive Services" (kind=SpeechServices). Once the key is available
	// in the Azure portal, push it into an environment variable (export AZUREKEY=SYS64738).
	// By default the free tier keys are served out of West US2
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		exit(fmt.Errorf("Please set your AZURE_KEY environment variable"))
	}

	az, err := tts.NewClient(apiKey, model.RegionEastUS)
	if err != nil {
		exit(fmt.Errorf("failed to create new client, received %v", err))
	}
	defer close(az.TokenRefreshDoneCh)

	// Digitize a text string using the enUS locale, female voice and specify the
	// audio format of a 16Khz, 32kbit mp3 file.
	ctx := context.Background()
	req := model.TextToSpeechRequest{
		SpeechText:  "Hello world",
		Locale:      model.LocaleZhTW,
		Gender:      model.GenderFemale,
		VoiceName:   "en-US-ChristopherNeural",
		AudioOutput: model.Audio16khz32kbitrateMonoMp3,
	}
	b, err := az.TextToSpeech(ctx, req)

	if err != nil {
		exit(fmt.Errorf("unable to synthesize, received: %v", err))
	}

	// send results to disk.
	err = ioutil.WriteFile("audio.mp3", b, 0644)
	if err != nil {
		exit(fmt.Errorf("unable to write file, received %v", err))
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
