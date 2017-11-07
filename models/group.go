// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:generate goqueryset -in group.go

package models

import "github.com/jinzhu/gorm"

// Group represent group of wallets
// gen:qs
type Group struct {
	gorm.Model

	Name    string
	OwnerID uint
}

// CreateGroup save new group in db
func CreateGroup(name string, ownerID uint) (group *Group, err error) {
	group = new(Group)
	group.Name = name
	group.OwnerID = ownerID

	err = group.Create(db)

	return
}

// GetUserGroups return groups where user has access
func GetUserGroups(userID uint) (res []Group, err error) {
	err = NewGroupQuerySet(db).OwnerIDEq(userID).All(&res)
	return
}

// GetGroup returns group by id
func GetGroup(id uint) (group *Group, err error) {
	group = new(Group)
	err = NewGroupQuerySet(db).IDEq(id).One(group)
	return
}
