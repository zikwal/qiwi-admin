// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package localbitcoins

import (
	"fmt"
	"regexp"
	"time"
)

type messagesResponse struct {
	Data struct {
		MessageCount int `json:"message_count"`
		MessageList  []struct {
			Msg       string    `json:"msg"`
			CreatedAt time.Time `json:"created_at"`
			IsAdmin   bool      `json:"is_admin"`
			Sender    struct {
				Username      string    `json:"username"`
				FeedbackScore int       `json:"feedback_score"`
				TradeCount    string    `json:"trade_count"`
				LastOnline    time.Time `json:"last_online"`
				Name          string    `json:"name"`
			} `json:"sender"`
		} `json:"message_list"`
	} `json:"data"`
}

var (
	phoneRe = regexp.MustCompile(`((\+7|7|8)?9([0-9]){9})`)
)

func checkPaymentWalletID(key, secret, tradeID string) (string, error) {
	path := fmt.Sprintf("contact_messages/%s/", tradeID)
	var res messagesResponse
	err := SendAuthenticatedHTTPRequest(key, secret, "GET", path, nil, &res)
	if err != nil {
		return "", err
	}
	if (res.Data.MessageCount) == 0 {
		return "", fmt.Errorf("wallet not given")
	}
	arr := phoneRe.FindStringSubmatch(res.Data.MessageList[0].Msg)
	if len(arr) > 1 {
		return "+" + arr[1], nil
	}
	return "", fmt.Errorf("wallet not given")
}

func GetPaymentWalletID(key, secret, tradeID string) (walletID string, err error) {
	walletID, _ = checkPaymentWalletID(key, secret, tradeID)
	if walletID != "" {
		return walletID, nil
	}
	ticker := time.NewTicker(15 * time.Second)
	timeout := time.NewTimer(time.Hour)
	for {
		select {
		case <-ticker.C:
			walletID, err = checkPaymentWalletID(key, secret, tradeID)
			if err != nil {
				continue
			} else {
				return
			}
		case <-timeout.C:
			err = fmt.Errorf("Истекло время ожидания номера qiwi-кошелька")
			return
		}
	}
}
