// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package groups

import (
	"fmt"

	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Setting show list of wallets
func Setting(ctx *context.Context) {

	if ctx.Req.Method == "POST" {
		saveSetting(ctx)
		return
	}

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

	groups, err := models.GetUserGroups(ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	autothistory, err := models.AutotransferLast(groupID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["auto_history"] = autothistory
	ctx.Data["group"] = group
	ctx.Data["groups"] = groups
	ctx.Data["wallets"] = wallets
	ctx.Data["Title"] = "Группа " + group.Name
	ctx.HTML(200, "groups/setting")
}

func saveSetting(ctx *context.Context) {
	var (
		groupID             = uint(ctx.ParamsInt64(":groupID"))
		autotransferGroupID = uint(ctx.QueryInt("autotransfer_group"))
	)

	group, err := models.GetGroup(groupID, ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	group.AutTransferObjectType = 0
	if autotransferGroupID != 0 {
		group.AutTransferObjectType = models.ObjectGroup
	}
	group.AutoTransferObjectID = autotransferGroupID

	err = group.Update(models.DB(), models.GroupDBSchema.AutTransferObjectType,
		models.GroupDBSchema.AutoTransferObjectID)
	if ctx.HasError(err) {
		return
	}

	ctx.Flash.Success("Настройки успешно сохранены")
	ctx.Redirect(fmt.Sprintf("/groups/%d/setting", groupID))
}
