// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:generate goqueryset -in wallet.go

package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Wallet qiwi wallet credentials and setting
// gen:qs
type Wallet struct {
	gorm.Model

	Name string
	// Phone number
	WalletID    uint64
	Blocked     bool
	Token       string `gorm:"unique_index"`
	TokenExpiry time.Time

	Balance float64
	Limit   uint

	TotalMonthIncoming float64
	TotalMonthOutgoing float64
	TotalSynced        time.Time

	// userID
	OwnerID uint

	GroupID uint
}

func (w Wallet) String() string {
	if w.Name == "" {
		return fmt.Sprintf("+%d", w.WalletID)
	}
	return w.Name
}

// GroupWallets returns all group wallet
func GroupWallets(groupID uint) (res []Wallet, err error) {
	err = NewWalletQuerySet(db).GroupIDEq(groupID).All(&res)
	return
}

// GetAllWallets returns all wallets. Used for synchronizer
func GetAllWallets() (res []Wallet, err error) {
	err = NewWalletQuerySet(db).All(&res)
	return
}

// GetWallet returns wallet by their ID
func GetWallet(walletID uint) (wallet *Wallet, err error) {
	wallet = new(Wallet)
	err = NewWalletQuerySet(db).IDEq(walletID).One(wallet)
	return
}

// CreateWallet create an wallet
func CreateWallet(wallet *Wallet) (err error) {
	err = wallet.Create(db)
	if err != nil {
		return
	}

	return
}
