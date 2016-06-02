// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// AddAccount tbd
func (c *Client) AddAccount(req *AddAccountRequest) (accountID int64, err error) {
	var result AddAccountResult
	response, err := c.Request("AddAccount", req, NewReqID())
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Fatal(err)
		return 0, err
	}
	return result.Result.AccountID, nil
}

// GetAccountByName tbd
func (c *Client) GetAccountByName(req *GetAccountByNameRequest) (account Account, err error) {
	response, err := c.Request("GetAccountByName", req, NewReqID())
	if err != nil {
		return
	}

	var result GetAccountResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Fatal(err)
		return Account{}, err
	}
	log.Debugf("Returning account: %v", result.Result.Account)
	return result.Result.Account, err
}

// GetAccountByID tbd
func (c *Client) GetAccountByID(req *GetAccountByIDRequest) (account Account, err error) {
	var result GetAccountResult
	response, err := c.Request("GetAccountByID", req, NewReqID())
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Fatal(err)
		return account, err
	}
	return result.Result.Account, err
}
