package minimax

import (
	"fmt"
	relaycommon "one-api/relay/common"
	relayconstant "one-api/relay/constant"
)

func GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	switch info.RelayMode {
	case relayconstant.RelayModeAudioSpeech:
		return fmt.Sprintf("%s/v1/t2a_v2", info.ChannelBaseUrl), nil
	case relayconstant.RelayModeAudioTranscription:
		return fmt.Sprintf("%s/v1/audio/transcription", info.ChannelBaseUrl), nil
	case relayconstant.RelayModeAudioTranslation:
		return fmt.Sprintf("%s/v1/audio/translation", info.ChannelBaseUrl), nil
	default:
		return fmt.Sprintf("%s/v1/text/chatcompletion_v2", info.ChannelBaseUrl), nil
	}
}
