package util

import "github.com/sirupsen/logrus"

// Wrapper to provide log for GRPC
type GrpcLogger struct {
	*logrus.Entry
}

func (_ *GrpcLogger) V(_ int) bool {
	return false
}