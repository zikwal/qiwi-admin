// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package groups

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// List show list of wallets
func List(ctx *context.Context) {
	groups, err := models.GetUserGroups(ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["groups"] = groups
	ctx.Data["Title"] = "Мои группы"
	ctx.HTML(200, "groups/list")
}
