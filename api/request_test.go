package api

import (
	"testing"

	"github.com/barkingdog-ai/azure-tts/model"
	"github.com/stretchr/testify/assert"
)

func Test_voiceXML(t *testing.T) {
	tests := []struct {
		name        string
		speechText  string
		description string
		locale      model.Locale
		gender      model.Gender
		rate        string
		pitch       string
		want        string
	}{
		{
			name:        "IP地址测试",
			speechText:  "服务器IP是192.168.1.1和10.0.0.1",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>服务器IP是<say-as interpret-as="characters">192.168.1.1</say-as>和<say-as interpret-as="characters">10.0.0.1</say-as></prosody></voice></speak>`,
		},
		{
			name:        "普通URL测试",
			speechText:  "请访问http://api.ai-amaze.com查看详情",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>请访问api.ai-amaze.com查看详情</prosody></voice></speak>`,
		},
		{
			name:        "IP形式URL测试",
			speechText:  "请访问http://192.168.1.1查看详情",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>请访问<say-as interpret-as="characters">192.168.1.1</say-as>查看详情</prosody></voice></speak>`,
		},
		{
			name:        "Markdown URL测试",
			speechText:  "点击[官方网站](http://api.ai-amaze.com)了解更多",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>点击官方网站了解更多</prosody></voice></speak>`,
		},
		{
			name:        "混合测试",
			speechText:  "服务器IP是192.168.1.1，请访问http://api.ai-amaze.com或[官方网站](http://api.ai-amaze.com)了解更多",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>服务器IP是<say-as interpret-as="characters">192.168.1.1</say-as>，请访问api.ai-amaze.com或官方网站了解更多</prosody></voice></speak>`,
		},
		{
			name:        "有效日期格式测试",
			speechText:  "日期是2024.03.15和24.03.15",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>日期是2024.03.15和24.03.15</prosody></voice></speak>`,
		},
		{
			name:        "无效日期格式测试",
			speechText:  "版本号是8.2.3和2024.13.32",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>版本号是<say-as interpret-as="characters">8.2.3</say-as>和<say-as interpret-as="characters">2024.13.32</say-as></prosody></voice></speak>`,
		},
		{
			name:        "混合日期和版本号测试",
			speechText:  "更新时间是2024.03.15，当前版本8.2.3",
			description: "test",
			locale:      model.LocaleZhCN,
			gender:      model.GenderFemale,
			rate:        "0%",
			pitch:       "0%",
			want:        `<speak version='1.0' xml:lang='zh-CN'><voice xml:lang='zh-CN' xml:gender='Female' name='test'><prosody rate='0%' pitch='0%'>更新时间是2024.03.15，当前版本<say-as interpret-as="characters">8.2.3</say-as></prosody></voice></speak>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := voiceXML(tt.speechText, tt.description, tt.locale, tt.gender, tt.rate, tt.pitch)
			assert.Equal(t, tt.want, got)
		})
	}
}
