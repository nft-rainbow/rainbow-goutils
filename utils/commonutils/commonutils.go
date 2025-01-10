package commonutils

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func Retry(count int, interval time.Duration, description string, fn func() error) error {
	for i := 0; i < count; i++ {
		err := fn()
		if err != nil {
			if i == count-1 {
				logrus.WithError(err).WithField("description", description).WithField("stack", string(debug.Stack())).WithField("retry cnt", i).Info("retry function error")
				return errors.WithMessage(err, fmt.Sprintf("retry '%s' %d times and failed", description, count))
			} else {
				time.Sleep(interval)
				continue
			}
		}
		return nil
	}
	return nil
}

// 重试直到
// 1. 遇到预期错误 (使用errors.Is判断)
// 2. 重试次数达到上限
func RetryOrExpectError(count int, interval time.Duration, description string, expectErr error, fn func() error) error {
	for i := 0; i < count; i++ {
		err := fn()
		if errors.Is(err, expectErr) {
			return err
		}
		if err != nil {
			if i == count-1 {
				logrus.WithError(err).WithField("description", description).WithField("stack", string(debug.Stack())).WithField("retry cnt", i).Info("retry function error")
				return errors.WithMessage(err, fmt.Sprintf("retry '%s' %d times and failed", description, count))
			} else {
				time.Sleep(interval)
				continue
			}
		}
		return nil
	}
	return nil
}

func MapSlice[T any, R any](items []T, mapFunc func(item T) (R, error)) ([]R, error) {
	var result []R
	for _, item := range items {
		r, err := mapFunc(item)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func MustMapSlice[T any, R any](items []T, mapFunc func(item T) R) []R {
	return lo.Map(items, func(v T, i int) R {
		return mapFunc(v)
	})
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
