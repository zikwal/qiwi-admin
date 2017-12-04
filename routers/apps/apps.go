// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package apps

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/notifier"
)

// Apps shows auth form
func Apps(ctx *context.Context) {
	apps, err := models.Apps.List(ctx.User.ID)
	if ctx.HasError(err) {
		return
	}

	ctx.Data["apps"] = apps
	ctx.Data["Title"] = "Приложения"
	ctx.HTML(200, "apps/apps")
}

func Create(ctx *context.Context) {
	name := ctx.Query("name")
	_, err := models.Apps.Create(ctx.User.ID, name)
	if ctx.HasError(err) {
		return
	}
	ctx.Redirect("/dashboard/apps")
}

func SaveWebHookURL(ctx *context.Context) {
	uri := ctx.Query("webhook_url")
	appID := ctx.ParamsInt(":appID")
	err := models.Apps.SetWebhook(uint(appID), uri)
	if ctx.HasError(err) {
		return
	}
	ctx.Redirect("/dashboard/apps")
}

// Test sends test notification
func Test(ctx *context.Context) {
	var (
		uri    = ctx.Query("webhook_url")
		coment = ctx.Query("comment")
	)

	txn := models.Txn{
		Amount:    100.0,
		Comment:   coment,
		Model:     gorm.Model{ID: uint(time.Now().UnixNano())},
		QiwiTxnID: uint(time.Now().UnixNano()),
	}

	cmd := notifier.NewCmd(uri, txn)

	notifier.Notify(cmd)

	ctx.Flash.Success("Запрос успешно отправлен")
	ctx.Redirect("/dashboard/apps")
}
