// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncronizer

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/notifier"
	"github.com/zhuharev/qiwi-admin/pkg/qiwi"
)

var (
	checkInterval = 2 * time.Minute
)

// NewContext init synchronizer adn run it in background
func NewContext() (err error) {
	go func() {
		for {
			wallets, _ := models.GetAllWallets()
			for _, wallet := range wallets {
				err = Sync(wallet.ID)
				if err != nil {
					color.Red("%s", err)
				}
			}
			time.Sleep(checkInterval)
		}
	}()
	return
}

// Sync specified wallet
func Sync(walletID uint) (err error) {
	wallet, err := models.GetWallet(walletID)
	if err != nil {
		return
	}

	_, blocked, balance, err := qiwi.CheckToken(wallet.Token)
	if err != nil {
		return
	}

	if wallet.Balance == balance && time.Since(wallet.UpdatedAt) < time.Hour {
		return
	}

	// todo notify change
	wallet.Blocked = blocked
	wallet.Balance = balance
	err = wallet.Update(models.DB(),
		models.WalletDBSchema.Blocked,
		models.WalletDBSchema.Balance)
	if err != nil {
		return
	}

	lastTxnID, _ := models.GetLastQiwiTxn(wallet.ID)

	txns, err := qiwi.GetLastTxns(wallet.Token, wallet.WalletID)
	if err != nil {
		return err
	}

	group, err := models.GetGroup(wallet.GroupID)
	if err != nil {
		return err
	}

	var (
		insertTxns []models.Txn
	)

	for _, txn := range txns {
		if txn.QiwiTxnID > lastTxnID {
			insertTxns = append(insertTxns, txn)

			// make webhook
			err = notifier.NotifyTxn(wallet, txn)
			if err != nil {
				color.Red("Error when making webhook: %s", err)
			}

			if txn.TxnType == models.In && group.AutTransferObjectType == models.ObjectGroup && group.AutoTransferObjectID != 0 {
				targetWallet, err := models.GetGroupFreeWallet(group.AutoTransferObjectID, uint(txn.Amount))
				if err != nil {
					color.Red("Error when getting free master-group wallet: %s", err)
					continue
				}
				color.Green("[autotransfer] from wallet %d to %d amount: %f", wallet.ID, targetWallet.ID, txn.Amount)
				_, err = qiwi.Transfer(wallet.Token, fmt.Sprintf("+%d", targetWallet.WalletID), txn.Amount)
				if err != nil {
					log.Printf("[autotransfer] error transfer from group: %d", err)
					continue
				}

				autotransferLogEntry := models.Autotransfer{
					SourceID:   wallet.ID,
					SourceType: models.ObjectGroup,
					TargetID:   fmt.Sprint(targetWallet.WalletID),
					TargetType: models.ObjectWallet,
					Amount:     uint(txn.Amount),
				}

				err = models.AutotransferSave(&autotransferLogEntry)
				if err != nil {
					log.Printf("[autotransfer] error save autotransfer from group: %d", err)
				}

			}

		}
	}

	err = models.CreateMultipleTxns(walletID, insertTxns)
	if err != nil {
		return
	}

	if wallet.TotalSynced.IsZero() || time.Since(wallet.TotalSynced) > time.Minute*10 {
		inc, out, err := qiwi.GetStat(wallet.Token, wallet.WalletID)
		if err != nil {
			return err
		}
		wallet.TotalMonthIncoming = inc
		wallet.TotalMonthOutgoing = out
		wallet.TotalSynced = time.Now()
		err = wallet.Update(models.DB(),
			models.WalletDBSchema.TotalMonthOutgoing,
			models.WalletDBSchema.TotalMonthIncoming,
			models.WalletDBSchema.TotalSynced)
		if err != nil {
			return err
		}
	}

	return
}
