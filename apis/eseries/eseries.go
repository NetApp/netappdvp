// Copyright 2016 NetApp, Inc. All Rights Reserved.

// This package provides a high-level interface to the E-series Web Services Proxy REST API.
package eseries

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/utils"
)

const MAX_NAME_LENGTH int = 30
const HTTP_TIMEOUT_SECONDS = 10
const NULL_REF string = "0000000000000000000000000000000000000000"
const HOST_MAPPING_TYPE = "host"
const HOST_GROUP_MAPPING_TYPE = "cluster"
const DEFAULT_POOL_SEARCH_PATTERN = ".+"

// DriverConfig holds configuration data for the API driver object.
type DriverConfig struct {
	// Web Proxy Services Info
	WebProxyHostname  string
	WebProxyPort      string
	WebProxyUseHTTP   bool
	WebProxyVerifyTLS bool
	Username          string
	Password          string

	// Array Info
	ControllerA   string
	ControllerB   string
	PasswordArray string

	// Options
	PoolNameSearchPattern string
	DebugTraceFlags       map[string]bool

	// Host Connectivity
	HostDataIP string //for iSCSI with multipathing this can be either IP or host

	// Internal Config Variables
	ArrayID                       string // Unique ID for array once added to web proxy services
	CompiledPoolNameSearchPattern *regexp.Regexp

	// Storage protocol of the driver (iSCSI, FC, etc)
	Protocol    string
	AccessGroup string
	HostType    string

	DriverName    string
	DriverVersion string
	ConfigVersion int
}

// Driver is the object to use for interacting with the E-series API.
type ESeriesAPIDriver struct {
	config *DriverConfig
	m      *sync.Mutex
}

// NewDriver is a factory method for creating a new instance.
func NewDriver(config DriverConfig) *ESeriesAPIDriver {
	d := &ESeriesAPIDriver{
		config: &config,
		m:      &sync.Mutex{},
	}

	// Initialize internal config variables
	d.config.ArrayID = ""

	compiledRegex, err := regexp.Compile(d.config.PoolNameSearchPattern)
	if err != nil {
		log.WithFields(log.Fields{
			"PoolNameSearchPattern": d.config.PoolNameSearchPattern,
			"DefaultSearchPattern":  DEFAULT_POOL_SEARCH_PATTERN,
			"Error":                 err,
		}).Warn("Could not compile PoolNameSearchPattern regular expression, using default pattern.")
		compiledRegex, _ = regexp.Compile(".+")
	}
	d.config.CompiledPoolNameSearchPattern = compiledRegex

	volumeTags = []VolumeTag{
		{Key: "API", Value: "netappdvp-" + config.DriverVersion},
		{Key: "eBI", Value: "Containers-Docker"},
		{Key: "IF", Value: d.config.Protocol},
		{Key: "netappdvp", Value: config.DriverName},
	}

	return d
}

var volumeTags []VolumeTag

// InvokeAPI makes a REST call to the Web Services Proxy. The body must be a marshaled JSON byte array (or nil).
// The method is the HTTP verb (i.e. GET, POST, ...).  The resource path is appended to the base URL to identify
// the desired server resource; it should start with '/'.
func (d ESeriesAPIDriver) InvokeAPI(requestBody []byte, method string, resourcePath string) (*http.Response, []byte, error) {

	// Default to secure connection
	scheme, port := "https", "8443"

	// Allow insecure override
	if d.config.WebProxyUseHTTP {
		scheme, port = "http", "8080"
	}

	// Allow port override
	if d.config.WebProxyPort != "" {
		port = d.config.WebProxyPort
	}

	// Build URL
	url := scheme + "://" + d.config.WebProxyHostname + ":" + port + "/devmgr/v2/storage-systems/" + d.config.ArrayID + resourcePath

	var request *http.Request
	var err error

	if requestBody == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	}
	if err != nil {
		return nil, nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(d.config.Username, d.config.Password)

	// Allow certificate validation override
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !d.config.WebProxyVerifyTLS,
		},
	}

	if d.config.DebugTraceFlags["api"] {
		if method == "POST" && resourcePath == "" {
			// Suppress the empty POST body since it contains the array password
			d.logHttpRequest(request, []byte("<suppressed>"))
		} else {
			d.logHttpRequest(request, requestBody)
		}
	}

	// Send the request
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT_SECONDS * time.Second),
	}
	response, err := client.Do(request)
	if err != nil {
		log.Warnf("Error communicating with Web Services Proxy. %v", err)
		return nil, nil, err
	}
	defer response.Body.Close()

	responseBody := []byte{}
	if err == nil {

		responseBody, err = ioutil.ReadAll(response.Body)

		if method == "GET" && resourcePath == "/volumes" {
			// Suppress the potentially huge GET /volumes body unless asked for explicitly
			if d.config.DebugTraceFlags["api_get_volumes"] {
				d.logHttpResponse(response, responseBody)
			} else if d.config.DebugTraceFlags["api"] {
				d.logHttpResponse(response, []byte("<suppressed>"))
			}
		} else {
			if d.config.DebugTraceFlags["api"] {
				d.logHttpResponse(response, responseBody)
			}
		}
	}

	return response, responseBody, err
}

func (d ESeriesAPIDriver) logHttpRequest(request *http.Request, requestBody []byte) {
	header := ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
	footer := "--------------------------------------------------------------------------------"
	var body string
	if requestBody == nil {
		body = "<nil>"
	} else {
		body = string(requestBody)
	}
	log.Debugf("\n%s\n%s %s\nHeaders: %v\nBody: %s\n%s", header, request.Method, request.URL, request.Header, body, footer)
}

func (d ESeriesAPIDriver) logHttpResponse(response *http.Response, responseBody []byte) {
	header := "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	footer := "================================================================================"
	var body string
	if responseBody == nil {
		body = "<nil>"
	} else {
		body = string(responseBody)
	}
	log.Debugf("\n%s\nStatus: %s\nHeaders: %v\nBody: %s\n%s", header, response.Status, response.Header, body, footer)
}

// Connect connects to the Web Services Proxy and registers the array with it.
func (d ESeriesAPIDriver) Connect() (string, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "Connect",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> Connect")
		defer log.WithFields(fields).Debug("<<<< Connect")
	}

	// Send a login/connect request for array to web services proxy
	request := MsgConnect{[]string{d.config.ControllerA, d.config.ControllerB}, d.config.PasswordArray}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("Could not marshal JSON request: %v. %v", request, err)
	}

	// Send the message
	response, responseBody, err := d.InvokeAPI(jsonRequest, "POST", "")
	if err != nil {
		return "", fmt.Errorf("Could not log into the Web Services Proxy. %v", err)
	}

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Could not add storage array to Web Services Proxy. Status code: %d", response.StatusCode)
	}

	// Parse JSON data
	responseData := MsgConnectResponse{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return "", fmt.Errorf("Could not parse connect response: %s. %v", string(responseBody), err)
	}

	if responseData.ArrayID == "" {
		return "", errors.New("Invalid ArrayID received from Web Services Proxy.")
	}

	d.config.ArrayID = responseData.ArrayID
	AlreadyRegistered := responseData.AlreadyExists

	log.WithFields(log.Fields{
		"ArrayID":           d.config.ArrayID,
		"AlreadyRegistered": AlreadyRegistered,
	}).Debug("Connected to Web Services Proxy.")

	return d.config.ArrayID, nil
}

// GetControllers returns an array containing all the controllers in the storage system.
func (d ESeriesAPIDriver) GetControllers() ([]Controller, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetControllers",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> GetControllers")
		defer log.WithFields(fields).Debug("<<<< GetControllers")
	}

	// Query volumes on array
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/controllers")
	if err != nil {
		return nil, errors.New("Failed to read controllers.")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to read controllers. Status code: %d", response.StatusCode)
	}

	controllers := make([]Controller, 0)
	err = json.Unmarshal(responseBody, &controllers)
	if err != nil {
		return nil, fmt.Errorf("Could not parse controller data: %s. %v", string(responseBody), err)
	}

	log.WithField("Count", len(controllers)).Debug("Read controllers.")

	return controllers, nil
}

// ListNodeSerialNumbers returns an array containing the controller serial numbers for this storage system.
func (d ESeriesAPIDriver) ListNodeSerialNumbers() ([]string, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "ListNodeSerialNumbers",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> ListNodeSerialNumbers")
		defer log.WithFields(fields).Debug("<<<< ListNodeSerialNumbers")
	}

	serialNumbers := make([]string, 0, 0)

	controllers, err := d.GetControllers()
	if err != nil {
		return serialNumbers, err
	}

	// Get the serial numbers
	for _, controller := range controllers {
		serialNumber := strings.TrimSpace(controller.SerialNumber)
		if serialNumber != "" {
			serialNumbers = append(serialNumbers, serialNumber)
		}
	}

	log.WithFields(log.Fields{
		"Count":         len(serialNumbers),
		"SerialNumbers": serialNumbers,
	}).Debug("Read serial numbers.")

	return serialNumbers, nil
}

// GetVolumePools reads all pools on the array, including volume groups and dynamic disk pools. It then
// filters them based on several selection parameters and returns the ones that match.
func (d ESeriesAPIDriver) GetVolumePools(mediaType string, minFreeSpaceBytes uint64, poolName string) ([]VolumeGroupEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":            "GetVolumePools",
			"Type":              "ESeriesAPIDriver",
			"mediaType":         mediaType,
			"minFreeSpaceBytes": minFreeSpaceBytes,
			"poolName":          poolName,
		}
		log.WithFields(fields).Debug(">>>> GetVolumePools")
		defer log.WithFields(fields).Debug("<<<< GetVolumePools")
	}

	// Get the storage pools (includes volume RAID groups and dynamic disk pools)
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/storage-pools")
	if err != nil {
		return nil, fmt.Errorf("Could not get storage pools. %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not get storage pools. Status code: %d", response.StatusCode)
	}

	// Parse JSON data
	allPools := make([]VolumeGroupEx, 0)
	if err := json.Unmarshal(responseBody, &allPools); err != nil {
		return nil, fmt.Errorf("Could not parse storage pool data: %s. %v", string(responseBody), err)
	}

	// Return only pools that match the requested criteria
	matchingPools := make([]VolumeGroupEx, 0)
	for _, pool := range allPools {

		log.WithFields(log.Fields{
			"Name": pool.Label,
			"Pool": pool,
		}).Debug("Considering pool.")

		// Pool must match regex from config
		if !d.config.CompiledPoolNameSearchPattern.MatchString(pool.Label) {
			log.WithFields(log.Fields{"Name": pool.Label}).Debug("Pool does not match search pattern.")
			continue
		}

		// Pool must be online
		if pool.IsOffline {
			log.WithFields(log.Fields{"Name": pool.Label}).Debug("Pool is offline.")
			continue
		}

		// Pool name
		if poolName != "" {
			if poolName != pool.Label {
				log.WithFields(log.Fields{
					"Name":          pool.Label,
					"RequestedName": poolName,
				}).Debug("Pool does not match requested pool name.")
				continue
			}
		}

		// Drive media type
		if mediaType != "" {
			if mediaType != pool.DriveMediaType {
				log.WithFields(log.Fields{
					"Name":               pool.Label,
					"MediaType":          pool.DriveMediaType,
					"RequestedMediaType": mediaType,
				}).Debug("Pool does not match requested media type.")
				continue
			}
		}

		// Free space
		if minFreeSpaceBytes > 0 {
			poolFreeSpace, err := strconv.ParseUint(pool.FreeSpace, 10, 64)
			if err != nil {
				log.WithFields(log.Fields{
					"Name":  pool.Label,
					"Error": err,
				}).Warn("Could not parse free space for pool.")
				continue
			}
			if poolFreeSpace < minFreeSpaceBytes {
				log.WithFields(log.Fields{
					"Name":           pool.Label,
					"FreeSpace":      poolFreeSpace,
					"RequestedSpace": minFreeSpaceBytes,
				}).Debug("Pool does not have sufficient free space.")
				continue
			}
		}

		// Everything matched
		matchingPools = append(matchingPools, pool)
	}

	return matchingPools, nil
}

// GetVolumes returns an array containing all the volumes on the array.
func (d ESeriesAPIDriver) GetVolumes() ([]VolumeEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetVolumes",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> GetVolumes")
		defer log.WithFields(fields).Debug("<<<< GetVolumes")
	}

	// Query volumes on array
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/volumes")
	if err != nil {
		return nil, errors.New("Failed to read volumes.")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to read volumes. Status code: %d", response.StatusCode)
	}

	volumes := make([]VolumeEx, 0)
	err = json.Unmarshal(responseBody, &volumes)
	if err != nil {
		return nil, fmt.Errorf("Could not parse volume data: %s. %v", string(responseBody), err)
	}

	log.WithField("Count", len(volumes)).Debug("Read volumes.")

	return volumes, nil
}

// ListVolumes returns an array containing all the volume names on the array.
func (d ESeriesAPIDriver) ListVolumes() ([]string, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "ListVolumes",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> ListVolumes")
		defer log.WithFields(fields).Debug("<<<< ListVolumes")
	}

	volumes, err := d.GetVolumes()
	if err != nil {
		return nil, err
	}

	volumeNames := make([]string, 0, len(volumes))

	for _, vol := range volumes {
		volumeNames = append(volumeNames, vol.Label)
	}

	log.WithField("Count", len(volumeNames)).Debug("Read volume names.")

	return volumeNames, nil
}

// GetVolume returns a volume structure from the array whose label matches the specified name. Use this method sparingly,
// at most once per workflow, because the Web Services Proxy does not support server-side filtering so the only choice is to
// read all volumes to find the one of interest. Most methods in this module operate on the returned VolumeEx structure, not
// the volume name, to minimize the need for calling this method.
func (d ESeriesAPIDriver) GetVolume(name string) (VolumeEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetVolume",
			"Type":   "ESeriesAPIDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> GetVolume")
		defer log.WithFields(fields).Debug("<<<< GetVolume")
	}

	volumes, err := d.GetVolumes()
	if err != nil {
		return VolumeEx{}, err
	}

	for _, vol := range volumes {
		if vol.Label == name {
			return vol, nil
		}
	}

	return VolumeEx{}, nil
}

// CreateVolume creates a volume (i.e. a LUN) on the array, and it returns the resulting VolumeEx structure.
func (d ESeriesAPIDriver) CreateVolume(
	name string, volumeGroupRef string, size uint64, mediaType, fstype string,
) (VolumeEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":         "CreateVolume",
			"Type":           "ESeriesAPIDriver",
			"name":           name,
			"volumeGroupRef": volumeGroupRef,
			"size":           size,
			"mediaType":      mediaType,
		}
		log.WithFields(fields).Debug(">>>> CreateVolume")
		defer log.WithFields(fields).Debug("<<<< CreateVolume")
	}

	// Ensure that we do not exceed the maximum allowed volume length
	if len(name) > MAX_NAME_LENGTH {
		return VolumeEx{}, fmt.Errorf("The volume name %v exceeds the maximum length of %d characters", name, MAX_NAME_LENGTH)
	}

	// Copy static volume metadata and add fstype
	tags := append([]VolumeTag(nil), volumeTags...)
	tags = append(tags, VolumeTag{"fstype", fstype})

	// Set up the volume create request
	request := VolumeCreateRequest{
		VolumeGroupRef: volumeGroupRef,
		Name:           name,
		SizeUnit:       "kb",
		Size:           int(size / 1024), // The API requires Size to be an int (not int64) so pass as an int but in KB.
		SegmentSize:    128,
		VolumeTags:     tags,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return VolumeEx{}, fmt.Errorf("Could not marshal JSON request: %v. %v", request, err)
	}

	// Create the volume
	response, responseBody, err := d.InvokeAPI(jsonRequest, "POST", "/volumes")
	if err != nil {
		return VolumeEx{}, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK {

		err = d.getErrorFromHTTPResponse(response, responseBody)
		return VolumeEx{}, fmt.Errorf("Could not create volume %s. %v", name, err)

	} else {

		// Parse JSON volume data
		vol := VolumeEx{}
		if err := json.Unmarshal(responseBody, &vol); err != nil {
			return VolumeEx{}, fmt.Errorf("Could not parse API response: %s. %v", string(responseBody), err)
		}

		log.WithFields(log.Fields{
			"Name":           vol.Label,
			"VolumeRef":      vol.VolumeRef,
			"VolumeGroupRef": vol.VolumeGroupRef,
		}).Debug("Created volume.")

		return vol, nil
	}
}

// DeleteVolume deletes a volume from the array.
func (d ESeriesAPIDriver) DeleteVolume(volume VolumeEx) error {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "DeleteVolume",
			"Type":   "ESeriesAPIDriver",
			"name":   volume.Label,
		}
		log.WithFields(fields).Debug(">>>> DeleteVolume")
		defer log.WithFields(fields).Debug("<<<< DeleteVolume")
	}

	// Remove this volume from storage array
	response, responseBody, err := d.InvokeAPI(nil, "DELETE", "/volumes/"+volume.VolumeRef)
	if err != nil {
		return fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {

		err = d.getErrorFromHTTPResponse(response, responseBody)
		return fmt.Errorf("Could not destroy volume %s. %v", volume.Label, err)
	}

	log.WithFields(log.Fields{
		"Name":           volume.Label,
		"VolumeRef":      volume.VolumeRef,
		"VolumeGroupRef": volume.VolumeGroupRef,
	}).Debug("Deleted volume.")

	return nil
}

// EnsureHostForIQN handles automatic E-series Host and Host Group creation. Given the IQN of a host, this method
// verifies whether a Host is already configured on the array. If so, the Host info is returned and no further action is
// taken. If not, this method chooses a unique name for the Host and creates it on the array. Once the Host is created,
// it is placed in the Host Group used for nDVP volumes.
func (d ESeriesAPIDriver) EnsureHostForIQN(iqn string) (HostEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "EnsureHostForIQN",
			"Type":   "ESeriesAPIDriver",
			"iqn":    iqn,
		}
		log.WithFields(fields).Debug(">>>> EnsureHostForIQN")
		defer log.WithFields(fields).Debug("<<<< EnsureHostForIQN")
	}

	// See if there is already a host for the specified IQN
	host, err := d.GetHostForIQN(iqn)
	if err != nil {
		return HostEx{}, fmt.Errorf("Could not ensure host for IQN %s. %v", iqn, err)
	}

	// If we found a host, return it and leave well enough alone, since the host could have been defined
	// by nDVP or by an admin. Otherwise, we'll create a host for our purposes.
	if host.HostRef != "" {

		log.WithFields(log.Fields{
			"IQN":        iqn,
			"HostRef":    host.HostRef,
			"ClusterRef": host.ClusterRef,
		}).Debug("Host already exists for IQN.")

		return host, nil
	}

	// Pick a host name
	hostname := d.createNameForHost(iqn)

	// Ensure we have a group for the new host. If for any reason this fails, we'll create the host anyway with no group.
	hostGroup, err := d.EnsureHostGroup()
	if err != nil {
		log.Warn("Could not ensure host group for new host.")
	}

	// Create the new host in the group
	return d.CreateHost(hostname, iqn, d.config.HostType, hostGroup)
}

func (d ESeriesAPIDriver) createNameForHost(iqn string) string {

	// Get unique hostname suffix up to 10 chars, either the last part of the IQN or a random sequence
	var uniqueSuffix = utils.RandomString(10)
	index := strings.LastIndex(iqn, ":")
	if (index >= 0) && (len(iqn) > index+2) {
		uniqueSuffix = iqn[index+1:]
	}
	if len(uniqueSuffix) > 10 {
		uniqueSuffix = uniqueSuffix[:10]
	}

	// Pick a host name, incorporating the local hostname value if possible
	hostname, err := os.Hostname()
	if err != nil {
		hostname = utils.RandomString(MAX_NAME_LENGTH)
	}

	// Use as much of the hostname as will fit
	maxLength := MAX_NAME_LENGTH - (len(uniqueSuffix) + 1)
	if len(hostname) > maxLength {
		hostname = hostname[0:maxLength]
	}

	return hostname + "_" + uniqueSuffix
}

func (d ESeriesAPIDriver) createNameForPort(host string) string {

	suffix := "_port"
	hostname := host

	maxLength := MAX_NAME_LENGTH - len(suffix)
	if len(hostname) > maxLength {
		hostname = hostname[0:maxLength]
	}

	return hostname + suffix
}

// GetHostForIQN queries the Host objects on the array an returns one matching the supplied IQN. An empty struct is
// returned if a matching host is not found, so the caller should check for empty values in the result.
func (d ESeriesAPIDriver) GetHostForIQN(iqn string) (HostEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetHostForIQN",
			"Type":   "ESeriesAPIDriver",
			"iqn":    iqn,
		}
		log.WithFields(fields).Debug(">>>> GetHostForIQN")
		defer log.WithFields(fields).Debug("<<<< GetHostForIQN")
	}

	// Get hosts
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/hosts")
	if err != nil {
		return HostEx{}, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return HostEx{}, fmt.Errorf("Could not get hosts from array. Status code %d", response.StatusCode)
	}

	// Parse JSON data
	hosts := make([]HostEx, 0)
	if err := json.Unmarshal(responseBody, &hosts); err != nil {
		return HostEx{}, fmt.Errorf("Could not parse host data: %s. %v", string(responseBody), err)
	}

	// Find initiator with matching IQN
	for _, host := range hosts {
		for _, f := range host.Initiators {
			if f.NodeName.IoInterfaceType == "iscsi" && f.NodeName.IscsiNodeName == iqn {

				log.WithFields(log.Fields{
					"Name": host.Label,
					"IQN":  iqn,
				}).Debug("Found host.")

				return host, nil
			}
		}
	}

	// Nothing failed, so return an empty structure if we didn't find anything
	log.WithField("IQN", iqn).Debug("No host found.")
	return HostEx{}, nil
}

// CreateHost creates a Host on the array. If a HostGroup is specified, the Host is placed in that group.
func (d ESeriesAPIDriver) CreateHost(name string, iqn string, hostType string, hostGroup HostGroup) (HostEx, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":    "CreateHost",
			"Type":      "ESeriesAPIDriver",
			"name":      name,
			"iqn":       iqn,
			"hostType":  hostType,
			"hostGroup": hostGroup.Label,
		}
		log.WithFields(fields).Debug(">>>> CreateHost")
		defer log.WithFields(fields).Debug("<<<< CreateHost")
	}

	// Set up the host create request
	var request HostCreateRequest
	request.Name = name
	request.HostType.Index = d.getBestIndexForHostType(hostType)
	request.GroupID = hostGroup.ClusterRef
	request.Ports = make([]HostPort, 1)
	request.Ports[0].Label = d.createNameForPort(name)
	request.Ports[0].Port = iqn
	request.Ports[0].Type = "iscsi"

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return HostEx{}, fmt.Errorf("Could not marshal JSON request: %v. %v", request, err)
	}

	log.WithFields(log.Fields{
		"Name":  name,
		"Group": hostGroup.Label,
		"IQN":   iqn,
	}).Debug("Creating host.")

	// Create the host
	response, responseBody, err := d.InvokeAPI(jsonRequest, "POST", "/hosts")
	if err != nil {
		return HostEx{}, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		return HostEx{}, fmt.Errorf("Could not create host %s. Status code %d", name, response.StatusCode)
	}

	// Parse JSON data
	host := HostEx{}
	if err := json.Unmarshal(responseBody, &host); err != nil {
		return HostEx{}, fmt.Errorf("Could not parse host data: %s. %v", string(responseBody), err)
	}

	log.WithFields(log.Fields{
		"Name":          host.Label,
		"HostRef":       host.HostRef,
		"ClusterRef":    host.ClusterRef,
		"HostTypeIndex": host.HostTypeIndex,
	}).Debug("Created host.")

	return host, nil
}

func (d ESeriesAPIDriver) getBestIndexForHostType(hostType string) int {

	hostTypeIndex := -1

	// Try the mapped values first
	_, ok := HOST_TYPES[hostType]
	if ok {
		hostTypeIndex, _ = d.getIndexForHostType(HOST_TYPES[hostType])
	}

	// If not found, try matching the E-series host type codes directly
	if hostTypeIndex == -1 {
		hostTypeIndex, _ = d.getIndexForHostType(hostType)
	}

	// If still not found, fall back to standard Linux DM-MPIO multipath driver
	if hostTypeIndex == -1 {
		hostTypeIndex, _ = d.getIndexForHostType("LnxALUA")
	}

	// Failing that, use index 0, which should be the factory default
	if hostTypeIndex == -1 {
		hostTypeIndex = 0
	}

	return hostTypeIndex
}

// getIndexForHostType queries the array for a host type matching the specified value. If found, it returns the
// index by which the type is known on the array. If not found, it returns -1.
func (d ESeriesAPIDriver) getIndexForHostType(hostTypeCode string) (int, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":       "getIndexForHostType",
			"Type":         "ESeriesAPIDriver",
			"hostTypeCode": hostTypeCode,
		}
		log.WithFields(fields).Debug(">>>> getIndexForHostType")
		defer log.WithFields(fields).Debug("<<<< getIndexForHostType")
	}

	// Get host types
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/host-types")
	if err != nil {
		return -1, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("Could not get host types from array. Status code %d", response.StatusCode)
	}

	// Parse JSON data
	hostTypes := make([]HostType, 0)
	if err := json.Unmarshal(responseBody, &hostTypes); err != nil {
		return -1, fmt.Errorf("Could not parse host type data: %s. %v", string(responseBody), err)
	}

	// Find host type with matching code
	for _, hostType := range hostTypes {
		if hostType.Code == hostTypeCode {

			log.WithFields(log.Fields{
				"Name":  hostType.Name,
				"Index": hostType.Index,
				"Code":  hostType.Code,
			}).Debug("Host type found.")

			return hostType.Index, nil
		}
	}

	log.WithField("Code", hostTypeCode).Debug("Host type not found.")
	return -1, nil
}

// EnsureHostGroup ensures that an E-series HostGroup exists to contain all Host objects created by the nDVP E-series driver.
// The group name is taken from the config structure. If the group exists, the group structure is returned and no further
// action is taken. If not, this method creates the group and returns the resulting group structure.
func (d ESeriesAPIDriver) EnsureHostGroup() (HostGroup, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "EnsureHostGroup",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> EnsureHostGroup")
		defer log.WithFields(fields).Debug("<<<< EnsureHostGroup")
	}

	hostGroupName := d.config.AccessGroup
	if len(hostGroupName) > MAX_NAME_LENGTH {
		hostGroupName = hostGroupName[:MAX_NAME_LENGTH]
	}

	// Get the group with the preconfigured name
	hostGroup, err := d.GetHostGroup(hostGroupName)
	if err != nil {
		return HostGroup{}, fmt.Errorf("Could not ensure host group %s. %v", hostGroupName, err)
	}

	// Group found, so use it for host creation
	if hostGroup.ClusterRef != "" {
		log.WithFields(log.Fields{"Name": hostGroup.Label}).Debug("Host group found.")
		return hostGroup, nil
	}

	// Create the group
	hostGroup, err = d.CreateHostGroup(hostGroupName)
	if err != nil {
		return HostGroup{}, fmt.Errorf("Could not create host group %s. %v", hostGroupName, err)
	}

	log.WithFields(log.Fields{
		"Name":       hostGroup.Label,
		"ClusterRef": hostGroup.ClusterRef,
	}).Debug("Host group created.")

	return hostGroup, nil
}

// GetHostGroup returns an E-series HostGroup structure with the specified name. If no matching group is found, an
// empty structure is returned, so the caller should check for empty values in the result.
func (d ESeriesAPIDriver) GetHostGroup(name string) (HostGroup, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetHostGroup",
			"Type":   "ESeriesAPIDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> GetHostGroup")
		defer log.WithFields(fields).Debug("<<<< GetHostGroup")
	}

	// Get host groups
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/host-groups")
	if err != nil {
		return HostGroup{}, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return HostGroup{}, fmt.Errorf("Could not get host groups from array. Status code %d", response.StatusCode)
	}

	// Parse JSON data
	hostGroups := make([]HostGroup, 0)
	if err := json.Unmarshal(responseBody, &hostGroups); err != nil {
		return HostGroup{}, fmt.Errorf("Could not parse host group data: %s. %v", string(responseBody), err)
	}

	for _, hostGroup := range hostGroups {
		if hostGroup.Label == name {
			return hostGroup, nil
		}
	}

	// Nothing failed, so return an empty structure if we didn't find anything
	return HostGroup{}, nil
}

// CreateHostGroup creates an E-series HostGroup object with the specified name and returns the resulting HostGroup structure.
func (d ESeriesAPIDriver) CreateHostGroup(name string) (HostGroup, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "CreateHostGroup",
			"Type":   "ESeriesAPIDriver",
			"name":   name,
		}
		log.WithFields(fields).Debug(">>>> CreateHostGroup")
		defer log.WithFields(fields).Debug("<<<< CreateHostGroup")
	}

	// Set up the host group create request
	request := HostGroupCreateRequest{
		Name:  name,
		Hosts: make([]string, 0),
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return HostGroup{}, fmt.Errorf("Could not marshal JSON request: %v. %v", request, err)
	}

	// Create the host group
	response, responseBody, err := d.InvokeAPI(jsonRequest, "POST", "/host-groups")
	if err != nil {
		return HostGroup{}, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		return HostGroup{}, fmt.Errorf("Could not create host group %s. Status code %d", name, response.StatusCode)
	}

	// Parse JSON data
	hostGroup := HostGroup{}
	if err := json.Unmarshal(responseBody, &hostGroup); err != nil {
		return HostGroup{}, fmt.Errorf("Could not parse host data: %s. %v", string(responseBody), err)
	}

	return hostGroup, nil
}

// MapVolume maps a volume to the specified host and returns the resulting LUN mapping. If the volume is already mapped to the
// specified host, either directly or to the containing host group, no action is taken. If the volume is mapped to a different host,
// the method returns an error. Note that if the host is in a group, the volume will actually be mapped to the group instead of the
// individual host.
func (d ESeriesAPIDriver) MapVolume(volume VolumeEx, host HostEx) (LUNMapping, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "MapVolume",
			"Type":       "ESeriesAPIDriver",
			"volumeName": volume.Label,
			"hostName":   host.Label,
		}
		log.WithFields(fields).Debug(">>>> MapVolume")
		defer log.WithFields(fields).Debug("<<<< MapVolume")
	}

	if !volume.IsMapped {

		// Volume is not already mapped, so map it now
		return d.mapVolume(volume, host)

	} else {

		// Volume is mapped, so check if it is mapped to this host/group
		mappedToHost, mapping := d.volumeIsMappedToHost(volume, host)

		if mappedToHost {

			// Mapped here, so nothing to do
			log.WithFields(log.Fields{
				"Name": volume.Label,
				"Host": host.Label,
			}).Debug("Volume already mapped to host.")

			return mapping, nil

		} else {

			// Mapped elsewhere, so return an error
			return LUNMapping{}, fmt.Errorf("Volume %s is already mapped to a different host or host group", volume.Label)
		}
	}
}

// mapVolume maps a volume to a host with no checks for an existing mapping. If the host is in a host group, the volume is
// mapped to the group instead. The resulting mapping structure is returned.
func (d ESeriesAPIDriver) mapVolume(volume VolumeEx, host HostEx) (LUNMapping, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "mapVolume",
			"Type":       "ESeriesAPIDriver",
			"volumeName": volume.Label,
			"hostName":   host.Label,
		}
		log.WithFields(fields).Debug(">>>> mapVolume")
		defer log.WithFields(fields).Debug("<<<< mapVolume")
	}

	// Map to host group if available, else a host
	targetID := host.HostRef
	if d.IsRefValid(host.ClusterRef) {
		targetID = host.ClusterRef
	}

	// Create a map request. The API proxy will pick a non-zero LUN number automatically.
	request := VolumeMappingCreateRequest{
		MappableObjectID: volume.VolumeRef,
		TargetID:         targetID,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return LUNMapping{}, fmt.Errorf("Could not marshal JSON request: %v. %v", request, err)
	}

	// Create the mapping
	response, responseBody, err := d.InvokeAPI(jsonRequest, "POST", "/volume-mappings")
	if err != nil {
		return LUNMapping{}, fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return LUNMapping{}, fmt.Errorf("Could not map volume %s. Status code %d", volume.Label, response.StatusCode)
	}

	// Parse JSON data
	mapping := LUNMapping{}
	if err := json.Unmarshal(responseBody, &mapping); err != nil {
		return LUNMapping{}, fmt.Errorf("Could not parse volume mapping data: %s. %v", string(responseBody), err)
	}

	log.WithFields(log.Fields{
		"Name":      volume.Label,
		"VolumeRef": volume.VolumeRef,
		"MapRef":    mapping.MapRef,
		"Type":      mapping.Type,
		"LunNumber": mapping.LunNumber,
	}).Debug("Volume mapped.")

	return mapping, nil
}

// volumeIsMappedToHost checks whether a volume is mapped to the specified host (or containing host group). If the mapping
// exists, the method returns true with the associated mapping structure. If no mapping exists, or if the volume is mapped
// elsewhere, the method returns false with an empty structure.
func (d ESeriesAPIDriver) volumeIsMappedToHost(volume VolumeEx, host HostEx) (bool, LUNMapping) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method":     "volumeIsMappedToHost",
			"Type":       "ESeriesAPIDriver",
			"volumeName": volume.Label,
			"hostName":   host.Label,
		}
		log.WithFields(fields).Debug(">>>> volumeIsMappedToHost")
		defer log.WithFields(fields).Debug("<<<< volumeIsMappedToHost")
	}

	// Quick check to skip everything else
	if !volume.IsMapped {
		log.WithField("Name", volume.Label).Debug("Volume is not mapped.")
		return false, LUNMapping{}
	}

	// E-series only supports a single mapping per volume
	if len(volume.Mappings) == 0 {
		log.WithField("Name", volume.Label).Debug("Volume has no mappings.")
		return false, LUNMapping{}
	}
	mapping := volume.Mappings[0]

	// Double check we're looking at the right volume
	if mapping.VolumeRef != volume.VolumeRef {
		return false, LUNMapping{}
	}

	// Match either a host or a host group
	switch mapping.Type {

	case HOST_GROUP_MAPPING_TYPE: // "cluster"

		if mapping.MapRef == host.ClusterRef {

			log.WithFields(log.Fields{
				"volumeName": volume.Label,
				"hostName":   host.Label,
			}).Debug("Volume is mapped to the host's enclosing group.")

			return true, mapping
		}

	case HOST_MAPPING_TYPE: // "host"

		if mapping.MapRef == host.HostRef {

			log.WithFields(log.Fields{
				"volumeName": volume.Label,
				"hostName":   host.Label,
			}).Debug("Volume is mapped to the host.")

			return true, mapping
		}
	}

	log.WithFields(log.Fields{
		"volumeName": volume.Label,
		"hostName":   host.Label,
	}).Debug("Volume is mapped to a different host or group.")

	return false, LUNMapping{}
}

// UnmapVolume removes a mapping from the specified volume. If no map exists, no action is taken.
func (d ESeriesAPIDriver) UnmapVolume(volume VolumeEx) error {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "UnmapVolume",
			"Type":   "ESeriesAPIDriver",
			"volume": volume.Label,
		}
		log.WithFields(fields).Debug(">>>> UnmapVolume")
		defer log.WithFields(fields).Debug("<<<< UnmapVolume")
	}

	// If volume isn't mapped, there's nothing to do
	if !volume.IsMapped || len(volume.Mappings) == 0 {

		log.WithFields(log.Fields{
			"Name":      volume.Label,
			"VolumeRef": volume.VolumeRef,
		}).Warn("Volume unmap requested, but volume is not mapped.")

		return nil
	}
	mapping := volume.Mappings[0]

	// Remove this volume mapping from storage array
	response, _, err := d.InvokeAPI(nil, "DELETE", "/volume-mappings/"+mapping.LunMappingRef)
	if err != nil {
		return fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Could not unmap volume %s. Status code %d", volume.Label, response.StatusCode)
	}

	log.WithFields(log.Fields{
		"Name":      volume.Label,
		"VolumeRef": volume.VolumeRef,
		"MapRef":    mapping.MapRef,
		"Type":      mapping.Type,
		"LunNumber": mapping.LunNumber,
	}).Debug("Volume unmapped.")

	return nil
}

// GetTargetIqn returns the IQN for the array.
func (d *ESeriesAPIDriver) GetTargetIQN() (string, error) {

	if d.config.DebugTraceFlags["method"] {
		fields := log.Fields{
			"Method": "GetTargetIqn",
			"Type":   "ESeriesAPIDriver",
		}
		log.WithFields(fields).Debug(">>>> GetTargetIqn")
		defer log.WithFields(fields).Debug("<<<< GetTargetIqn")
	}

	// Query iSCSI target settings
	response, responseBody, err := d.InvokeAPI(nil, "GET", "/iscsi/target-settings")
	if err != nil {
		return "", fmt.Errorf("API invocation failed. %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Could not read iSCSI settings. Status code %d", response.StatusCode)
	}

	var settings IscsiTargetSettings
	err = json.Unmarshal(responseBody, &settings)
	if err != nil {
		return "", fmt.Errorf("Could not parse iSCSI settings data: %s. %v", string(responseBody), err)
	}

	log.WithFields(log.Fields{
		"IargetIQN": settings.NodeName.IscsiNodeName,
	}).Debug("Got target iSCSI node name.")

	return settings.NodeName.IscsiNodeName, nil
}

// isRefValid checks whether the supplied string is a valid E-series object reference as used by its REST API.
// Ref values are strings of all numerical digits that aren't all zeros (i.e. the null ref).
func (d ESeriesAPIDriver) IsRefValid(ref string) bool {

	switch ref {
	case "", NULL_REF:
		return false
	default:
		return true
	}
}

// getErrorFromHTTPResponse converts error information from some E-series API responses into GoLang error objects that
// embed the additional error text.
func (d ESeriesAPIDriver) getErrorFromHTTPResponse(response *http.Response, responseBody []byte) error {

	if response.StatusCode == http.StatusNotFound || response.StatusCode == HTTP_UNPROCESSABLE_ENTITY {

		// Parse JSON error data
		responseData := CallResponseError{}
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			return fmt.Errorf("Could not parse API error response: %s. %v", string(responseBody), err)
		}

		return fmt.Errorf("API failed. Status code: %d. Error: %s Localized: %s",
			response.StatusCode, responseData.ErrorMsg, responseData.LocalizedMsg)
	} else {

		// Other error
		return fmt.Errorf("API failed. Status code: %d", response.StatusCode)
	}
}
