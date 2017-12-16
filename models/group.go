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

	AutoTransferObjectID  uint
	AutTransferObjectType ObjectType

	Counters GroupCounters `gorm:"-"`
}

type GroupWithCounters struct {
	Group
	Balance float64
	Count   int
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

// GetUserGroupsWithCounters return groups where user has access
func GetUserGroupsWithCounters(userID uint) (res []GroupWithCounters, err error) {
	sql := `SELECT groups.*, balance, count
	FROM groups
	LEFT JOIN (
		SELECT
		  SUM(balance) as balance,
			count() as count,
			group_id
		FROM wallets
		WHERE "wallets"."deleted_at" IS NULL
		)
	ON group_id = id
	WHERE owner_id = ?
	  AND "groups"."deleted_at" IS NULL`
	err = db.Raw(sql, userID).Scan(&res).Error
	return
}

// GetGroup returns group by id
func GetGroup(id uint, userIDs ...uint) (group *Group, err error) {
	group = new(Group)
	query := NewGroupQuerySet(db).IDEq(id)
	if len(userIDs) > 0 {
		query = query.OwnerIDEq(userIDs[0])
	}
	err = query.One(group)
	if err != nil {
		return
	}

	counters, err := GetGroupCounters(id)
	group.Counters = counters

	return
}

// GroupCounters response of aggregate sql request
type GroupCounters struct {
	Balance float64 `gorm:"balance"`
	Count   int     `gorm:"count"`
}

// GetGroupCounters agregate stat
func GetGroupCounters(groupID uint) (res GroupCounters, err error) {
	sql := `select sum(wallets.balance) as balance,count() as count from wallets where "wallets"."deleted_at" IS NULL AND group_id = ?`
	err = db.Raw(sql, groupID).Scan(&res).Error
	return
}

// GetGroupFreeWallet returns wallet, which the cat receive amount size payment
func GetGroupFreeWallet(groupID uint, amount uint) (wallet *Wallet, err error) {
	wallet = new(Wallet)
	err = NewWalletQuerySet(db).GroupIDEq(groupID).BalanceLt(15000 - float64(amount)).One(wallet)
	return
}

// DeleteGroup remove group from db
func DeleteGroup(groupID uint) error {
	gr := new(Group)
	gr.ID = groupID
	return db.Delete(gr).Error
}
