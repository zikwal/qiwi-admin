// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package exchange

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/trader"
)

// Trade starts trade from group
func Trade(ctx *context.Context) {
	var (
		groupID = uint(ctx.QueryInt("group_id"))
	)

	wallets, err := models.GroupWallets(groupID)
	if ctx.HasError(err) {
		return
	}

	err = trader.Trade(ctx.User.LocalBitcoinsKey, ctx.User.LocalBitcoinsSecret, wallets)
	if ctx.HasError(err) {
		return
	}

	ctx.Flash.Success("Деньги успешно переведены")
	ctx.Redirect("/exchange")
}

// Wallet starts trade from group
func Wallet(ctx *context.Context) {
	var (
		walletID = uint(ctx.ParamsInt(":id"))
	)

	wallet, err := models.GetWallet(walletID, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	_, err = trader.TradeWallet(ctx.User.LocalBitcoinsKey, ctx.User.LocalBitcoinsSecret, *wallet)
	if ctx.HasError(err) {
		return
	}

	// TODO: sync wallet after trade

	ctx.Flash.Success("Деньги успешно переведены")
	ctx.Redirect("/exchange")
}

// Trades starts trade from group
func Trades(ctx *context.Context) {
	trades, err := models.GetUserTrades(ctx.User.ID)
	if ctx.HasError(err) {
		return
	}
	ctx.Data["trades"] = trades
	ctx.HTML(200, "exchange/trades")
}
