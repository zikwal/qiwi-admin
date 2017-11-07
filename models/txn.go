// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:generate goqueryset -in txn.go

package models

import "time"

// TxnType of txn
type TxnType uint

const (
	// In incominc txn
	In TxnType = iota + 1
	// Out outgoing txn
	Out
	// QiwiCard payment from Qiwi card
	QiwiCard
)

// Txn qiwi transaction
// gen:qs
type Txn struct {
	ID uint `gorm:"unique_index"`

	TxnType    TxnType
	ProviderID uint // ?
	Amount     float64
	CreatedAt  time.Time `gorm:"index"`
	Fee        float64
	Status     Status

	WalletID uint
}

// Status represent status of txn
type Status uint

const (
	// Waiting txn created but not processed
	Waiting Status = iota + 1
	// Success txn
	Success
	// Error represent txn with an error
	Error
)

// CreateMultipleTxns insert transactions in on txn
func CreateMultipleTxns(walletID uint, txns []Txn) (err error) {
	tx := db.Begin()
	for _, txn := range txns {
		txn.WalletID = walletID
		err = tx.Create(txn).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	return
}

// GetWalletTxns get lasts wallet txns
func GetWalletTxns(walletID uint) (res []Txn, err error) {
	err = NewTxnQuerySet(db).WalletIDEq(walletID).OrderDescByID().Limit(50).All(&res)
	return
}

// GetLastTxn return last txn on wallet with walletID
func GetLastTxn(walletID uint) (txnID uint, err error) {
	txn := new(Txn)
	err = NewTxnQuerySet(db).WalletIDEq(walletID).OrderDescByID().One(txn)
	if err != nil {
		return
	}
	txnID = txn.ID
	return
}
