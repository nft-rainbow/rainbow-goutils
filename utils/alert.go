package utils

import (
	"context"
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/alert"
	"github.com/sirupsen/logrus"
)

var defaultChannel *alert.DingTalkChannel
var financeChannel *alert.DingTalkChannel
var enabled bool

func SetDingEnabled(e bool) {
	enabled = e
}

func MustInitDefaultDingChannel(config alert.DingTalkConfig, tags []string) *alert.DingTalkChannel {
	defaultChannel = mustInitDingChannel("default", config, tags)
	return defaultChannel
}

func MustInitFinanceDingChannel(config alert.DingTalkConfig, tags []string) *alert.DingTalkChannel {
	financeChannel = mustInitDingChannel("finance", config, tags)
	return financeChannel
}

func mustInitDingChannel(id string, config alert.DingTalkConfig, tags []string) *alert.DingTalkChannel {
	fmt, err := alert.NewDingtalkMarkdownFormatter(tags)
	if err != nil {
		panic(err)
	}
	financeChannel = alert.NewDingTalkChannel(id, fmt, config)
	return financeChannel
}

func DingInfo(pattern string, args ...any) {
	logrus.WithField("args", args).Info(pattern)
	DingText(alert.SeverityLow, "info", fmt.Sprintf(pattern, args...))
}

func DingWarnf(pattern string, args ...any) {
	logrus.WithField("args", args).Warn(pattern)
	DingText(alert.SeverityMedium, "warn", fmt.Sprintf(pattern, args...))
}

func DingFinanceWarnf(pattern string, args ...any) {
	logrus.WithField("args", args).Warn(pattern)
	DingTextToFinance(alert.SeverityMedium, "warn", fmt.Sprintf(pattern, args...))
}

func DingError(err error, describe ...string) {
	if len(describe) == 0 {
		describe = append(describe, "unexpected error")
	}
	logrus.WithError(err).Error(describe)
	DingText(alert.SeverityCritical, describe[0], err.Error())
}

func DingPanicf(err error, description string, args ...any) {
	logrus.WithField("args", args).WithError(err).Error(description)
	DingText(alert.SeverityCritical, description, fmt.Sprintf("%+v", err))
	panic(err)
}

func DingText(level alert.Severity, brief, detail string) {
	DingTextToChannel(defaultChannel, level, brief, detail)
}

func DingTextToFinance(level alert.Severity, brief, detail string) {
	DingTextToChannel(financeChannel, level, brief, detail)
}

func DingTextToChannel(channel *alert.DingTalkChannel, level alert.Severity, brief, detail string) {
	if !enabled {
		return
	}
	// if at mobiles is not empty, add to detail
	if len(channel.Config.AtMobiles) > 0 {
		for _, mobile := range channel.Config.AtMobiles {
			detail += fmt.Sprintf("\n@%v", mobile)
		}
	}

	if err := channel.Send(context.Background(), &alert.Notification{
		Title:    brief,
		Content:  detail,
		Severity: level,
	}); err != nil {
		logrus.WithError(err).WithField("detail", detail).Error("failed to send dingding")
	}
}
