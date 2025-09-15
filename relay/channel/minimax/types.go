package minimax

// MiniMax TTS Request structures
type MiniMaxTTSRequest struct {
	Model        string        `json:"model"`
	Text         string        `json:"text"`
	VoiceSetting VoiceSetting  `json:"voice_setting"`
	AudioSetting *AudioSetting `json:"audio_setting,omitempty"`
}

type VoiceSetting struct {
	VoiceID string  `json:"voice_id"`
	Speed   float64 `json:"speed,omitempty"`
	Vol     float64 `json:"vol,omitempty"`
	Pitch   int     `json:"pitch,omitempty"`
}

type AudioSetting struct {
	SampleRate int64  `json:"sample_rate,omitempty"`
	Bitrate    int64  `json:"bitrate,omitempty"`
	Format     string `json:"format"`
	Channel    int64  `json:"channel,omitempty"`
}

// MiniMax TTS Response structures
type MiniMaxTTSResponse struct {
	BaseResp  BaseResp  `json:"base_resp"`
	Data      Data      `json:"data"`
	ExtraInfo ExtraInfo `json:"extra_info"`
}

type BaseResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type Data struct {
	Audio  string `json:"audio"`  // hex-encoded audio data
	Status int    `json:"status"`
}

type ExtraInfo struct {
	UsageCharacters int    `json:"usage_characters"`
	AudioFormat     string `json:"audio_format"`
}