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
	// 檢查是否有參數指定要比較語音
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "compare":
			compareVoices()
			return
		case "chinese":
			compareChineseVoices()
			return
		case "test":
			testEnthusiasticSentences()
			return
		}
	}
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
			text: "您好，歡迎致電東東婚宴會館。請撥分機號碼, 或播9由總機為您服務!",
			style: &model.TTSStyle{
				Style:       "chat",
				StyleDegree: "1.1", // 更強的開心熱情風格
			},
			filename: "welcome.wav",
		},
		// {
		// 	name: "活潑熱情客服",
		// 	text: "哈囉～您好！感謝您撥打我們的客服專線！今天過得還順利嗎？我很開心能為您服務，請問需要什麼幫忙呢？",
		// 	style: &model.TTSStyle{
		// 		Style:       "cheerful",
		// 		StyleDegree: "1.6", // 更強的開心熱情風格
		// 	},
		// 	filename: "enthusiastic.mp3",
		// },
		// {
		// 	name: "親切溫暖客服",
		// 	text: "您好，我是客服專員，很高興為您服務。請您放心，把問題交給我，我會盡快幫您處理好。",
		// 	style: &model.TTSStyle{
		// 		Style:       "friendly",
		// 		StyleDegree: "1.3", // 更明顯的溫暖上揚感
		// 	},
		// 	filename: "warm.mp3",
		// },
		// {
		// 	name: "正式專業客服",
		// 	text: "您好，感謝您聯繫本公司客服中心。我很樂意協助您解決問題，請告訴我需要的服務內容。",
		// 	style: &model.TTSStyle{
		// 		Style:       "chat",
		// 		StyleDegree: "1.1", // 更明顯的積極上揚風格
		// 	},
		// 	filename: "professional.mp3",
		// },
		// {
		// 	name:     "無風格（基本）",
		// 	text:     "您好，這是基本的語音合成，沒有特殊風格。",
		// 	style:    nil,
		// 	filename: "basic.mp3",
		// },
	}

	// 統一使用 zh-TW-HsiaoChenNeural 語音，調整不同風格
	voiceName := "zh-TW-HsiaoChenNeural"

	for _, example := range examples {
		fmt.Printf("正在生成: %s\n", example.name)

		// 根據不同風格調整語音參數 - 尾調上揚，開心熱情
		var rate, pitch string
		switch example.name {
		case "活潑熱情客服":
			rate = "1.25"  // 稍微慢一點，保持順暢
			pitch = "1.15" // 提高音調，尾調上揚
		case "親切溫暖客服":
			rate = "1.2"   // 適中速度，保持親切
			pitch = "1.12" // 提高音調，溫暖上揚
		case "正式專業客服":
			rate = "1.15"  // 稍微快一點，保持專業
			pitch = "1.08" // 提高音調，積極上揚
		default:
			rate = "1.2"   // 預設適中速度
			pitch = "1.12" // 預設提高音調，尾調上揚
		}

		req := model.TextToSpeechRequest{
			SpeechText:  example.text, // 直接使用文本，讓 API 內部處理 SSML
			Locale:      model.LocaleZhTW,
			Gender:      model.GenderFemale,
			VoiceName:   voiceName,
			AudioOutput: model.Audio16khz32kbitrateMonoMp3,
			Rate:        rate,
			Pitch:       pitch,
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
		Rate:        "1.2",  // 適中速度，保持順暢
		Pitch:       "1.12", // 提高音調，尾調上揚
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

	// 語音選項說明
	fmt.Println("\n=== 語音選項說明 ===")
	fmt.Println("目前使用: zh-TW-HsiaoChenNeural (統一語音，不同風格)")
	fmt.Println("特色: 順暢、尾調上揚、開心熱情語調")
	fmt.Println("\n風格強度調整:")
	fmt.Println("- 活潑熱情: StyleDegree 3.2 (最開心熱情)")
	fmt.Println("- 親切溫暖: StyleDegree 2.5 (溫暖上揚)")
	fmt.Println("- 正式專業: StyleDegree 2.0 (積極上揚)")
	fmt.Println("\n語音參數調整 (尾調上揚，開心熱情):")
	fmt.Println("- 活潑熱情: Rate=1.25, Pitch=1.15 (順暢且上揚)")
	fmt.Println("- 親切溫暖: Rate=1.2, Pitch=1.12 (溫暖且上揚)")
	fmt.Println("- 正式專業: Rate=1.15, Pitch=1.08 (專業且上揚)")
}

func testEnthusiasticSentences() {
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		fmt.Fprintf(os.Stderr, "ERROR: please set your AZURE_KEY environment variable\n")
		os.Exit(1)
	}

	az, err := azuretts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create new client, received %v\n", err)
		os.Exit(1)
	}
	defer close(az.TokenRefreshDoneCh)

	ctx := context.Background()

	// 熱情調皮測試句子
	testSentences := []struct {
		name     string
		text     string
		filename string
	}{
		{
			name:     "熱情開場",
			text:     "哈囉～歡迎來到我們的客服中心！我是您的小助手，今天要為您帶來最棒的服務體驗！",
			filename: "test_enthusiastic_1.mp3",
		},
		{
			name:     "調皮互動",
			text:     "嘿嘿～您問得真好！這個問題我已經準備好答案了，準備好被驚艷到了嗎？",
			filename: "test_enthusiastic_2.mp3",
		},
		{
			name:     "溫暖關懷",
			text:     "親愛的客戶，您今天過得還好嗎？有什麼煩惱都可以告訴我，我會用心傾聽的！",
			filename: "test_enthusiastic_3.mp3",
		},
		{
			name:     "積極鼓勵",
			text:     "太棒了！您做得非常好！我為您感到驕傲！繼續保持這個好狀態吧！",
			filename: "test_enthusiastic_4.mp3",
		},
		{
			name:     "調皮結尾",
			text:     "好了～問題解決了！您是不是覺得我很厲害呢？嘿嘿～期待下次再為您服務！",
			filename: "test_enthusiastic_5.mp3",
		},
		{
			name:     "特殊情境",
			text:     "哎呀～您真是太可愛了！這個問題問得我都要笑出來了！讓我來為您解答吧！",
			filename: "test_enthusiastic_6.mp3",
		},
		{
			name:     "情感表達",
			text:     "太感動了！能為您服務我真的好開心！您讓我覺得今天特別有意義！",
			filename: "test_enthusiastic_7.mp3",
		},
		{
			name:     "調皮問候",
			text:     "嗨嗨～我是您的小太陽客服！今天要為您帶來滿滿的正能量！準備好了嗎？",
			filename: "test_enthusiastic_8.mp3",
		},
	}

	voiceName := "zh-TW-HsiaoChenNeural"

	for i, sentence := range testSentences {
		fmt.Printf("正在生成: %s (%d/8)\n", sentence.name, i+1)

		req := model.TextToSpeechRequest{
			SpeechText:  sentence.text,
			Locale:      model.LocaleZhTW,
			Gender:      model.GenderFemale,
			VoiceName:   voiceName,
			AudioOutput: model.Audio16khz32kbitrateMonoMp3,
			Rate:        "1.25", // 順暢語速
			Pitch:       "1.15", // 尾調上揚
			Style: &model.TTSStyle{
				Style:       "cheerful",
				StyleDegree: "3.2", // 最開心熱情
			},
		}

		b, err := az.TextToSpeech(ctx, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to synthesize %s, received: %v\n", sentence.name, err)
			os.Exit(1)
		}

		const filePermission = 0600
		err = ioutil.WriteFile(sentence.filename, b, filePermission)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to write file %s, received %v\n", sentence.filename, err)
			os.Exit(1)
		}

		fmt.Printf("✓ 已生成: %s\n\n", sentence.filename)
	}

	fmt.Println("=== 熱情調皮測試完成 ===")
	fmt.Println("請播放以下檔案來測試熱情調皮語音:")
	for _, sentence := range testSentences {
		fmt.Printf("- %s: %s\n", sentence.filename, sentence.name)
	}
	fmt.Println("\n所有句子都使用以下參數:")
	fmt.Println("- 語音: zh-TW-HsiaoChenNeural")
	fmt.Println("- 語速: 1.25 (順暢)")
	fmt.Println("- 音調: 1.15 (尾調上揚)")
	fmt.Println("- 風格: cheerful, StyleDegree=3.2 (最開心熱情)")
}

// compareVoices 比較不同台灣語音選項
func compareVoices() {
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		fmt.Fprintf(os.Stderr, "ERROR: please set your AZURE_KEY environment variable\n")
		os.Exit(1)
	}

	az, err := azuretts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create new client, received %v\n", err)
		os.Exit(1)
	}
	defer close(az.TokenRefreshDoneCh)

	ctx := context.Background()

	// 台灣語音選項比較
	voices := []struct {
		name     string
		voiceID  string
		desc     string
		filename string
	}{
		{
			name:     "年輕活潑",
			voiceID:  "zh-TW-HsiaoYuNeural",
			desc:     "年輕、活潑、親切，適合客服",
			filename: "voice_hsiaoyu.mp3",
		},
		{
			name:     "成熟專業",
			voiceID:  "zh-TW-YunJheNeural",
			desc:     "成熟、專業、穩重，適合正式場合",
			filename: "voice_yunjhe.mp3",
		},
		{
			name:     "字正腔圓",
			voiceID:  "zh-TW-HsiaoChenNeural",
			desc:     "字正腔圓、清晰，適合正式客服",
			filename: "voice_hsiaochen.mp3",
		},
	}

	// 測試文本
	testText := "您好，歡迎致電客服中心。我是您的專屬客服，很高興為您服務。請問有什麼可以幫助您的嗎？"

	for _, voice := range voices {
		fmt.Printf("正在生成: %s (%s)\n", voice.name, voice.desc)

		req := model.TextToSpeechRequest{
			SpeechText:  testText,
			Locale:      model.LocaleZhTW,
			Gender:      model.GenderFemale,
			VoiceName:   voice.voiceID,
			AudioOutput: model.Audio16khz32kbitrateMonoMp3,
			Rate:        "1.2", // 加快速度
			Pitch:       "1.1", // 提高音調，更開心
			Style: &model.TTSStyle{
				Style:       "friendly",
				StyleDegree: "1.8", // 提高風格強度，更明顯
			},
		}

		b, err := az.TextToSpeech(ctx, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to synthesize %s, received: %v\n", voice.name, err)
			os.Exit(1)
		}

		const filePermission = 0600
		err = ioutil.WriteFile(voice.filename, b, filePermission)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to write file %s, received %v\n", voice.filename, err)
			os.Exit(1)
		}

		fmt.Printf("✓ 已生成: %s\n\n", voice.filename)
	}

	fmt.Println("=== 語音比較完成 ===")
	fmt.Println("請播放以下檔案來比較不同語音:")
	for _, voice := range voices {
		fmt.Printf("- %s: %s\n", voice.filename, voice.desc)
	}
}

// compareChineseVoices 比較不同中文語音選項 (包含大陸語音)
func compareChineseVoices() {
	var apiKey string
	if apiKey = os.Getenv("AZURE_KEY"); apiKey == "" {
		fmt.Fprintf(os.Stderr, "ERROR: please set your AZURE_KEY environment variable\n")
		os.Exit(1)
	}

	az, err := azuretts.NewClient(apiKey, model.RegionEastAsia)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create new client, received %v\n", err)
		os.Exit(1)
	}
	defer close(az.TokenRefreshDoneCh)

	ctx := context.Background()

	// 中文語音選項比較 - 包含大陸和台灣語音
	voices := []struct {
		name     string
		voiceID  string
		desc     string
		filename string
		locale   model.Locale
	}{
		// 台灣語音
		{
			name:     "台灣-年輕活潑",
			voiceID:  "zh-TW-HsiaoYuNeural",
			desc:     "台灣年輕活潑，親切自然",
			filename: "voice_tw_hsiaoyu.mp3",
			locale:   model.LocaleZhTW,
		},
		{
			name:     "台灣-成熟專業",
			voiceID:  "zh-TW-YunJheNeural",
			desc:     "台灣成熟專業，穩重可靠",
			filename: "voice_tw_yunjhe.mp3",
			locale:   model.LocaleZhTW,
		},
		// 大陸語音 - 更自然的選項
		{
			name:     "大陸-年輕活潑",
			voiceID:  "zh-CN-XiaoxiaoNeural",
			desc:     "大陸年輕活潑，自然親切",
			filename: "voice_cn_xiaoxiao.mp3",
			locale:   model.LocaleZhCN,
		},
		{
			name:     "大陸-溫暖親切",
			voiceID:  "zh-CN-XiaoyiNeural",
			desc:     "大陸溫暖親切，不字正腔圓",
			filename: "voice_cn_xiaoyi.mp3",
			locale:   model.LocaleZhCN,
		},
		{
			name:     "大陸-成熟穩重",
			voiceID:  "zh-CN-YunyangNeural",
			desc:     "大陸成熟穩重，自然專業",
			filename: "voice_cn_yunyang.mp3",
			locale:   model.LocaleZhCN,
		},
		{
			name:     "大陸-親切客服",
			voiceID:  "zh-CN-XiaochenNeural",
			desc:     "大陸親切客服，溫暖自然",
			filename: "voice_cn_xiaochen.mp3",
			locale:   model.LocaleZhCN,
		},
	}

	// 測試文本 - 客服常用語句
	testText := "您好，歡迎致電客服中心。我是您的專屬客服，很高興為您服務。請問有什麼可以幫助您的嗎？"

	for _, voice := range voices {
		fmt.Printf("正在生成: %s (%s)\n", voice.name, voice.desc)

		req := model.TextToSpeechRequest{
			SpeechText:  testText,
			Locale:      voice.locale,
			Gender:      model.GenderFemale,
			VoiceName:   voice.voiceID,
			AudioOutput: model.Audio16khz32kbitrateMonoMp3,
			Rate:        "1.2", // 加快速度
			Pitch:       "1.1", // 提高音調，更開心
			Style: &model.TTSStyle{
				Style:       "friendly",
				StyleDegree: "1.8", // 提高風格強度，更明顯
			},
		}

		b, err := az.TextToSpeech(ctx, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to synthesize %s, received: %v\n", voice.name, err)
			os.Exit(1)
		}

		const filePermission = 0600
		err = ioutil.WriteFile(voice.filename, b, filePermission)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to write file %s, received %v\n", voice.filename, err)
			os.Exit(1)
		}

		fmt.Printf("✓ 已生成: %s\n\n", voice.filename)
	}

	fmt.Println("=== 中文語音比較完成 ===")
	fmt.Println("請播放以下檔案來比較不同語音:")
	fmt.Println("\n台灣語音:")
	fmt.Println("- voice_tw_hsiaoyu.mp3: 台灣年輕活潑")
	fmt.Println("- voice_tw_yunjhe.mp3: 台灣成熟專業")
	fmt.Println("\n大陸語音 (推薦，不字正腔圓):")
	fmt.Println("- voice_cn_xiaoxiao.mp3: 大陸年輕活潑")
	fmt.Println("- voice_cn_xiaoyi.mp3: 大陸溫暖親切 ⭐")
	fmt.Println("- voice_cn_yunyang.mp3: 大陸成熟穩重")
	fmt.Println("- voice_cn_xiaochen.mp3: 大陸親切客服 ⭐")
	fmt.Println("\n建議優先試聽標記 ⭐ 的語音選項")
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}


ffmpeg -i /home/claire/repo/gittea/teleagent/gopkg/azure-tts/example/welcome.wav -ar 8000 -ac 1 -sample_fmt s16 /home/claire/repo/gittea/teleagent/gopkg/azure-tts/example/welcome_8k.wav
