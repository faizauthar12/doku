package helper

import (
	"github.com/sirupsen/logrus"
)

func Recover(location string) {
	if r := recover(); r != nil {
		logrus.Debugf("recover panic action from %s : %s\n", location, r)
	}
}
