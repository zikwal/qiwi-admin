// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wallets

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Setting show and keep wallet setting
func Setting(ctx *context.Context) {

	var (
		walletID = uint(ctx.ParamsInt(":id"))
	)

	_, err := models.GetWallet(walletID)
	if err != nil {
		return
	}

	ctx.HTML(200, "wallets/setting")
}
