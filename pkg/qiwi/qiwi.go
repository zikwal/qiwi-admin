// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package qiwi

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/zhuharev/qiwi"
	"github.com/zhuharev/qiwi-admin/models"
)

// CheckToken simple token checker
func CheckToken(token string) (walletID uint64, blocked bool, balance float64, err error) {
	client := qiwi.New(token, qiwi.Debug)
	profile, err := client.Profile.Current()
	if err != nil {
		return
	}
	color.Green("%v", profile)
	walletID = uint64(profile.ContractInfo.ContractID)
	blocked = profile.ContractInfo.Blocked
	if blocked {
		return
	}

	client.SetWallet(fmt.Sprint(walletID))
	balanceResp, err := client.Balance.Current()
	if err != nil {
		return
	}

	if len(balanceResp.Accounts) == 0 {
		return
	}
	balance = balanceResp.Accounts[0].Balance.Amount
	return
}

func convertQiwiTxn(qiwiTxn qiwi.Txn) (txn models.Txn) {
	txn.ID = uint(qiwiTxn.TxnID)
	// TxnType    TxnType
	switch qiwiTxn.Type {
	case "IN":
		txn.TxnType = models.In
	case "OUT":
		txn.TxnType = models.Out
	case "QIWI_CARD":
		txn.TxnType = models.QiwiCard
	}
	// ProviderID uint // ?
	txn.ProviderID = uint(qiwiTxn.Provider.ID)
	// Amount     float64
	txn.Amount = qiwiTxn.Sum.Amount
	// CreatedAt  time.Time
	txn.CreatedAt = qiwiTxn.Date
	// Fee        float64
	txn.Fee = qiwiTxn.Commission.Amount
	// Status     Status
	switch qiwiTxn.Status {
	case "WAITING":
		txn.Status = models.Waiting
	case "SUCCESS":
		txn.Status = models.Success
	case "ERROR":
		txn.Status = models.Error
	}
	return
}

// GetLastTxns call qiwi api and returns last 50 txns
func GetLastTxns(token string, walletID uint64) (res []models.Txn, err error) {
	client := qiwi.New(token, qiwi.Debug, qiwi.Wallet(fmt.Sprint(walletID)))
	payments, err := client.History.Payments(50)
	if err != nil {
		return
	}
	for _, qiwiTxn := range payments.Data {
		res = append(res, convertQiwiTxn(qiwiTxn))
	}
	return
}
