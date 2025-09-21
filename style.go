package azuretts

import (
	"fmt"
	"strings"

	"github.com/barkingdog-ai/azure-tts/model"
)

// 預定義的 TTS 風格
const (
	StyleCheerful = "cheerful" // 活潑熱情
	StyleFriendly = "friendly" // 親切溫暖
	StyleChat     = "chat"     // 對話聊天
	StyleCalm     = "calm"     // 冷靜專業
)

// TTS 風格配置
type StyleConfig struct {
	Name        string
	Style       string
	StyleDegree string
	Description string
}

// 預定義風格配置
var PredefinedStyles = map[string]StyleConfig{
	"enthusiastic": {
		Name:        "enthusiastic",
		Style:       StyleCheerful,
		StyleDegree: "2",
		Description: "活潑熱情客服，適合第一時間接聽",
	},
	"warm": {
		Name:        "warm",
		Style:       StyleFriendly,
		StyleDegree: "1",
		Description: "親切溫暖客服，給客戶安心感",
	},
	"professional": {
		Name:        "professional",
		Style:       StyleChat,
		StyleDegree: "1",
		Description: "正式又保持活潑，適合企業客服",
	},
	"calm": {
		Name:        "calm",
		Style:       StyleCalm,
		StyleDegree: "1",
		Description: "冷靜專業，適合處理複雜問題",
	},
}

// 創建帶風格的 SSML 文本
func CreateStyledSSML(text, voiceName string, style *model.TTSStyle, locale string) string {
	if style == nil {
		// 如果沒有風格，返回基本的 SSML
		return fmt.Sprintf(`<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang="%s">
  <voice name="%s">%s</voice>
</speak>`, locale, voiceName, text)
	}

	// 帶風格的 SSML
	return fmt.Sprintf(`<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis"
       xmlns:mstts="https://www.w3.org/mstts"
       xml:lang="%s">
  <voice name="%s">
    <mstts:express-as style="%s" styledegree="%s">
      %s
    </mstts:express-as>
  </voice>
</speak>`, locale, voiceName, style.Style, style.StyleDegree, text)
}

// 從預定義風格名稱獲取風格配置
func GetStyleFromName(styleName string) (*model.TTSStyle, error) {
	config, exists := PredefinedStyles[styleName]
	if !exists {
		return nil, fmt.Errorf("unknown style: %s. Available styles: %s",
			styleName, strings.Join(getAvailableStyleNames(), ", "))
	}

	return &model.TTSStyle{
		Style:       config.Style,
		StyleDegree: config.StyleDegree,
	}, nil
}

// 獲取所有可用風格名稱
func getAvailableStyleNames() []string {
	names := make([]string, 0, len(PredefinedStyles))
	for name := range PredefinedStyles {
		names = append(names, name)
	}
	return names
}

// 獲取風格描述
func GetStyleDescription(styleName string) string {
	if config, exists := PredefinedStyles[styleName]; exists {
		return config.Description
	}
	return "未知風格"
}

// 列出所有可用風格
func ListAvailableStyles() map[string]string {
	styles := make(map[string]string)
	for name, config := range PredefinedStyles {
		styles[name] = config.Description
	}
	return styles
}
