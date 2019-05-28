package raft

import "github.com/sirupsen/logrus"

var (
	log     = logrus.WithFields(logrus.Fields{"app": "blockchain", "package": "raft"})
	etcdLog = logrus.WithFields(logrus.Fields{"app": "blockchain", "package": "etcd/raft"})
)
