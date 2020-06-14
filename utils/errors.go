package utils

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func FatalOnError(err error) {
	if err != nil {
		logrus.WithError(err).
			WithField("stack", debug.Stack()).
			Fatal()
	}
}

func PanicOnError(err error) {
	if err != nil {
		logrus.WithError(err).
			WithField("stack", debug.Stack()).
			Panic()
	}
}

func LogOnError(err error) {
	if err != nil {
		logrus.WithError(err).Error()
	}
}

//func WarnOnError(err error) {
//	if err != nil {
//		logrus.WithError(err).Warn()
//	}
//}
//
//func InfoOnError(err error) {
//	if err != nil {
//		logrus.WithError(err).Info()
//	}
//}
