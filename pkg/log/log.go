// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

// Info log simple text (fmt)
func Info(text string, args ...interface{}) {
	logger.Infof(text, args...)
}

// Trace log error with message
func Trace(err error, args ...interface{}) {
	logger.WithError(err).Debug(args...)
}

// Warn error with message
func Warn(err error, args ...interface{}) {
	logger.WithError(err).Warn(args...)
}
