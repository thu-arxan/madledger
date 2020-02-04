package raft

import "github.com/sirupsen/logrus"

var (
	log     = logrus.WithFields(logrus.Fields{"app": "consensus", "package": "raft"})
	etcdLog = logrus.WithFields(logrus.Fields{"app": "consensus", "package": "etcd/raft"})
)
