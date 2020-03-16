// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package eraft

import (
	"madledger/common/util"
)

// Here defines some status
const (
	Stopped int32 = iota
	OnStarting
	Running
)

// randNode random return a value of a cluster
func randNode(cluster map[uint64]string) string {
	var i = util.RandNum(len(cluster))
	for id := range cluster {
		if i == 0 {
			return cluster[id]
		}
		i--
	}
	return ""
}

// getLeaderFromError try to parse leader address from error
// func getLeaderFromError(err error) string {
// 	e := strings.Replace(err.Error(), "rpc error: code = Unknown desc =", "", -1)
// 	if strings.Contains(e, "Leader is") {
// 		return strings.Replace(strings.Replace(e, "Leader is", "", 1), " ", "", -1)
// 	}

// 	if strings.Contains(e, "Please send tx to") {
// 		return strings.Replace(strings.Replace(e, "Please send tx to", "", 1), " ", "", -1)
// 	}
// 	return ""
// }
