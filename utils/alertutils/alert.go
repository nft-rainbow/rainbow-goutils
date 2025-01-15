package alertutils

import (
	"context"
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/alert"
	"github.com/nft-rainbow/rainbow-goutils/utils/commonutils"
	"github.com/sirupsen/logrus"
)

type DingHelper struct {
	channel *alert.DingTalkChannel
}

var defaultDingHelper *DingHelper

func DingInfof(pattern string, args ...any) error {
	return defaultDingHelper.DingInfof(pattern, args...)
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

func (d *DingHelper) DingInfof(pattern string, args ...any) error {
	logrus.WithField("args", args).Info(pattern)
	return d.DingText(alert.SeverityLow, "", fmt.Sprintf(pattern, args...))
}

func (d *DingHelper) DingWarnf(pattern string, args ...any) error {
	logrus.WithField("args", args).Info(pattern)
	return d.DingText(alert.SeverityMedium, "", fmt.Sprintf(pattern, args...))
}

func (d *DingHelper) DingError(err error, describe ...string) error {
	if len(describe) == 0 {
		describe = append(describe, "unexpected error")
	}
	logrus.WithError(err).Error(describe)
	return d.DingText(alert.SeverityHigh, describe[0], fmt.Sprintf("%+v", err))
}

func (d *DingHelper) DingPanicf(err error, description string, args ...any) {
	logrus.WithField("args", args).WithError(err).Error(description)
	err = d.DingText(alert.SeverityCritical, fmt.Sprintf(description, args...), fmt.Sprintf("%+v", err))
	panic(err)
}

func (d *DingHelper) DingText(level alert.Severity, brief, detail string) error {
	if d == nil {
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

	return d.channel.Send(context.Background(), &alert.Notification{
		Title:    title,
		Content:  content,
		Severity: alert.Severity(level),
	})
}

// Initialize alert configuration using Viper.
// Channels with name set to default and platform set to dingtalk are considered the default channel.
// If no channel with name set to default exists, no notifications will be sent when call DingXXX methods, and no errors will be raised.
func MustInitFromViper() {
	alert.MustInitFromViper()
	dh, err := GetDingHelper("default")
	if err != nil {
		return
	}
	defaultDingHelper = dh
}

func GetDingHelper(name string) (*DingHelper, error) {
	ch, ok := alert.DefaultManager().Channel(name)
	if !ok {
		return nil, fmt.Errorf("ding channel not found")
	}
	if ch.Type() != alert.ChannelTypeDingTalk {
		return nil, fmt.Errorf("channel is not dingtalk channel")
	}
	return &DingHelper{channel: ch.(*alert.DingTalkChannel)}, nil
}

func MustGetDingHelper(name string) *DingHelper {
	return commonutils.Must(GetDingHelper(name))
}
