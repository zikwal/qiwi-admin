// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

// ObjectType represend kind of object
type ObjectType = AppLevel

const (
  // ObjectWallet object is wallet
  ObjectWallet ObjectType = iota +1
  // ObjectGroup object is group
  ObjectGroup
  // ObjectAccount object is account
  ObjectAccount
)
