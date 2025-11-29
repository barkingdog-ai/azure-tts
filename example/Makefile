# Azure TTS 語音測試 Makefile

.PHONY: help run run-comparison clean

help:
	@echo "Azure TTS 語音測試工具"
	@echo ""
	@echo "可用指令:"
	@echo "  make run           - 執行主要語音風格測試"
	@echo "  make run-comparison - 執行台灣語音選項比較測試"
	@echo "  make run-chinese   - 執行中文語音選項比較測試 (包含大陸語音)"
	@echo "  make run-test      - 執行熱情調皮句子測試"
	@echo "  make clean         - 清理生成的音檔"
	@echo ""
	@echo "環境變數:"
	@echo "  AZURE_KEY         - Azure 語音服務 API Key (必需)"

run:
	@echo "執行語音風格測試..."
	go run main.go

run-comparison:
	@echo "執行語音選項比較測試..."
	go run main.go compare

run-chinese:
	@echo "執行中文語音選項比較測試..."
	go run main.go chinese

run-test:
	@echo "執行熱情調皮句子測試..."
	go run main.go test

clean:
	@echo "清理生成的音檔..."
	rm -f *.mp3
	@echo "清理完成"

# 檢查環境變數
check-env:
	@if [ -z "$$AZURE_KEY" ]; then \
		echo "錯誤: 請設定 AZURE_KEY 環境變數"; \
		echo "例如: export AZURE_KEY=your_api_key_here"; \
		exit 1; \
	fi
	@echo "環境變數檢查通過"
