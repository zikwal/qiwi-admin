// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package localbitcoins

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-/gocryptotrader/common"
	"github.com/thrasher-/gocryptotrader/exchanges/nonce"
	dry "github.com/ungerik/go-dry"
)

const (
	localbitcoinsAPIURL = "https://localbitcoins.net"

	EndpointTradeCreate = "contact_create/"
)

var (
	non nonce.Nonce
)

func GetAds() ([]Ad, error) {
	return getAds()
}

func getAds() ([]Ad, error) {
	var res AdsResponse
	err := dry.FileUnmarshallJSON("https://localbitcoins.net/buy-bitcoins-online/qiwi/.json", &res)
	return res.Data.AdList, err
}

type AdsResponse struct {
	Data struct {
		AdList  []Ad `json:"ad_list"`
		AdCount int  `json:"ad_count"`
	} `json:"data"`
}

type Ad struct {
	Data struct {
		Profile struct {
			Username      string    `json:"username"`
			FeedbackScore int       `json:"feedback_score"`
			TradeCount    string    `json:"trade_count"`
			LastOnline    time.Time `json:"last_online"`
			Name          string    `json:"name"`
		} `json:"profile"`
		RequireFeedbackScore       int         `json:"require_feedback_score"`
		HiddenByOpeningHours       bool        `json:"hidden_by_opening_hours"`
		TradeType                  string      `json:"trade_type"`
		AdID                       int         `json:"ad_id"`
		TempPrice                  string      `json:"temp_price"`
		BankName                   string      `json:"bank_name"`
		PaymentWindowMinutes       int         `json:"payment_window_minutes"`
		TrustedRequired            bool        `json:"trusted_required"`
		MinAmount                  string      `json:"min_amount"`
		Visible                    bool        `json:"visible"`
		RequireTrustedByAdvertiser bool        `json:"require_trusted_by_advertiser"`
		TempPriceUsd               string      `json:"temp_price_usd"`
		Lat                        float64     `json:"lat"`
		AgeDaysCoefficientLimit    string      `json:"age_days_coefficient_limit"`
		IsLocalOffice              bool        `json:"is_local_office"`
		FirstTimeLimitBtc          interface{} `json:"first_time_limit_btc"`
		AtmModel                   interface{} `json:"atm_model"`
		City                       string      `json:"city"`
		LocationString             string      `json:"location_string"`
		Countrycode                string      `json:"countrycode"`
		Currency                   string      `json:"currency"`
		LimitToFiatAmounts         string      `json:"limit_to_fiat_amounts"`
		CreatedAt                  time.Time   `json:"created_at"`
		MaxAmount                  string      `json:"max_amount"`
		Lon                        float64     `json:"lon"`
		SmsVerificationRequired    bool        `json:"sms_verification_required"`
		RequireTradeVolume         float64     `json:"require_trade_volume"`
		OnlineProvider             string      `json:"online_provider"`
		MaxAmountAvailable         string      `json:"max_amount_available"`
		Msg                        string      `json:"msg"`
		RequireIdentification      bool        `json:"require_identification"`
		Email                      interface{} `json:"email"`
		VolumeCoefficientBtc       string      `json:"volume_coefficient_btc"`
	} `json:"data"`
	Actions struct {
		PublicView string `json:"public_view"`
	} `json:"actions"`
}

type CreateTradeResponse struct {
	Data struct {
		Message   string `json:"message"`
		Funded    bool   `json:"funded"`
		ContactID int    `json:"contact_id"`
	} `json:"data"`
	Actions struct {
		ContactURL string `json:"contact_url"`
	} `json:"actions"`
}

// GeneralError is an error capture type
type GeneralError struct {
	Error struct {
		Message   string `json:"message"`
		ErrorCode int    `json:"error_code"`
	} `json:"error"`
}

func StartTrade(key, secret string, ad Ad, amount uint) (res CreateTradeResponse, err error) {
	path := fmt.Sprintf("%s%d/", EndpointTradeCreate, ad.Data.AdID)
	err = SendAuthenticatedHTTPRequest(key, secret, "POST", path, url.Values{"amount": {strconv.Itoa(int(amount))}}, &res)
	if err != nil {
		return
	}
	// TODO:
	return
}

// SendAuthenticatedHTTPRequest sends an authenticated HTTP request to
// localbitcoins
func SendAuthenticatedHTTPRequest(key, secret, method, path string, values url.Values, result interface{}) (err error) {
	if non.Get() == 0 {
		non.Set(time.Now().UnixNano())
	} else {
		non.Inc()
	}

	payload := ""
	path = "/api/" + path

	if len(values) > 0 {
		payload = values.Encode()
	}

	message := non.String() + key + path + payload
	hmac := common.GetHMAC(common.HashSHA256, []byte(message), []byte(secret))
	headers := make(map[string]string)
	headers["Apiauth-Key"] = key
	headers["Apiauth-Nonce"] = non.String()
	headers["Apiauth-Signature"] = common.StringToUpper(common.HexEncodeToString(hmac))
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	log.Printf("Raw Path: \n%s\n", path)

	resp, err := common.SendHTTPRequest(method, localbitcoinsAPIURL+path, headers, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	log.Printf("Received raw: \n%s\n", resp)

	var errCapture GeneralError

	if err = common.JSONDecode([]byte(resp), &errCapture); err == nil {
		if len(errCapture.Error.Message) != 0 {
			return errors.New(errCapture.Error.Message)
		}
	}

	err = common.JSONDecode([]byte(resp), &result)
	if err != nil {
		return err
	}

	return nil
}

func MarkAsPaid(key, secret, tradeID string) (err error) {
	path := fmt.Sprintf("contact_mark_as_paid/%s/", tradeID)
	err = SendAuthenticatedHTTPRequest(key, secret, "POST", path, nil, nil)
	if err != nil {
		return
	}
	return nil
}
