package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	azuretts "github.com/barkingdog-ai/azure-tts"
	"github.com/barkingdog-ai/azure-tts/model"
)

func main() {
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		exit(fmt.Errorf("please set your AZURE_KEY environment variable"))
	}

	az, err := azuretts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		exit(fmt.Errorf("failed to create new client, received %v", err))
	}
	defer close(az.TokenRefreshDoneCh)

	ctx := context.Background()

	// 示範不同風格的使用方式
	examples := []struct {
		name     string
		text     string
		style    *model.TTSStyle
		filename string
	}{
		{
			name: "活潑熱情客服",
			text: "哈囉～您好！感謝您撥打我們的客服專線！今天過得還順利嗎？我很開心能為您服務，請問需要什麼幫忙呢？",
			style: &model.TTSStyle{
				Style:       "cheerful",
				StyleDegree: "2",
			},
			filename: "enthusiastic.mp3",
		},
		{
			name: "親切溫暖客服",
			text: "您好，我是客服專員，很高興為您服務。請您放心，把問題交給我，我會盡快幫您處理好。",
			style: &model.TTSStyle{
				Style:       "friendly",
				StyleDegree: "1",
			},
			filename: "warm.mp3",
		},
		{
			name: "正式專業客服",
			text: "您好，感謝您聯繫本公司客服中心。我很樂意協助您解決問題，請告訴我需要的服務內容。",
			style: &model.TTSStyle{
				Style:       "chat",
				StyleDegree: "1",
			},
			filename: "professional.mp3",
		},
		{
			name:     "無風格（基本）",
			text:     "您好，這是基本的語音合成，沒有特殊風格。",
			style:    nil,
			filename: "basic.mp3",
		},
	}

	voiceName := "zh-TW-HsiaoChenNeural"
	locale := "zh-TW"

	for _, example := range examples {
		fmt.Printf("正在生成: %s\n", example.name)

		req := model.TextToSpeechRequest{
			SpeechText:  example.text, // 直接使用文本，讓 API 內部處理 SSML
			Locale:      model.LocaleZhTW,
			Gender:      model.GenderFemale,
			VoiceName:   voiceName,
			AudioOutput: model.Audio16khz32kbitrateMonoMp3,
			Rate:        "1",
			Pitch:       "1",
			Style:       example.style,
		}

		b, err := az.TextToSpeech(ctx, &req)
		if err != nil {
			exit(fmt.Errorf("unable to synthesize %s, received: %v", example.name, err))
		}

		const filePermission = 0600
		err = ioutil.WriteFile(example.filename, b, filePermission)
		if err != nil {
			exit(fmt.Errorf("unable to write file %s, received %v", example.filename, err))
		}

		fmt.Printf("✓ 已生成: %s\n\n", example.filename)
	}

	// 示範預定義風格的使用
	fmt.Println("=== 使用預定義風格 ===")

	styleName := "enthusiastic" // 可以從外部參數傳入
	style, err := azuretts.GetStyleFromName(styleName)
	if err != nil {
		exit(fmt.Errorf("failed to get style: %v", err))
	}

	text := "歡迎使用我們的服務！"

	req := model.TextToSpeechRequest{
		SpeechText:  text, // 直接使用文本
		Locale:      model.LocaleZhTW,
		Gender:      model.GenderFemale,
		VoiceName:   voiceName,
		AudioOutput: model.Audio16khz32kbitrateMonoMp3,
		Rate:        "1",
		Pitch:       "1",
		Style:       style,
	}

	b, err := az.TextToSpeech(ctx, &req)
	if err != nil {
		exit(fmt.Errorf("unable to synthesize with predefined style, received: %v", err))
	}

	err = ioutil.WriteFile("predefined_style.mp3", b, 0600)
	if err != nil {
		exit(fmt.Errorf("unable to write predefined style file, received %v", err))
	}

	fmt.Printf("✓ 使用預定義風格 '%s' 生成: predefined_style.mp3\n", styleName)
	fmt.Printf("風格描述: %s\n", azuretts.GetStyleDescription(styleName))

	// 列出所有可用風格
	fmt.Println("\n=== 可用風格列表 ===")
	styles := azuretts.ListAvailableStyles()
	for name, desc := range styles {
		fmt.Printf("- %s: %s\n", name, desc)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
