// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package apps

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Apps shows auth form
func Apps(ctx *context.Context) {
	apps,err:=models.Apps.List(ctx.User.ID)
  if  ctx.HasError(err) {
    return
  }

ctx.Data["apps"] = apps
	ctx.Data["Title"] = "Приложения"
	ctx.HTML(200, "apps/apps")
}

func Create(ctx *context.Context) {
  name:= ctx.Query("name")
  _,err := models.Apps.Create(ctx.User.ID, name)
  if ctx.HasError(err) {
    return
  }
  ctx.Redirect("/dashboard/apps")
}

func SaveWebHookURL(ctx *context.Context) {
  uri := ctx.Query("webhook_url")
  appID := ctx.ParamsInt(":appID")
  err := models.Apps.SetWebhook(uint(appID), uri)
  if  ctx.HasError(err) {
    return
  }
  ctx.Redirect("/dashboard/apps")
}
