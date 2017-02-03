// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/units"
)

// Client is used to send API requests to a SolidFire system system
type Client struct {
	SVIP              string
	Endpoint          string
	DefaultAPIPort    int
	DefaultVolSize    int64 //bytes
	DefaultAccountID  int64
	DefaultTenantName string
	VolumeTypes       *[]VolType
	Config            *Config
}

// Config holds the configuration data for the Client to communicate with a SolidFire storage system
type Config struct {
	TenantName       string
	EndPoint         string
	DefaultVolSz     int64 //Default volume size in GiB
	MountPoint       string
	SVIP             string
	InitiatorIFace   string //iface to use of iSCSI initiator
	Types            *[]VolType
	LegacyNamePrefix string
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
	defaultSizeGiB    int64
	cfg               Config
)

// NewFromParameters is a factory method to create a new sfapi.Client object using the supplied parameters
func NewFromParameters(pendpoint string, pdefaultSizeGiB int64, psvip string, pcfg Config, pdefaultTenantName string) (c *Client, err error) {
	rand.Seed(time.Now().UTC().UnixNano())
	defSize := pdefaultSizeGiB * int64(units.GiB)
	SFClient := &Client{
		Endpoint:          pendpoint,
		DefaultVolSize:    defSize,
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
	log.Debug("sending request to SolidFire endpoint")
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	log.Debugf("sending request to: %+v", c.Endpoint)
	_ = json.Indent(&prettyJSON, data, "", "  ")
	log.Debugf("request json: %+v", string(prettyJSON.Bytes()))

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
	log.Debugf("response body: %+v", string(prettyJSON.Bytes()))
	errresp := APIError{}
	json.Unmarshal([]byte(body), &errresp)
	if errresp.Error.Code != 0 {
		log.Errorf("error detected in API response: %+v", errresp)
		return body, errors.New("device API error")
	}
	return body, nil
}

// NewReqID generates a random id for a request
func NewReqID() int {
	return rand.Intn(1000-1) + 1
}
