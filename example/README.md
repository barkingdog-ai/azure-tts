# Azure TTS 台灣語音測試

這個專案用於測試 Azure Text-to-Speech 服務的台灣語音選項，特別針對客服場景進行優化。

## 問題解決

### 為什麼不同風格聽起來都一樣？

1. **風格強度不足**: 原本的 `StyleDegree` 設定太低 (1-2)，現在調整為 1.2-2.5
2. **語音參數固定**: 現在根據不同風格調整 `Rate` 和 `Pitch` 參數
3. **語音選擇**: 改用更適合台灣客服的 `zh-TW-HsiaoYuNeural`

### 語音選項比較

| 語音 ID | 特色 | 適用場景 |
|---------|------|----------|
| `zh-TW-HsiaoYuNeural` | 年輕活潑、親切自然 | 一般客服、親切服務 |
| `zh-TW-YunJheNeural` | 成熟專業、穩重可靠 | 正式客服、專業服務 |
| `zh-TW-HsiaoChenNeural` | 字正腔圓、清晰標準 | 正式場合、標準服務 |

## 使用方法

### 1. 設定環境變數

```bash
export AZURE_KEY=your_azure_speech_api_key
```

### 2. 執行語音風格測試

```bash
make run
```

這會生成以下檔案：
- `enthusiastic.mp3` - 活潑熱情客服
- `warm.mp3` - 親切溫暖客服  
- `professional.mp3` - 正式專業客服
- `basic.mp3` - 基本語音
- `predefined_style.mp3` - 預定義風格

### 3. 執行語音選項比較

```bash
# 台灣語音選項比較
make run-comparison

# 中文語音選項比較 (包含大陸語音，推薦)
make run-chinese
```

**台灣語音選項**：
- `voice_hsiaoyu.mp3` - 年輕活潑語音
- `voice_yunjhe.mp3` - 成熟專業語音
- `voice_hsiaochen.mp3` - 字正腔圓語音

**中文語音選項 (包含大陸語音)**：
- `voice_tw_hsiaoyu.mp3` - 台灣年輕活潑
- `voice_tw_yunjhe.mp3` - 台灣成熟專業
- `voice_cn_xiaoxiao.mp3` - 大陸年輕活潑
- `voice_cn_xiaoyi.mp3` - 大陸溫暖親切 ⭐ (推薦)
- `voice_cn_yunyang.mp3` - 大陸成熟穩重
- `voice_cn_xiaochen.mp3` - 大陸親切客服 ⭐ (推薦)

### 4. 清理生成的檔案

```bash
make clean
```

## 參數調整建議

### 風格強度 (StyleDegree)
- **活潑熱情**: 2.0-2.5
- **親切溫暖**: 1.5-2.0  
- **正式專業**: 1.0-1.5

### 語音參數
- **活潑熱情**: Rate=1.1, Pitch=1.05
- **親切溫暖**: Rate=0.95, Pitch=0.98
- **正式專業**: Rate=1.0, Pitch=1.0

## 檔案說明

- `main.go` - 主要語音風格測試程式
- `voice_comparison.go` - 語音選項比較程式
- `Makefile` - 便利的執行指令
- `README.md` - 本說明文件

## 注意事項

1. 確保已設定正確的 `AZURE_KEY` 環境變數
2. 需要網路連線來呼叫 Azure 語音服務
3. 生成的 MP3 檔案會保存在當前目錄
4. 建議先執行 `make run-comparison` 來選擇最適合的語音
