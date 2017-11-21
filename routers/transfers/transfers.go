// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package transfers

import (
	"strings"

	"github.com/zhuharev/qiwi-admin/models"

	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/qiwi"
	"github.com/zhuharev/qiwi-admin/pkg/syncronizer"
)

// Transfer transfer money
func Transfer(ctx *context.Context) {

	var (
		//from = ctx.Query("from")
		to          = ctx.Query("to")
		amount      = ctx.QueryFloat64("amount")
		walletID    = uint(ctx.QueryInt("wallet_id"))
		needShowFee = ctx.QueryBool("show_fee")
		comment     = ctx.Query("comment")
	)

	ctx.Data["to"] = to
	ctx.Data["amount"] = amount
	ctx.Data["comment"] = comment

	if ctx.Req.Method == "POST" {
		wallet, err := models.GetWallet(walletID, ctx.User.ID)
		if ctx.HasError(err) {
			return
		}

		if needShowFee {
			ctx.Data["wallet"] = wallet

			if strings.HasPrefix(to, "+") {
				ctx.Data["fee"] = 0
			} else {
				ctx.Data["fee"], err = qiwi.DetectFee(wallet.Token, to, amount)
			}
			if ctx.HasError(err) {
				return
			}
		} else {
			_, err = qiwi.Transfer(wallet.Token, to, amount, comment)
			if ctx.HasError(err) {
				return
			}
			err = syncronizer.Sync(wallet.ID)
			if ctx.HasError(err) {
				return
			}
			ctx.Redirect("/dashboard")
			return
		}
	}

	ctx.Data["walletID"] = walletID
	ctx.HTML(200, "transfers/transfer")
}
