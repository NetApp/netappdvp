// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Client is used to send API requests to a SolidFire system system
type Client struct {
	SVIP              string
	Endpoint          string
	DefaultAPIPort    int
	DefaultAccountID  int64
	DefaultTenantName string
	VolumeTypes       *[]VolType
	Config            *Config
	AccessGroups      []int64
}

// Config holds the configuration data for the Client to communicate with a SolidFire storage system
type Config struct {
	TenantName       string
	EndPoint         string
	MountPoint       string
	SVIP             string
	InitiatorIFace   string //iface to use of iSCSI initiator
	Types            *[]VolType
	LegacyNamePrefix string
	AccessGroups     []int64
}

// VolType holds quality of service configuration data
type VolType struct {
	Type string
	QOS  QoS
}

var (
	endpoint          string
	svip              string
	configFile        string
	defaultTenantName string
	cfg               Config
)

// NewFromParameters is a factory method to create a new sfapi.Client object using the supplied parameters
func NewFromParameters(pendpoint string, psvip string, pcfg Config, pdefaultTenantName string) (c *Client, err error) {
	rand.Seed(time.Now().UTC().UnixNano())
	SFClient := &Client{
		Endpoint:          pendpoint,
		SVIP:              psvip,
		Config:            &pcfg,
		DefaultAPIPort:    443,
		VolumeTypes:       pcfg.Types,
		DefaultTenantName: pdefaultTenantName,
	}
	return SFClient, nil
}

// Request performs a json-rpc POST to the configured endpoint
func (c *Client) Request(method string, params interface{}, id int) (response []byte, err error) {
	var prettyJSON bytes.Buffer
	if c.Endpoint == "" {
		log.Error("endpoint is not set, unable to issue json-rpc requests")
		err = errors.New("no endpoint set")
		return nil, err
	}
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	log.Debugf("issuing request to SolidFire endpoint:  %+v", string(data))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	_ = json.Indent(&prettyJSON, data, "", "  ")
	http := &http.Client{Transport: tr}
	resp, err := http.Post(c.Endpoint,
		"json-rpc",
		strings.NewReader(string(data)))
	if err != nil {
		log.Errorf("error response from SolidFire API request: %v", err)
		return nil, errors.New("device API error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	_ = json.Indent(&prettyJSON, body, "", "  ")
	errresp := APIError{}
	json.Unmarshal([]byte(body), &errresp)

	// NOTE(jdg): We removed the raw dump of the json response, was mostly just
	// more noise than anything useful and clogged up the logfile.  Might be
	// cool to add a config option that lets one turn this on/off for special
	// debug purposes, but in general I think it's just too much noise,
	// especially on things like a List
	if errresp.Error.Code != 0 {
		log.Warningf("error detected in API response: %+v", errresp)
		return body, fmt.Errorf("device API error: %+v", errresp.Error.Name)
	}
	return body, nil
}

// NewReqID generates a random id for a request
func NewReqID() int {
	return rand.Intn(1000-1) + 1
}
