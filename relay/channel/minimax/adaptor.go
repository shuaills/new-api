package minimax

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"one-api/dto"
	"one-api/relay/channel"
	relaycommon "one-api/relay/common"
	relayconstant "one-api/relay/constant"
	"one-api/types"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Adaptor struct {
	ChannelType int
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo, request dto.GeneralOpenAIRequest) {
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	return GetRequestURL(info)
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	switch info.RelayMode {
	case relayconstant.RelayModeAudioSpeech:
		return a.ConvertTTSRequest(c, info, request)
	}

	return request, nil
}

func (a *Adaptor) ConvertTTSRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	// Parse AudioRequest from the request
	audioReq := &dto.AudioRequest{}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	err = json.Unmarshal(jsonData, audioReq)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling audio request: %w", err)
	}

	// Convert OpenAI format to MiniMax format
	minimaxReq := &MiniMaxTTSRequest{
		Model: audioReq.Model,
		Text:  audioReq.Input,
		VoiceSetting: VoiceSetting{
			VoiceID: a.convertVoice(audioReq.Voice),
			Speed:   audioReq.Speed,
		},
	}

	// Parse response format if provided (e.g., "mp3", "mp3-1-32000-128000")
	if audioReq.ResponseFormat != "" {
		minimaxReq.AudioSetting = a.parseAudioFormat(audioReq.ResponseFormat)
	} else {
		// Default audio settings
		minimaxReq.AudioSetting = &AudioSetting{
			Format:     "mp3",
			Channel:    1,
			SampleRate: 32000,
			Bitrate:    128000,
		}
	}

	return minimaxReq, nil
}

// convertVoice maps OpenAI voice names to MiniMax voice IDs
func (a *Adaptor) convertVoice(voice string) string {
	voiceMap := map[string]string{
		"alloy":   "female-chengshu",
		"echo":    "male-qn-qingse",
		"fable":   "male-qn-jingying",
		"onyx":    "presenter_male",
		"nova":    "presenter_female",
		"shimmer": "audiobook_female_1",
	}

	if mappedVoice, exists := voiceMap[voice]; exists {
		return mappedVoice
	}
	// Default to first voice if not found
	return "female-chengshu"
}

// parseAudioFormat parses response_format like "mp3-1-32000-128000"
func (a *Adaptor) parseAudioFormat(format string) *AudioSetting {
	parts := strings.Split(format, "-")

	setting := &AudioSetting{
		Format:     "mp3",
		Channel:    1,
		SampleRate: 32000,
		Bitrate:    128000,
	}

	if len(parts) >= 1 {
		setting.Format = parts[0]
	}
	if len(parts) >= 2 {
		if channel, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			setting.Channel = channel
		}
	}
	if len(parts) >= 3 {
		if sampleRate, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
			setting.SampleRate = sampleRate
		}
	}
	if len(parts) >= 4 {
		if bitrate, err := strconv.ParseInt(parts[3], 10, 64); err == nil {
			setting.Bitrate = bitrate
		}
	}

	return setting
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	switch info.RelayMode {
	case relayconstant.RelayModeAudioSpeech:
		// Convert to MiniMax TTS format
		minimaxReq := &MiniMaxTTSRequest{
			Model: request.Model,
			Text:  request.Input,
			VoiceSetting: VoiceSetting{
				VoiceID: a.convertVoice(request.Voice),
				Speed:   request.Speed,
			},
		}

		// Parse response format
		if request.ResponseFormat != "" {
			minimaxReq.AudioSetting = a.parseAudioFormat(request.ResponseFormat)
		} else {
			minimaxReq.AudioSetting = &AudioSetting{
				Format:     "mp3",
				Channel:    1,
				SampleRate: 32000,
				Bitrate:    128000,
			}
		}

		jsonData, err := json.Marshal(minimaxReq)
		if err != nil {
			return nil, fmt.Errorf("error marshalling MiniMax TTS request: %w", err)
		}
		return bytes.NewReader(jsonData), nil
	}

	return nil, fmt.Errorf("unsupported relay mode for MiniMax audio: %d", info.RelayMode)
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error) {
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage *dto.Usage, err *types.OpenAIErrorWithStatusCode) {
	switch info.RelayMode {
	case relayconstant.RelayModeAudioSpeech:
		return a.TTSDoResponse(c, resp, info)
	}

	return nil, &types.OpenAIErrorWithStatusCode{
		StatusCode: http.StatusInternalServerError,
		Error: types.OpenAIError{
			Message: "unsupported relay mode for MiniMax",
			Type:    "minimax_error",
		},
	}
}

func (a *Adaptor) TTSDoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage *dto.Usage, err *types.OpenAIErrorWithStatusCode) {
	defer resp.Body.Close()

	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, &types.OpenAIErrorWithStatusCode{
			StatusCode: http.StatusInternalServerError,
			Error: types.OpenAIError{
				Message: "error reading response body",
				Type:    "minimax_error",
			},
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &types.OpenAIErrorWithStatusCode{
			StatusCode: resp.StatusCode,
			Error: types.OpenAIError{
				Message: string(body),
				Type:    "minimax_error",
			},
		}
	}

	// Parse MiniMax response
	var minimaxResp MiniMaxTTSResponse
	err2 = json.Unmarshal(body, &minimaxResp)
	if err2 != nil {
		return nil, &types.OpenAIErrorWithStatusCode{
			StatusCode: http.StatusInternalServerError,
			Error: types.OpenAIError{
				Message: "error parsing MiniMax response",
				Type:    "minimax_error",
			},
		}
	}

	// Check for API errors
	if minimaxResp.BaseResp.StatusCode != 0 {
		return nil, &types.OpenAIErrorWithStatusCode{
			StatusCode: http.StatusBadRequest,
			Error: types.OpenAIError{
				Message: minimaxResp.BaseResp.StatusMsg,
				Type:    "minimax_error",
			},
		}
	}

	// Decode hex audio data
	audioBytes, err2 := hex.DecodeString(minimaxResp.Data.Audio)
	if err2 != nil {
		return nil, &types.OpenAIErrorWithStatusCode{
			StatusCode: http.StatusInternalServerError,
			Error: types.OpenAIError{
				Message: "error decoding audio data",
				Type:    "minimax_error",
			},
		}
	}

	// Set response headers and return audio data
	c.Header("Content-Type", "audio/"+minimaxResp.ExtraInfo.AudioFormat)
	c.Header("Content-Length", strconv.Itoa(len(audioBytes)))
	c.Data(http.StatusOK, "audio/"+minimaxResp.ExtraInfo.AudioFormat, audioBytes)

	// Return usage information
	usage = &dto.Usage{
		PromptTokens:     minimaxResp.ExtraInfo.UsageCharacters,
		CompletionTokens: 0,
		TotalTokens:      minimaxResp.ExtraInfo.UsageCharacters,
	}

	return usage, nil
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}