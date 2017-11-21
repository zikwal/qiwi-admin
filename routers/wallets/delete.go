// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wallets

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Delete soft delete wallet from db
func Delete(ctx *context.Context) {
	var (
		id = uint(ctx.ParamsInt(":id"))
	)

	err := models.NewWalletQuerySet(models.DB()).IDEq(id).Delete()
	if ctx.HasError(err) {
		return
	}

	ctx.Flash.Success("Кошелёк успешно удалён")
	ctx.Redirect("/dashboard")
}
