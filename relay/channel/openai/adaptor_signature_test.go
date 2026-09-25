package openai

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomChannelSignatureHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{
		ChannelType: constant.ChannelTypeCustom,
		ApiKey:      "1234567890abcdef",
		ChannelOtherSettings: dto.ChannelOtherSettings{
			SignatureType:  "cmc_sh",
			SignatureAppID: "app_1",
		},
	}}
	header := http.Header{}
	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	assert.True(t, strings.HasPrefix(header.Get("Authorization"), "Bearer app_1."))
}

func TestCustomChannelSignatureHonorsAuthorizationOverride(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{
		ChannelType:          constant.ChannelTypeCustom,
		ApiKey:               "1234567890abcdef",
		HeadersOverride:      map[string]any{"authorization": "Bearer override"},
		ChannelOtherSettings: dto.ChannelOtherSettings{SignatureType: "cmc_sh", SignatureAppID: "app_1"},
	}}
	header := http.Header{}
	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	assert.Empty(t, header.Get("Authorization"))
}

func TestCustomChannelSignatureRealtimeSubprotocol(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/realtime", nil)
	c.Request.Header.Set("Sec-WebSocket-Protocol", "realtime")
	info := &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeRealtime,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType: constant.ChannelTypeCustom,
			ApiKey:      "1234567890abcdef",
			ChannelOtherSettings: dto.ChannelOtherSettings{
				SignatureType:  "cmc_sh",
				SignatureAppID: "app_1",
			},
		},
	}
	header := http.Header{}
	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	protocol := header.Get("Sec-WebSocket-Protocol")
	assert.Contains(t, protocol, "openai-insecure-api-key.app_1.")
	assert.NotContains(t, protocol, "1234567890abcdef")
}
