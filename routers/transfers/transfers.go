// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package transfers

import (
	"github.com/zhuharev/qiwi-admin/models"

	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/qiwi"
)

// Transfer transfer money
func Transfer(ctx *context.Context) {

	var (
		//from = ctx.Query("from")
		to       = ctx.Query("to")
		amount   = ctx.QueryFloat64("amount")
		walletID = uint(ctx.QueryInt("wallet_id"))
	)

	if ctx.Req.Method == "POST" {
		wallet, err := models.GetWallet(walletID)
		if ctx.HasError(err) {
			return
		}

		_, err = qiwi.Transfer(wallet.Token, to, amount)
		if ctx.HasError(err) {
			return
		}
	}

	ctx.Data["walletID"] = walletID
	ctx.HTML(200, "transfers/transfer")
}
