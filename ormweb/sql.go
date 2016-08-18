// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ormweb

import (
	"github.com/dotcoo/orm"
)

func NewSQL(table ...string) *orm.SQL {
	return orm.NewSQL(table...)
}
