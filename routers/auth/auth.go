// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Auth shows auth form
func Auth(ctx *context.Context, af models.AuthForm) {
	if ctx.Req.Method == "POST" {
		user, err := models.GetUserByAuthForm(af)
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

	ctx.Data["Title"] = "Авторизация"
	ctx.HTML(200, "auth/auth")
}

// RedirectAutorized if user authorized, and try open index page
// he will be redirected to dashboard
func RedirectAutorized(ctx *context.Context) {
	if ctx.Autorized() {
		ctx.Redirect("/dashboard")
		return
	}
}

// MustAuthorized if user authorized, and try open index page
// he will be redirected to dashboard
func MustAuthorized(ctx *context.Context) {
	if !ctx.Autorized() {
		ctx.Redirect("/auth")
		return
	}
}

// Logout delete req cookies and session user_id
func Logout(ctx *context.Context) {
	ctx.Session.Delete("user_id")
	ctx.SetCookie("s", "")
	ctx.Redirect("/auth")
}
