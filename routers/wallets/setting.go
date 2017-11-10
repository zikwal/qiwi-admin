// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wallets

import (
	"fmt"
	"time"

	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Setting show and keep wallet setting
func Setting(ctx *context.Context) {

	var (
		walletID = uint(ctx.ParamsInt(":id"))
	)

	wallet, err := models.GetWallet(walletID, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	if ctx.Req.Method == "POST" {
		wallet.Token = ctx.Query("token")
		wallet.Name = ctx.Query("name")
		wallet.Limit = uint(ctx.QueryInt("limit"))
		wallet.TokenExpiry, _ = time.Parse("02.01.2006", ctx.Query("expiry"))
		err = wallet.Update(models.DB(),
			models.WalletDBSchema.Name,
			models.WalletDBSchema.Token,
			models.WalletDBSchema.TokenExpiry,
			models.WalletDBSchema.Limit)
		if ctx.HasError(err) {
			return
		}
		ctx.Redirect(fmt.Sprintf("/wallets/%d/setting", walletID))
		return
	}

	group, err := models.GetGroup(wallet.GroupID, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["group"] = group
	ctx.Data["wallet"] = wallet
	ctx.HTML(200, "wallets/setting")
}
