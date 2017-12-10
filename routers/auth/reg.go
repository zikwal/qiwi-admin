// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/setting"
)

// Reg shows registration form
func Reg(ctx *context.Context, af models.AuthForm) {
	if setting.App.Reg.Disabled {
		ctx.Error(404, "Страница не найдена")
		return
	}
	if ctx.Req.Method == "POST" {
		user, err := models.CreateUser(af)
		if ctx.HasError(err) {
			return
		}

		err = ctx.Session.Set("user_id", user.ID)
		if ctx.HasError(err) {
			return
		}
		ctx.Redirect("/dashboard")
		return
	}

	ctx.Data["Title"] = "Регистрация"
	ctx.HTML(200, "auth/reg")
}
