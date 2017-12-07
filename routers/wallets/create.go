// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wallets

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/qiwi"
	"github.com/zhuharev/qiwi-admin/pkg/syncronizer"
)

// Create create an wallet
func Create(ctx *context.Context) {
	var (
		name    = ctx.Query("name")
		token   = ctx.Query("token")
		ownerID = ctx.User.ID
	)

	walletID, blocked, balance, err := qiwi.CheckToken(token)
	if ctx.HasError(err) {
		color.Red("Ошибка при проверке кошелька: %s", err)
		return
	}
	wallet := new(models.Wallet)
	wallet.Name = name
	wallet.Balance = balance
	wallet.Blocked = blocked
	wallet.WalletID = walletID
	wallet.Token = token
	wallet.OwnerID = ownerID
	wallet.GroupID = uint(ctx.QueryInt64("group"))
	wallet.Limit = 15000

	err = models.CreateWallet(wallet)
	if ctx.HasError(err) {
		return
	}

	err = syncronizer.Sync(wallet.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["Title"] = "Мои кошельки"
	ctx.Redirect(fmt.Sprintf("/wallets/%d", wallet.ID))
}
