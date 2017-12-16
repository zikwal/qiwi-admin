// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"path/filepath"

	"github.com/zhuharev/qiwi-admin/pkg/setting"

	"github.com/asdine/storm"
)

var (
	stormDB *storm.DB
)

func NewStormContext() (err error) {
	stormDB, err = storm.Open(filepath.Join(setting.App.DataDir, "storm.bolt"))
	if err != nil {
		return
	}
	return
}
