// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wallets

import (
	"fmt"

	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Show create an wallet
func Show(ctx *context.Context) {
	var (
		id = uint(ctx.ParamsInt64(":id"))
	)
	wallet, err := models.GetWallet(id, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	txns, err := models.GetWalletTxns(wallet.ID)
	if ctx.HasError(err) {
		return
	}

	group, err := models.GetGroup(wallet.GroupID, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["group"] = group
	ctx.Data["transactions"] = txns
	ctx.Data["wallet"] = wallet
	ctx.Data["Title"] = fmt.Sprintf("Кошелёк %s", wallet)
	ctx.HTML(200, "wallets/get")
}
