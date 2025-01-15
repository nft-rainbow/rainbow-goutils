package alertutils

import (
	"context"
	"testing"

	"github.com/Conflux-Chain/go-conflux-util/alert"
	"github.com/nft-rainbow/rainbow-goutils/utils/commonutils"
	"github.com/stretchr/testify/assert"
)

func TestRawDingTalkChannel(t *testing.T) {
	err := alert.NewDingTalkChannel("1", commonutils.Must(alert.NewDingtalkMarkdownFormatter([]string{"Alert"}, []string{"17611422948"})), alert.DingTalkConfig{
		Webhook:   "xxx",
		AtMobiles: []string{"YOUR_PHONE"},
		MsgType:   "markdown",
	}).Send(context.Background(), &alert.Notification{
		Title:   "test",
		Content: "test",
	})
	assert.NoError(t, err)
}

func TestDingInfof(t *testing.T) {
	ch := alert.NewDingTalkChannel("1", commonutils.Must(alert.NewDingtalkMarkdownFormatter([]string{"Alert"}, []string{"17611422948"})), alert.DingTalkConfig{
		Webhook:   "YOUR_WEBHOOK",
		AtMobiles: []string{"YOUR_PHONE"},
		MsgType:   "markdown",
	})
	alert.DefaultManager().Add(ch)
	defaultDingHelper = &DingHelper{channel: ch}
	DingInfof("test %v", "test")
}
