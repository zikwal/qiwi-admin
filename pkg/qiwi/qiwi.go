// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package qiwi

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jinzhu/now"
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
	payments, err := client.Payments.History(50)
	if err != nil {
		return
	}
	for _, qiwiTxn := range payments.Data {
		res = append(res, convertQiwiTxn(qiwiTxn))
	}
	return
}

// GetStat return stat of current month
func GetStat(token string, walletID uint64) (incoming, outgoing float64, err error) {
	client := qiwi.New(token, qiwi.Debug, qiwi.Wallet(fmt.Sprint(walletID)))
	stat, err := client.Payments.Stat(now.BeginningOfMonth(), now.EndOfMonth())
	if err != nil {
		return
	}
	for _, a := range stat.IncomingTotal {
		incoming += a.Amount
	}
	for _, a := range stat.OutgoingTotal {
		outgoing += a.Amount
	}
	return
}

// DetectProvider detect provider
func DetectProvider(token string, to string) (id int, err error) {
	client := qiwi.New(token, qiwi.Debug)
	id, err = client.Cards.Detect(to)
	if err != nil {
		return
	}
	return
}

// Transfer transfer money
func Transfer(token, to string, amount float64, comments ...string) (transactionID uint, err error) {
	client := qiwi.New(token, qiwi.Debug)

	var (
		// qiwi to qiwi
		providerID = 99
	)

	if !strings.HasPrefix(to, "+") {
		providerID, err = client.Cards.Detect(to)
		if err != nil {
			return
		}
	}

	return TransferWithProvider(providerID, token, to, amount, comments...)
}

func TransferWithProvider(providerID int, token, to string, amount float64, comments ...string) (transactionID uint, err error) {
	client := qiwi.New(token, qiwi.Debug)
	_, err = client.Cards.Payment(providerID, amount, to, comments...)
	if err != nil {
		return
	}
	return
}

func calculateTransferAmount(balance float64, restAmount float64, comission qiwi.ComissionResponse) (amount float64) {
	for _, com := range comission.Content.Terms.Commission.Ranges {
		if com.Fixed != 0 {
			amount -= com.Fixed
		}
		if com.Rate != 0 {
			amount /= 1.0 + com.Rate
		}
	}
	return
}

// TransferFromGroup transfer from group wallets to target
func TransferFromGroup(groupID, userID uint, to string, restAmount float64) (errs []error) {
	wallets, err := models.GroupWallets(groupID)
	if err != nil {
		errs = []error{err}
		return
	}

	var (
		providerID        = 0
		comissionResponse qiwi.ComissionResponse
	)

	for _, wallet := range wallets {
		walletID, blocked, balance, err := CheckToken(wallet.Token)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if blocked {
			errs = append(errs, fmt.Errorf("Кошелёк %d заблокирован", walletID))
			continue
		}
		if balance < restAmount {
			errs = append(errs, fmt.Errorf("На кошельке %d не хватает средств для вывода", walletID))
			continue
		}
		client := qiwi.New(wallet.Token, qiwi.Debug)
		if providerID == 0 {
			providerID, err = client.Cards.Detect(to)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
		if len(comissionResponse.Content.Terms.Commission.Ranges) == 0 {
			comissionResponse, err = client.Payments.Comission(providerID)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		amount := calculateTransferAmount(balance, restAmount, comissionResponse)
		_, err = TransferWithProvider(providerID, wallet.Token, to, amount)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return
}

// Fee returns fee of payment
func Fee(token string, providerID int, to string, amount float64) (fee float64, err error) {
	client := qiwi.New(token, qiwi.Debug)
	feeResp, err := client.Payments.SpecialComission(providerID, to, amount)
	return feeResp.QwCommission.Amount, err
}

// DetectFee detet provider and after fee
func DetectFee(token string, to string, amount float64) (fee float64, err error) {
	providerID, err := DetectProvider(token, to)
	if err != nil {
		return
	}
	return Fee(token, providerID, to, amount)
}
