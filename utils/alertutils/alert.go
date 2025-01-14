package alertutils

import (
	"context"
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/alert"
	"github.com/nft-rainbow/rainbow-goutils/utils/commonutils"
	"github.com/sirupsen/logrus"
)

type dingHelper struct {
	channel *alert.DingTalkChannel
}

var defaultDingHelper *dingHelper

func DingInfo(pattern string, args ...any) error {
	return defaultDingHelper.DingInfo(pattern, args...)
}

func DingWarn(pattern string, args ...any) error {
	return defaultDingHelper.DingWarn(pattern, args...)
}

func DingWarnf(pattern string, args ...any) error {
	return defaultDingHelper.DingWarnf(pattern, args...)
}

func DingError(err error, describe ...string) error {
	return defaultDingHelper.DingError(err, describe...)
}

func DingPanicf(err error, description string, args ...any) {
	defaultDingHelper.DingPanicf(err, description, args...)
}

func (d *dingHelper) DingInfo(pattern string, args ...any) error {
	logrus.WithField("args", args).Info(pattern)
	return d.DingText(alert.SeverityLow, "", fmt.Sprintf(pattern, args...))
}

func (d *dingHelper) DingWarn(pattern string, args ...any) error {
	logrus.WithField("args", args).Info(pattern)
	return d.DingText(alert.SeverityLow, "", fmt.Sprintf(pattern, args...))
}

func (d *dingHelper) DingWarnf(pattern string, args ...any) error {
	logrus.WithField("args", args).Warn(pattern)
	return d.DingText(alert.SeverityMedium, "", fmt.Sprintf(pattern, args...))
}

func (d *dingHelper) DingError(err error, describe ...string) error {
	if len(describe) == 0 {
		describe = append(describe, "unexpected error")
	}
	logrus.WithError(err).Error(describe)
	return d.DingText(alert.SeverityHigh, describe[0], fmt.Sprintf("%+v", err))
}

func (d *dingHelper) DingPanicf(err error, description string, args ...any) {
	logrus.WithField("args", args).WithError(err).Error(description)
	err = d.DingText(alert.SeverityCritical, description, fmt.Sprintf("%+v", err))
	panic(err)
}

func (d *dingHelper) DingText(level alert.Severity, brief, detail string) error {
	if defaultDingHelper == nil {
		return nil
	}

	title := "ℹ️ Info"
	switch level {
	case alert.SeverityMedium:
		title = "⚠️ Warn"
	case alert.SeverityHigh:
		title = "💔 Error"
	case alert.SeverityCritical:
		title = "😱 Panic"
	}

	content := detail
	if brief != "" {
		content = fmt.Sprintf("%s: %s", brief, detail)
	}

	return defaultDingHelper.channel.Send(context.Background(), &alert.Notification{
		Title:    title,
		Content:  content,
		Severity: alert.Severity(level),
	})
}

func MustInitFromViper() {
	alert.MustInitFromViper()
	allCh := alert.DefaultManager().All("")
	for _, ch := range allCh {
		if ch.Type() == alert.ChannelTypeDingTalk {
			defaultDingHelper = &dingHelper{channel: ch.(*alert.DingTalkChannel)}
			return
		}
	}
}

func GetDingHelper(name string) (*dingHelper, error) {
	ch, ok := alert.DefaultManager().Channel(name)
	if !ok {
		return nil, fmt.Errorf("ding channel not found")
	}
	if ch.Type() != alert.ChannelTypeDingTalk {
		return nil, fmt.Errorf("channel is not dingtalk channel")
	}
	return &dingHelper{channel: ch.(*alert.DingTalkChannel)}, nil
}

func MustGetDingHelper(name string) *dingHelper {
	return commonutils.Must(GetDingHelper(name))
}
