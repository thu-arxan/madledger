package util

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

var(
	glog = logrus.WithFields(logrus.Fields{"app": "grpc", "package": "grpc"})
	glogFlag = false
)
// Wrapper to provide log for GRPC
type GrpcLogger struct {
	*logrus.Entry
}

func (_ *GrpcLogger) V(_ int) bool {
	return false
}

func MountLogger() {
	if !glogFlag { // Export GRPC's log, execute only one time.
		glog.Infof("Mount GRPC logger...")
		grpclog.SetLoggerV2(&GrpcLogger{Entry: glog})
		glogFlag = true
	}
}
