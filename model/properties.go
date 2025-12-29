package model

import "fmt"

// AudioOutput types represent the supported audio encoding formats for the text-to-speech endpoint.
// This type is required when requesting to azuretexttospeech.Synthesize text-to-speed request.
// Each incorporates a bitrate and encoding type. The Speech service supports 24 kHz, 16 kHz, and 8 kHz audio outputs.
// See: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#audio-outputs
type AudioOutput int

const (
	AudioRIFF8Bit8kHzMonoPCM AudioOutput = iota
	AudioRIFF16Bit16kHzMonoPCM
	AudioRIFF16khz16kbpsMonoSiren
	AudioRIFF24khz16bitMonoPcm
	AudioRAW8Bit8kHzMonoMulaw
	AudioRAW16Bit16kHzMonoPcm
	AudioRAW24khz16bitMonoPcm
	AudioRAW22050hz16bitMonoPcm
	AudioSsml16khz16bitMonoTts
	Audio16khz16kbpsMonoSiren
	Audio16khz32kbitrateMonoMp3
	Audio6khz64kbitrateMonoMp3
	Audio16khz128kbitrateMonoMp3
	Audio24khz48kbitrateMonoMp3
	Audio24khz96kbitrateMonoMp3
)

func (a AudioOutput) String() string {
	return []string{
		"riff-8khz-8bit-mono-mulaw",
		"riff-16khz-16bit-mono-pcm",
		"riff-16khz-16kbps-mono-siren",
		"riff-24khz-16bit-mono-pcm",
		"raw-8khz-8bit-mono-mulaw",
		"raw-16khz-16bit-mono-pcm",
		"raw-24khz-16bit-mono-pcm",
		"raw-22050hz-16bit-mono-pcm",
		"ssml-16khz-16bit-mono-tts",
		"audio-16khz-16kbps-mono-siren",
		"audio-16khz-32kbitrate-mono-mp3",
		"audio-16khz-64kbitrate-mono-mp3",
		"audio-16khz-128kbitrate-mono-mp3",
		"audio-24khz-48kbitrate-mono-mp3",
		"audio-24khz-96kbitrate-mono-mp3",
	}[a]
}

func StringToAudioOutput(s string) (AudioOutput, error) {
	audioMap := map[string]AudioOutput{
		"riff-8khz-8bit-mono-mulaw":        AudioRIFF8Bit8kHzMonoPCM,
		"riff-16khz-16bit-mono-pcm":        AudioRIFF16Bit16kHzMonoPCM,
		"riff-16khz-16kbps-mono-siren":     AudioRIFF16khz16kbpsMonoSiren,
		"riff-24khz-16bit-mono-pcm":        AudioRIFF24khz16bitMonoPcm,
		"raw-8khz-8bit-mono-mulaw":         AudioRAW8Bit8kHzMonoMulaw,
		"raw-16khz-16bit-mono-pcm":         AudioRAW16Bit16kHzMonoPcm,
		"raw-24khz-16bit-mono-pcm":         AudioRAW24khz16bitMonoPcm,
		"raw-22050hz-16bit-mono-pcm":       AudioRAW22050hz16bitMonoPcm,
		"ssml-16khz-16bit-mono-tts":        AudioSsml16khz16bitMonoTts,
		"audio-16khz-16kbps-mono-siren":    Audio16khz16kbpsMonoSiren,
		"audio-16khz-32kbitrate-mono-mp3":  Audio16khz32kbitrateMonoMp3,
		"audio-16khz-64kbitrate-mono-mp3":  Audio6khz64kbitrateMonoMp3,
		"audio-16khz-128kbitrate-mono-mp3": Audio16khz128kbitrateMonoMp3,
		"audio-24khz-48kbitrate-mono-mp3":  Audio24khz48kbitrateMonoMp3,
		"audio-24khz-96kbitrate-mono-mp3":  Audio24khz96kbitrateMonoMp3,
	}

	if audioOutput, exists := audioMap[s]; exists {
		return audioOutput, nil
	}

	return 0, fmt.Errorf("invalid audio output string: %s", s)
}

// Gender type for the digitized language
//
//go:generate enumer -type=Gender -linecomment -json
type Gender int

const (
	// GenderMale , GenderFemale are the static Gender constants for digitized voices.
	// See Gender in https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/language-support#standard-voices for breakdown.
	GenderMale   Gender = iota // Male
	GenderFemale               // Female
)

// Locale references the language or locale for text-to-speech.
// See "locale" in https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/language-support#standard-voices
//
//go:generate enumer -type=Locale -linecomment -json
type Locale int

const (
	LocaleArEG  Locale = 1  // ar-EG
	LocaleArSA  Locale = 2  // ar-SA
	LocaleBgBG  Locale = 3  // bg-BG
	LocaleCaES  Locale = 4  // ca-ES
	LocaleCsCZ  Locale = 5  // cs-CZ
	LocaleDaDK  Locale = 6  // da-DK
	LocaleDeAT  Locale = 7  // de-AT
	LocaleDeCH  Locale = 8  // de-CH
	LocaleDeDE  Locale = 9  // de-DE
	LocaleElGR  Locale = 10 // el-GR
	LocaleEnAU  Locale = 11 // en-AU
	LocaleEnCA  Locale = 12 // en-CA
	LocaleEnGB  Locale = 13 // en-GB
	LocaleEnIE  Locale = 14 // en-IE
	LocaleEnIN  Locale = 15 // en-IN
	LocaleEnUS  Locale = 16 // en-US
	LocaleEsES  Locale = 17 // es-ES
	LocaleEsMX  Locale = 18 // es-MX
	LocaleEtEE  Locale = 19 // et-EE
	LocaleFiFI  Locale = 20 // fi-FI
	LocaleFrCA  Locale = 21 // fr-CA
	LocaleFrCH  Locale = 22 // fr-CH
	LocaleFrFR  Locale = 23 // fr-FR
	LocaleGaIE  Locale = 24 // ga-IE
	LocaleHeIL  Locale = 25 // he-IL
	LocaleHiIN  Locale = 26 // hi-IN
	LocaleHrHR  Locale = 27 // hr-HR
	LocaleHuHU  Locale = 28 // hu-HU
	LocaleIDID  Locale = 29 // id-ID
	LocaleItIT  Locale = 30 // it-IT
	LocaleJaJP  Locale = 31 // ja-JP
	LocaleKoKR  Locale = 32 // ko-KR
	LocaleLtLT  Locale = 33 // lt-LT
	LocaleLvLV  Locale = 34 // lv-LV
	LocaleMtMT  Locale = 35 // mt-MT
	LocaleMrIN  Locale = 36 // mr-IN
	LocaleMsMY  Locale = 37 // ms-MY
	LocaleNbNO  Locale = 38 // nb-NO
	LocaleNlNL  Locale = 39 // nl-NL
	LocalePlPL  Locale = 40 // pl-PL
	LocalePtBR  Locale = 41 // pt-BR
	LocalePtPT  Locale = 42 // pt-PT
	LocaleRoRO  Locale = 43 // ro-RO
	LocaleRuRU  Locale = 44 // ru-RU
	LocaleSkSK  Locale = 45 // sk-SK
	LocaleSlSI  Locale = 46 // sl-SI
	LocaleSvSE  Locale = 47 // sv-SE
	LocaleTaIN  Locale = 48 // ta-IN
	LocaleTeIN  Locale = 49 // te-IN
	LocaleThTH  Locale = 50 // th-TH
	LocaleTrTR  Locale = 51 // tr-TR
	LocaleViVN  Locale = 52 // vi-VN
	LocaleZhCN  Locale = 53 // zh-CN
	LocaleZhHK  Locale = 54 // zh-HK
	LocaleZhTW  Locale = 55 // zh-TW
	LocaleFilPH Locale = 56 // fil-PH
	LocaleTaMY  Locale = 57 // ta-MY
)

// Region references the locations of the availability of standard voices.
// See https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/regions#standard-voices
type Region int

const (
	// Azure regions and their endpoints that support the Text To Speech service.
	RegionAustraliaEast Region = iota
	RegionBrazilSouth
	RegionCanadaCentral
	RegionCentralUS
	RegionEastAsia
	RegionEastUS
	RegionEastUS2
	RegionFranceCentral
	RegionIndiaCentral
	RegionJapanEast
	RegionJapanWest
	RegionKoreaCentral
	RegionNorthCentralUS
	RegionNorthEurope
	RegionSouthCentralUS
	RegionSoutheastAsia
	RegionUKSouth
	RegionWestEurope
	RegionWestUS
	RegionWestUS2
)

func (t Region) String() string {
	return [...]string{
		"australiaeast",
		"brazilsouth",
		"canadacentral",
		"centralus",
		"eastasia",
		"eastus",
		"eastus2",
		"francecentral",
		"indiacentral",
		"japaneast",
		"japanwest",
		"koreacentral",
		"northcentralus",
		"northeurope",
		"southcentralus",
		"southeastasia",
		"uksouth",
		"westeurope",
		"westus",
		"westus2",
	}[t]
}
