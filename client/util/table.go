// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package util

import (
	"fmt"

	"github.com/apcera/termtables"
)

// Table is the interface of table
type Table interface {
	SetHeader(headers ...interface{})
	AddRow(items ...interface{})
	Render()
	ToString() string
}

// TermTable is the wrapper of termtables and
type TermTable struct {
	table *termtables.Table
}

// NewTable is the constructor of Table
func NewTable() Table {
	return &TermTable{
		table: termtables.CreateTable(),
	}
}

// SetHeader is the implementation of interface
func (t *TermTable) SetHeader(headers ...interface{}) {
	t.table.AddHeaders(headers...)
}

// AddRow is the implementation of interface
func (t *TermTable) AddRow(items ...interface{}) {
	t.table.AddRow(items...)
}

// Render is the implementation of interface
func (t *TermTable) Render() {
	fmt.Print(t.table.Render())
}

// ToString is the implementation of interface
func (t *TermTable) ToString() string {
	return t.table.Render()
}
