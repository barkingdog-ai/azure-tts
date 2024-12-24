package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	tts "github.com/barkingdog-ai/azure-tts"
	"github.com/barkingdog-ai/azure-tts/model"
)

func main() {
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		exit(fmt.Errorf("please set your AZURE_KEY environment variable"))
	}

	az, err := tts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		exit(fmt.Errorf("failed to create new client, received %v", err))
	}
	defer close(az.TokenRefreshDoneCh)
	ctx := context.Background()
	req := model.TextToSpeechRequest{
		SpeechText:  "版本号是8.2.3和2024.13.32",
		Locale:      model.LocaleZhTW,
		Gender:      model.GenderFemale,
		VoiceName:   "zh-TW-HsiaoChenNeural",
		AudioOutput: model.Audio16khz32kbitrateMonoMp3,
		Rate:        "1",
		Pitch:       "1",
	}
	b, err := az.TextToSpeech(ctx, &req)
	if err != nil {
		exit(fmt.Errorf("unable to synthesize, received: %v", err))
	}

	const filePermission = 0600
	err = ioutil.WriteFile("audio1.mp3", b, filePermission)
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
