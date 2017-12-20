// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"time"

	"github.com/zhuharev/qiwi-admin/pkg/localbitcoins"
)

type GroupTrade struct {
	Trades []int
}

func SaveGroupTrade(gr *GroupTrade) (err error) {
	return stormDB.Save(gr)
}

type Trade struct {
	ID int `storm:"id,increment"`

	WalletID       uint `storm:"index"`
	GroupID        uint `storm:"index"`
	QiwiWalletID   string
	InitialBalance float64
	OutgoingAmount uint

	State TradeState

	Events []TradeEvent

	HasError bool
	Error    string

	// localbitcoins ID
	ExternalTradeID  string
	OutgoingWalletID string

	// TODO:
	RemoteAds []localbitcoins.Ad

	RejectedAds []RejectedAd

	ChoosenAd int

	CreatedBy uint      `storm:"index"`
	CreatedAt time.Time `storm:"index"`
}

func InitializeTrade(t *Trade) (err error) {
	t.CreatedAt = time.Now()
	t.State = TradeInitialized
	err = stormDB.Save(t)
	return
}

type TradeEvent struct {
	Type      TradeEventType
	Msg       string
	CreatedAt time.Time
}

// NewInfoEvent returns TradeEvent
func NewInfoEvent(msg string, args ...interface{}) TradeEvent {
	message := fmt.Sprintf(msg, args)
	return TradeEvent{
		Type:      TradeEventInfo,
		Msg:       message,
		CreatedAt: time.Now(),
	}
}

// NewErrorEvent returns TradeEvent
func NewErrorEvent(msg string, args ...interface{}) TradeEvent {
	message := fmt.Sprintf(msg, args)
	return TradeEvent{
		Type:      TradeEventError,
		Msg:       message,
		CreatedAt: time.Now(),
	}
}

type TradeEventType int

const (
	TradeEventInfo TradeEventType = 1 + iota
	TradeEventError
)

type TradeState int

const (
	TradeInitialized TradeState = iota + 1
	TradeRemoteAdsLoad
	TradeRemoteCreating
	TradeOutgoingWalletIDWaiting
	TradeOutgoingWalletIDRecieved
	TradePaid
	TradeReleased
	TradeCanceled
)

// TradeAddEvent ads avent to trade and save them to DB
func TradeAddEvent(trade *Trade, event TradeEvent, newStates ...TradeState) error {
	trade.Events = append(trade.Events, event)
	if len(newStates) > 0 {
		trade.State = newStates[0]
	}
	return stormDB.Save(trade)
}

// TradeAddError ads avent to trade and save them to DB
// func TradeAddError(trade *Trade, event TradeEvent, newStates ...TradeState) error {
// 	trade.Events = append(trade.Events, event)
// 	if len(newStates) > 0 {
// 		trade.State = newStates[0]
// 	}
// 	return stormDB.Save(trade)
// }

// RejectedAd rejected ad
type RejectedAd struct {
	ID     int
	Reason string
}

// GetUserTrades returns user trades
func GetUserTrades(userID uint) (trades []Trade, err error) {
	// TODO: filter by userID
	var tmpTrades []Trade
	err = stormDB.All(&trades)
	for _, trade := range tmpTrades {
		if trade.CreatedBy == userID {
			trades = append(trades, trade)
		}
	}
	return
}
