// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package trader

import (
	"fmt"
	"strconv"

	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/localbitcoins"
	"github.com/zhuharev/qiwi-admin/pkg/qiwi"
)

func Trade(key, secret string, wallets []models.Wallet) (err error) {
	var (
		tradeIDs []int
	)

	for _, wallet := range wallets {
		var trade *models.Trade
		trade, err = TradeWallet(key, secret, wallet)
		if err != nil {
			return err
		}
		tradeIDs = append(tradeIDs, trade.ID)
	}

	return
}

func TradeWallet(key, secret string, wallet models.Wallet) (trade *models.Trade, err error) {

	amount := uint(wallet.Balance)

	trade = &models.Trade{
		WalletID:       wallet.ID,
		QiwiWalletID:   strconv.Itoa(int(wallet.WalletID)),
		InitialBalance: wallet.Balance,
		OutgoingAmount: amount,
		GroupID:        wallet.GroupID,
		CreatedBy:      wallet.OwnerID,
	}
	err = models.InitializeTrade(trade)
	if err != nil {
		return
	}

	err = models.TradeAddEvent(trade, models.NewInfoEvent("Запрос списка доступных трэйдеров"))
	if err != nil {
		return
	}
	ads, err := localbitcoins.GetAds()
	if err != nil {
		return
	}
	trade.RemoteAds = ads
	err = models.TradeAddEvent(trade, models.NewInfoEvent("Список трэйдеров загружен"))
	if err != nil {
		return
	}

	partner, ok, rejected := MustStartTrade(trade, key, secret, ads, amount)
	if !ok {
		err = fmt.Errorf("Не могу выбрать партнёра или создать трэйд")
		models.TradeAddEvent(trade, models.NewErrorEvent(err.Error()))
		return
	}

	trade.ChoosenAd = partner.Data.AdID
	trade.RejectedAds = rejected
	err = models.TradeAddEvent(trade, models.NewInfoEvent("Партнёр выбран, трэйд создан (%d). Ожидаю адрес qiwi-кошелька", trade.ChoosenAd), models.TradeOutgoingWalletIDWaiting)
	if err != nil {
		return
	}

	walletID, err := localbitcoins.GetPaymentWalletID(key, secret, trade.ExternalTradeID)
	if err != nil {
		return
	}

	trade.OutgoingWalletID = walletID
	err = models.TradeAddEvent(trade, models.NewInfoEvent("Получен адрес qiwi-кошелька: %s", walletID), models.TradeOutgoingWalletIDRecieved)
	if err != nil {
		return
	}

	_, err = qiwi.Transfer(wallet.Token, walletID, float64(amount))
	if err != nil {
		models.TradeAddEvent(trade, models.NewErrorEvent("Ошибка отправки денег: %s", err), models.TradeCanceled)
		return
	}
	err = models.TradeAddEvent(trade, models.NewInfoEvent("Деньги отправлены, завершаем сделку"), models.TradePaid)
	if err != nil {
		return
	}

	err = localbitcoins.MarkAsPaid(key, secret, trade.ExternalTradeID)
	if err != nil {
		models.TradeAddEvent(trade, models.NewErrorEvent("Ошибка пометки об оплате: %s", err), models.TradeCanceled)
		return
	}

	// TODO: release trade
	err = models.TradeAddEvent(trade, models.NewInfoEvent("Трэйд заверщён"), models.TradeReleased)
	if err != nil {
		return
	}

	return
}

func MustStartTrade(trade *models.Trade, key, secret string, ads []localbitcoins.Ad, amount uint) (choosen localbitcoins.Ad, found bool, rejAds []models.RejectedAd) {
	for _, ad := range ads {
		if ad.Data.Currency != "RUB" {
			rejAds = append(rejAds, models.RejectedAd{ID: ad.Data.AdID, Reason: "Валюты не рубли"})
			continue
		}
		if len(ad.Data.Profile.TradeCount) < 4 {
			rejAds = append(rejAds, models.RejectedAd{ID: ad.Data.AdID, Reason: "Мало трэйдов: " + ad.Data.Profile.TradeCount})
			continue
		}
		if minAmo, err := strconv.Atoi(ad.Data.MinAmount); err != nil {
			rejAds = append(rejAds, models.RejectedAd{ID: ad.Data.AdID, Reason: "Не могу перевести в цифры: " + ad.Data.MinAmount})
			continue
		} else if uint(minAmo) > amount {
			rejAds = append(rejAds, models.RejectedAd{ID: ad.Data.AdID, Reason: "У нас не хватает на минимальный порог: " + ad.Data.MinAmount})
			continue
		}
		tradeResp, err := localbitcoins.StartTrade(key, secret, ad, amount)
		if err != nil {
			rejAds = append(rejAds, models.RejectedAd{ID: ad.Data.AdID, Reason: "Ошибка старта трэйдинга: " + err.Error()})
			continue
		}
		choosen = ad
		found = true
		trade.ExternalTradeID = strconv.Itoa(tradeResp.Data.ContactID)
		return
	}
	return
}
