package evm

import "github.com/sirupsen/logrus"

var (
	log = logrus.WithFields(logrus.Fields{"app": "evm", "package": "executor/evm"})
)
