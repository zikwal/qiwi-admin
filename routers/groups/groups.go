// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package groups

import (
	"fmt"

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

// Get show list of wallets
func Get(ctx *context.Context) {
	var (
		groupID = uint(ctx.ParamsInt64(":groupID"))
	)
	wallets, err := models.GroupWallets(groupID)
	if ctx.HasError(err) {
		return
	}

	group, err := models.GetGroup(groupID, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["group"] = group
	ctx.Data["wallets"] = wallets
	ctx.Data["Title"] = "Мои кошельки"
	ctx.HTML(200, "groups/get")
}

// Create create an group
func Create(ctx *context.Context) {
	g, err := models.CreateGroup(ctx.Query("name"), ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Flash.Success("Группа создана")
	ctx.Redirect("/groups/" + fmt.Sprint(g.ID))
}

// Delete soft delete group and all group wallets
func Delete(ctx *context.Context) {
	var (
		groupID = uint(ctx.ParamsInt(":id"))
	)

	err := models.DeleteGroup(groupID)
	if ctx.HasError(err) {
		return
	}

	wallets, err := models.GroupWallets(groupID)
	if ctx.HasError(err) {
		return
	}

	for _, wallet := range wallets {
		err = (&wallet).Delete(models.DB())
		if ctx.HasError(err) {
			return
		}
	}

	ctx.Flash.Success("Группа успешно удалёна")
	ctx.Redirect("/dashboard")
}
