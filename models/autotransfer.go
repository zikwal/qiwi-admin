// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:generate goqueryset -in autotransfer.go

package models

import "github.com/jinzhu/gorm"

// Autotransfer used for logging
// gen:qs
type Autotransfer struct {
	gorm.Model

	SourceTxnID uint
	TargetTxnID uint

	SourceID   uint
	SourceType ObjectType

	TargetType ObjectType
	TargetID   string

	Amount uint
}

// AutotransferSave save
func AutotransferSave(a *Autotransfer) (err error) {
	err = a.Create(db)
	return
}

// AutotransferLast returns last autotransfers
func AutotransferLast(groupID uint) (res []Autotransfer, err error) {
	err = NewAutotransferQuerySet(db).SourceIDEq(groupID).OrderDescByID().Limit(50).All(&res)
	if err != nil {
		return
	}
	return
}
