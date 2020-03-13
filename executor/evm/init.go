package evm

import "github.com/sirupsen/logrus"

var (
	log = logrus.WithFields(logrus.Fields{"app": "github.com/thu-arxan/evm", "package": "executor/evm"})
)
