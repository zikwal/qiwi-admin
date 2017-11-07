// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wallets

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// List show list of wallets
func List(ctx *context.Context) {
	var (
		groupID = uint(ctx.ParamsInt64(":groupID"))
	)
	wallets, err := models.GroupWallets(groupID)
	if ctx.HasError(err) {
		return
	}

	group, err := models.GetGroup(groupID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["group"] = group
	ctx.Data["wallets"] = wallets
	ctx.Data["Title"] = "Мои кошельки"
	ctx.HTML(200, "wallets/list")
}
