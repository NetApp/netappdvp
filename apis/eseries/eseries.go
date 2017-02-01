// Copyright 2016 NetApp, Inc. All Rights Reserved.

package eseries

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/utils"
)

const maxNameLength int = 30

// VolumeInfo hold all the information about a constructed volume on E-Series array and Docker Host Mapping
type VolumeInfo struct {
	VolumeGroupRef string
	VolumeRef      string

	VolumeSize  int64
	SegmentSize int
	UnitSize    string

	MediaType    string
	SecureVolume bool

	IsVolumeMapped bool
	LunMappingRef  string
	LunNumber      int
}

// DriverConfig holds the configuration data for Driver objects
type DriverConfig struct {
	//Web Proxy Services Info
	WebProxyHostname  string
	WebProxyPort      string
	WebProxyUseHTTP   bool
	WebProxyVerifyTLS bool
	Username          string
	Password          string

	//Array Info
	ControllerA     string
	ControllerB     string
	PasswordArray   string
	ArrayRegistered bool

	//Host Connectivity
	HostDataIP string //for iSCSI with multipathing this can be either IP on host

	//Internal Config Variables
	ArrayID string //Unique ID for array once added to web proxy services
	Volumes map[string]*VolumeInfo

	//Storage protocol of the driver (iSCSI, FC, etc)
	Protocol string

	DriverName string
	Version    int
}

// Driver is the object to use for interacting with the Array
type Driver struct {
	config *DriverConfig
	m      *sync.Mutex
}

// NewDriver is a factory method for creating a new instance
func NewDriver(config DriverConfig) *Driver {
	d := &Driver{
		config: &config,
		m:      &sync.Mutex{},
	}

	//Clear out internal config variables
	d.config.ArrayID = ""
	d.config.Volumes = make(map[string]*VolumeInfo)

	volumeTags = []VolumeTag{
		{Key: "API", Value: "netappdvp-" + strconv.Itoa(config.Version)},
		{Key: "eBI", Value: "Containers-Docker"},
		{Key: "IF", Value: d.config.Protocol},
		{Key: "netappdvp", Value: config.DriverName},
	}

	return d
}

var volumeTags []VolumeTag

// SendMsg sends the marshaled json byte array to the web services proxy
func (d Driver) SendMsg(data []byte, httpMethod string, msgType string) (*http.Response, error) {

	if data == nil && d.config.ArrayID == "" {
		return nil, fmt.Errorf("Data is nil and no ArrayID set!")
	}

	log.Debugf("Sending data to web services proxy @ '%s' json: \n%s", d.config.WebProxyHostname, string(data))

	// Default to secure connection
	addressPrefix := "https"
	addressPort := "8443"

	if d.config.WebProxyUseHTTP {
		addressPrefix = "http"
		addressPort = "8080"
	}

	// Allow port override
	if d.config.WebProxyPort != "" {
		log.Debugf("Setting web services proxy port to %s", d.config.WebProxyPort)
		addressPort = d.config.WebProxyPort
	}

	//Set up address to web services proxy
	url := addressPrefix + "://" + d.config.WebProxyHostname + ":" + addressPort + "/devmgr/v2/storage-systems/" + d.config.ArrayID + msgType
	log.Debugf("URL:> %s", url)

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(d.config.Username, d.config.Password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !d.config.WebProxyVerifyTLS,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	//At this point either the resp or the err could be nil
	if resp != nil {
		log.Debugf("Response Status: %s", resp.Status)
		log.Debugf("Response Headers: %s", resp.Header)
	}

	if err != nil {
		log.Warnf("Error communicating with Web Services: %v!", err)
	}

	return resp, err
}

func (d Driver) init() {
}

func (d Driver) populateVolumeCache(name string) (err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	// If a volume is specified, only re-populate if it doesn't already exist in the cache.
	if len(name) > 0 {
		//First check if volume is already in persistant map
		if _, isPresent := d.config.Volumes[name]; isPresent {
			return nil
		}
	}

	//If not in map then we need to query volumes on array
	resp, err := d.SendMsg(nil, "GET", "/volumes")
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay {
		return fmt.Errorf("ESeriesStorageDriver::populateVolumeCache - GET to obtain volume failed! StatusCode=%v Status=%s", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	//Next need to demarshal json data
	responseJSON := make([]MsgVolumeExResponse, 0)
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		panic(err)
	}

	for _, e := range responseJSON {
		if len(name) > 0 {
			if e.Label != name {
				continue
			}
		}

		//Create a VolumeInfo structure and add it to the map
		var tmpVolumeInfo VolumeInfo
		tmpVolumeInfo.VolumeGroupRef = e.VolumeGroupRef
		tmpVolumeInfo.VolumeRef = e.VolumeRef

		//Convert size of volume from string to int64
		tmpLunSize, atoiErr := strconv.ParseInt(e.VolumeSize, 10, 0)
		if atoiErr != nil {
			fmt.Errorf("Cannot convert size to bytes: %v error: %v", e.VolumeSize, atoiErr)
			return atoiErr
		}

		tmpVolumeInfo.VolumeSize = tmpLunSize
		tmpVolumeInfo.SegmentSize = e.SegmentSize
		tmpVolumeInfo.UnitSize = "b" //bytes

		//Need to figure out mediaType (whether this volume belongs to hdd or ssd group)
		volumeGroupRef, error := d.VerifyVolumePools("hdd", "1m") //1 megabyte is just a small unit to use while figuring out which media type this volume belongs to
		if error != nil {
			return error
		}

		if e.VolumeGroupRef == volumeGroupRef {
			//Found it! It is a hdd volume group.
			tmpVolumeInfo.MediaType = "hdd"
		} else {
			//Not a part of hdd volume group so see if it is part of ssd volume group
			volumeGroupRef1, error1 := d.VerifyVolumePools("ssd", "1m") //1 megabyte is just a small unit to use while figuring out which media type this volume belongs to
			if error != nil {
				return error1
			}

			if e.VolumeGroupRef == volumeGroupRef1 {
				//Found it! It is part of ssd volume group.
				tmpVolumeInfo.MediaType = "ssd"
			} else {
				//It isn't part of ssd nor hdd volume group!
				continue
			}
		}

		tmpVolumeInfo.SecureVolume = false //TODO: add this capability for FDE drives
		tmpVolumeInfo.IsVolumeMapped = e.IsMapped

		for j, f := range e.ListOfMappings {
			log.Debugf("%v) Volume with name %s has mapping reference %s", j, e.Label, f.LunMappingRef)
			tmpVolumeInfo.LunMappingRef = f.LunMappingRef //TODO - what if there are multiple mappings? Is this even possible outside 'Default Group'?
			tmpVolumeInfo.LunNumber = f.LunNumber
		}

		//Add it to map
		d.config.Volumes[e.Label] = &tmpVolumeInfo

		// If we were just looking for one, we just found it. Bail out.
		if len(name) > 0 {
			break
		}
	}

	return nil
}

// Connect to the array's Web Services Proxy
func (d Driver) Connect() (response string, err error) {

	//Send a login/connect request for array to web services proxy
	msgConnect := MsgConnect{[]string{d.config.ControllerA, d.config.ControllerB}, d.config.PasswordArray}

	jsonConnect, err := json.Marshal(msgConnect)
	if err != nil {
		return "", fmt.Errorf("Error defining JSON body: %v", err)
	}

	log.Debugf("jsonConnect=%s", string(jsonConnect))

	//Send off the message
	resp, err := d.SendMsg(jsonConnect, "POST", "")

	if err != nil {
		return "", fmt.Errorf("Error logging into the Web Services Proxy: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseSuccess && resp.StatusCode != GenericResponseOkay {
		return "", fmt.Errorf("Couldn't add storage array to web services proxy!")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	//Next need to demarshal json data
	responseData := MsgConnectResponse{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("Unable to deserialize login JSON response: %v", err)
	}

	if responseData.ArrayID == "" {
		return "", fmt.Errorf("We received an invalid ArrayID!")
	}

	d.config.ArrayID = responseData.ArrayID
	d.config.ArrayRegistered = responseData.AlreadyExists

	log.Debugf("ArrayID=%s alreadyRegistered=%v", d.config.ArrayID, d.config.ArrayRegistered)

	return d.config.ArrayID, nil
}

func (d Driver) VerifyVolumePools(mediaType string, size string) (VolumeGroupRef string, err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	//Do the GET to obtain volume pools
	resp, err := d.SendMsg(nil, "GET", "/storage-pools")
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay {
		return "", fmt.Errorf("ESeriesStorageDriver::VerifyVolumePools - GET to obtain volume pools failed! StatusCode=%v Status=%s", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	//Next need to demarshal json data
	responseJSON := make([]VolumeGroupExResponse, 0)
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		panic(err)
	}

	//Create search volume group label
	var volumeGroupLabel string = "netappdvp_"
	if mediaType != "ssd" && mediaType != "hdd" {
		return "", fmt.Errorf("ESeriesStorageDriver::VerifyVolumePools - mediaType specified is invalid! mediaType=%s", mediaType)
	}

	volumeGroupLabel += mediaType

	//Search for correct volume group label
	var volumeGroupRef string = ""
	var volumeFreeSpace string = "0"
	for i, e := range responseJSON {
		log.Debugf("%v) label=%s freeSpace=%s", i, e.VolumeLabel, e.FreeSpace)

		if e.VolumeLabel == volumeGroupLabel {
			volumeGroupRef = e.VolumeGroupRef
			volumeFreeSpace = e.FreeSpace
		}
	}

	if volumeGroupRef == "" {
		return "", fmt.Errorf("ESeriesStorageDriver::VerifyVolumePools - correct volume group not found for mediaType=%s!", mediaType)
	}

	//Verify volume group has enough space
	convertedSize, convertErr := utils.ConvertSizeToBytes64(size)
	if convertErr != nil {
		fmt.Errorf("Cannot convert size to bytes: %v error: %v", size, convertErr)
		return "", convertErr
	}
	lunSize, atoiErr := strconv.ParseInt(convertedSize, 10, 0)
	if atoiErr != nil {
		fmt.Errorf("Cannot convert size to bytes: %v error: %v", size, atoiErr)
		return "", atoiErr
	}

	if convertedVolumeFreeSpace, _ := strconv.ParseInt(volumeFreeSpace, 10, 0); lunSize > convertedVolumeFreeSpace {
		return "", fmt.Errorf("ESeriesStorageDriver::VerifyVolumePools - volume group doesn't have enough space! lunSize=%v > volumeFreeSpace=%v!", lunSize, volumeFreeSpace)
	}

	return volumeGroupRef, nil
}

func (d Driver) GetVolumeList() (vols []string, err error) {
	// TODO: Right now we re-build the whole list every time; not sure how painful this call is. Should
	// look into finding ways to invalidate the cache and instituting a TTL, although there will be
	// issues when this driver is installed on multiple hosts if we do that.
	err = d.populateVolumeCache("")
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve the list of volumes: %v", err)
	}

	for vol := range d.config.Volumes {
		vols = append(vols, vol)
	}

	return vols, nil
}

func (d Driver) VerifyVolumeExists(name string) (err error) {
	d.populateVolumeCache(name)

	if _, isPresent := d.config.Volumes[name]; isPresent {
		return nil
	}

	return fmt.Errorf("ESeriesStorageDriver::VerifyVolumeExists - volume with name %s not found on array! Are you sure you created a volume with this name?", name)
}

func (d Driver) IsVolumeAlreadyMappedToHost(name string, hostRef string) (isMapped bool, lunNumber int, err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	resp, err := d.SendMsg(nil, "GET", "/volume-mappings")
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay {
		return false, -1, fmt.Errorf("ESeriesStorageDriver::IsVolumeAlreadyMappedToHost - GET to obtain volume mappings failed! StatusCode=%v Status=%s", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	//Next need to demarshal json data
	responseJSON := make([]LUNMapping, 0)
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		panic(err)
	}

	//Need to find the correct volumeRef for the name argument and also verify that if is mapped already that it is indeed mapped to our host and not some other host
	var foundVolumeMapping bool = false

	//Make sure volume is already in persistant map
	tmpVolumeInfo, isPresent := d.config.Volumes[name]
	if !isPresent {
		panic("volume is not already apart of map! Check to see if you indeed created the docker volume with this name.")
	}

	//tmpVolumeInfo.VolumeRef
	for i, e := range responseJSON {
		log.Debugf("%v) Found host mapping with volumeRef=%s mapped to LUN number %v with hostRef=%s", i, e.VolumeRef, e.LunNumber, e.HostRef)

		if e.VolumeRef == tmpVolumeInfo.VolumeRef {
			//We found our volume and it is indeed mapped, but is it mapped to the correct host?
			foundVolumeMapping = true

			if e.HostRef == hostRef {
				//Yes, it is mapped to proper host
				return true, e.LunNumber, nil
			} else {
				//No, it is mapped to different host!
				return false, e.LunNumber, fmt.Errorf("ESeriesStorageDriver::IsVolumeAlreadyMappedToHost - found mapped volume with name %s but it is mapped to different host (%s) rather than requested host (%s)", name, e.HostRef, hostRef)
			}
		}
	}

	if foundVolumeMapping {
		panic("Unreachable code path")
	}

	//The volume is not mapped to any host
	return false, -1, nil
}

func (d Driver) CreateVolume(name string, volumeGroupRef string, size string, mediaType string) (volumeRef string, err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	// Ensure that we do not exceed the maximum allowed volume length
	if len(name) > maxNameLength {
		return "", fmt.Errorf("The volume name of %v exceeds the maximum allowed length of %d characters",
			name, maxNameLength)
	}

	//Lets create a volume message structure
	var msgCreateVolume MsgVolumeEx
	msgCreateVolume.VolumeGroupRef = volumeGroupRef
	msgCreateVolume.Name = name
	msgCreateVolume.SizeUnit = "kb" //bytes, b, kb, mb, gb, tb, pb, eb, zb, yb
	msgCreateVolume.SegmentSize = 128
	msgCreateVolume.VolumeTags = volumeTags

	//Convert size string to int64
	convertedSize, convertErr := utils.ConvertSizeToBytes64(size)
	if convertErr != nil {
		return "", fmt.Errorf("Cannot convert size to bytes: %v error: %v", size, convertErr)
	}
	lunSize, atoiErr := strconv.ParseInt(convertedSize, 10, 0)
	if atoiErr != nil {
		return "", fmt.Errorf("Cannot convert size to bytes: %v error: %v", size, atoiErr)
	}

	//Set lun size in kilobytes
	msgCreateVolume.Size = int(lunSize / 1024) //the json request requires lunSize to be an int not an int64 so we are passing it as an int but in kilobytes

	jsonCreateVolume, err := json.Marshal(msgCreateVolume)
	if err != nil {
		panic(err)
	}

	log.Debugf("jsonCreateVolume=%s", string(jsonCreateVolume))

	//Send off the message
	resp, err := d.SendMsg(jsonCreateVolume, "POST", "/volumes")
	defer resp.Body.Close()

	if resp.StatusCode == GenericResponseOkay {
		//Success!

		body, _ := ioutil.ReadAll(resp.Body)
		log.Debugf("response Body:\n%s", string(body))

		//Next need to demarshal json data
		responseData := MsgVolumeExResponse{}
		if err := json.Unmarshal(body, &responseData); err != nil {
			panic(err)
		}

		log.Debugf("Label=%s VolumeRef=%s", responseData.Label, responseData.VolumeRef)

		//Create a VolumeInfo structure and put it into the map
		var tmpVolumeInfo VolumeInfo
		tmpVolumeInfo.VolumeGroupRef = volumeGroupRef
		tmpVolumeInfo.VolumeRef = responseData.VolumeRef

		tmpVolumeInfo.VolumeSize = lunSize
		tmpVolumeInfo.SegmentSize = msgCreateVolume.SegmentSize * 1024 //convert from kilobytes to bytes
		tmpVolumeInfo.UnitSize = "b"                                   //bytes

		//Sanity check
		if tmpVolumeInfo.SegmentSize != responseData.SegmentSize {
			panic("Segment size specified in request doesn't equal returned volume segment size!")
		}

		tmpVolumeInfo.MediaType = mediaType
		tmpVolumeInfo.SecureVolume = false //TODO: add this capability for FDE drives
		tmpVolumeInfo.IsVolumeMapped = false
		tmpVolumeInfo.LunMappingRef = ""
		tmpVolumeInfo.LunNumber = -1

		//Add it to map
		d.config.Volumes[name] = &tmpVolumeInfo

		return responseData.VolumeRef, nil
	} else if resp.StatusCode == GenericResponseNotFound || resp.StatusCode == GenericResponseMalformed {
		//Known Error!

		body, _ := ioutil.ReadAll(resp.Body)
		log.Debugf("response Body:\n%s", string(body))

		//Next need to demarshal json data
		responseData := CallResponseError{}
		if err := json.Unmarshal(body, &responseData); err != nil {
			panic(err)
		}

		return "", fmt.Errorf("Error - not found code - ErrorMsg=%s LocalizedMsg=%s", responseData.ErrorMsg, responseData.LocalizedMsg)
	} else {

		return "", fmt.Errorf("Unknown error code - StatusCode=%v!", resp.StatusCode)
	}

	return "", fmt.Errorf("Unreachable Code Path!")
}

func (d Driver) VerifyHostIQN(iqn string) (hostRef string, err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	//Do a GET to obtain hosts on array
	resp, err := d.SendMsg(nil, "GET", "/hosts")
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay {
		panic("Couldn't obtain storage pools!")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	//Next need to demarshal json data
	responseJSON := make([]HostExResponse, 0)
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		panic(err)
	}

	var retHostRef string = ""
	for i, e := range responseJSON {
		log.Debugf("%v) HostRef=%s Label=%s", i, e.HostRef, e.Label)

		for j, f := range e.Initiators {
			log.Debugf("	%v) Host_Label=%s interface=%s iqn=%s", j, f.Label, f.NodeName.IoInterfaceType, f.NodeName.IscsiNodeName)

			if f.NodeName.IoInterfaceType == "iscsi" && f.NodeName.IscsiNodeName == iqn {
				retHostRef = e.HostRef
			}
		}
	}

	//Make sure we found the correct hostRef
	if retHostRef == "" {
		return "", fmt.Errorf("Host reference not found on array for host IQN %s!", iqn)
	}

	return retHostRef, nil
}

func (d Driver) MapVolume(name string, hostRef string) (lunNumber int, err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	//Look up volumeRef in persistant map, and error out if it is not found
	tmpVolumeInfo, isPresent := d.config.Volumes[name]
	if !isPresent {
		return -1, fmt.Errorf("name (%s) wasn't found in persistant map! Have you created a netappdvp volume group yet with that name?", name)
	}

	//Lets create a volume message structure
	var msgMapVolume VolumeMappingCreateRequest
	msgMapVolume.MappableObjectId = tmpVolumeInfo.VolumeRef
	msgMapVolume.TargetID = hostRef
	//msgMapVolume.LunNumber = 20	<--- optional parameter. Just let array choose the LUN #, and return it on response

	jsonMapVolume, err := json.Marshal(msgMapVolume)
	if err != nil {
		panic(err)
	}

	log.Debugf("jsonMapVolume=%s", string(jsonMapVolume))

	//Send off the message
	resp, err := d.SendMsg(jsonMapVolume, "POST", "/volume-mappings")
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay {
		return -1, fmt.Errorf("Error occured while mapping volume! ReturnCode=%v name=%s volumeRef=%s hostRef=%s", resp.StatusCode, name, tmpVolumeInfo.VolumeRef, hostRef)
	}

	//Got back success code
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("response Body:\n%s", string(body))

	//Next need to demarshal json data
	responseData := LUNMapping{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		panic(err)
	}

	log.Debugf("LunMappingRef=%s LunNumber=%v VolumeRef=%s", responseData.LunMappingRef, responseData.LunNumber, responseData.VolumeRef)

	//Sanity check to verify that returned volumeRef matches up with the volume we sent
	if tmpVolumeInfo.VolumeRef != responseData.VolumeRef {
		return -1, fmt.Errorf("Extremely odd case where the returned volumeRef (%s) doesn't equal the volumeRef (%s) we sent to array to map!", responseData.VolumeRef, tmpVolumeInfo.VolumeRef)
	}

	//Update volumeInfo map for this volume that we are not mapped as well as setting the LunMappingRef and LUN #
	if tmpVolumeInfo.IsVolumeMapped {
		panic("Volume is already mapped!")
	}

	tmpVolumeInfo.IsVolumeMapped = true
	tmpVolumeInfo.LunMappingRef = responseData.LunMappingRef
	tmpVolumeInfo.LunNumber = responseData.LunNumber

	return responseData.LunNumber, nil
}

func (d Driver) UnmapVolume(name string) (err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	//Need to lookup lunMappingRef for this volume from the volumeInfo map and do some sanity checking
	tmpVolumeInfo, isPresent := d.config.Volumes[name]
	if !isPresent {
		//If we are in this unmap function this volume better be in the volumeInfo map!
		panic("ERROR - volume wasn't found in volumeInfo map!")
	}

	//Sanity checks to make sure volume is mapped and have a LunMappingRef
	if !tmpVolumeInfo.IsVolumeMapped {
		//We are in unmap and the volume isn't even mapped!
		panic("ERROR - volume was found on array but it isn't mapped!")
	}

	if tmpVolumeInfo.LunMappingRef == "" {
		//Note: if this volume is indeed mapped it should have had its LunMappingRef filled either inside of CreateVolume function or inside of VerifyVolumeExists function!
		panic("ERROR - LunMappingRef is empty!")
	}

	//Send a DELETE to remove this LUN mapping from storage array
	resp, err := d.SendMsg(nil, "DELETE", "/volume-mappings/"+tmpVolumeInfo.LunMappingRef)
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay && resp.StatusCode != GenericResponseNoContent {
		return fmt.Errorf("Error occured while trying to remove LUN mapping for volume %s! Error Code (%v) and Status=%s", name, resp.StatusCode, resp.Status)
	}

	if resp.StatusCode != GenericResponseNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Debugf("response Body:\n%s", string(body))

		//Next need to demarshal json data
		responseData := LUNMapping{}
		if err := json.Unmarshal(body, &responseData); err != nil {
			panic(err)
		}

		//Sanity check to verify we remove LUN mapping for correct volume
		if responseData.LunMappingRef != tmpVolumeInfo.LunMappingRef {
			return fmt.Errorf("Error occured while trying to remove LUN mapping for volume %s! Very odd - responseData.LunMappingRef (%s) != tmpVolumeInfo.LunMappingRef (%s)...", name, responseData.LunMappingRef, tmpVolumeInfo.LunMappingRef)
		}
	}

	//Alter the state of the volumeInfo mapping to reflect this volume is no longer mapped
	tmpVolumeInfo.IsVolumeMapped = false
	tmpVolumeInfo.LunMappingRef = ""
	tmpVolumeInfo.LunNumber = -1

	//Return success!
	return nil
}

func (d Driver) DestroyVolume(name string) (err error) {

	//Verify we have a valid array id
	if d.config.ArrayID == "" {
		panic("ArrayID is invalid!")
	}

	//Make sure this volume is found in our volumeInfo map
	tmpVolumeInfo, isPresent := d.config.Volumes[name]
	if !isPresent {
		//If we are in this destroy function this volume better be in the volumeInfo map!
		panic("ERROR - volume wasn't found in volumeInfo map!")
	}

	//Sanity checks to make sure volume in map has a volumeRef
	if tmpVolumeInfo.VolumeRef == "" {
		//We are in destroy and the volume doesn't have a volumeRef!
		panic("ERROR - volume was found in volumeInfo map but has invalid volumeRef!")
	}

	//Send a DELETE to remove this volume from storage array
	resp, err := d.SendMsg(nil, "DELETE", "/volumes/"+tmpVolumeInfo.VolumeRef)
	defer resp.Body.Close()

	if resp.StatusCode != GenericResponseOkay && resp.StatusCode != GenericResponseNoContent {
		if resp.StatusCode == GenericResponseNotFound {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Debugf("error response Body:\n%s", string(body))

			//Next need to demarshal json data
			responseData := CallResponseError{}
			if err := json.Unmarshal(body, &responseData); err != nil {
				panic(err)
			}

			//Offer more information about what went wrong with request than just HTTP error information
			return fmt.Errorf("Error occured while trying to remove LUN mapping for volume %s! ErrorMsg=%s LocalizedMsg=%s RetCode=%s CodeType=%s Error Code (%v) and Status=%s", name, responseData.ErrorMsg, responseData.LocalizedMsg, responseData.ReturnCode, responseData.CodeType, resp.StatusCode, resp.Status)
		} else {

			//Different error than NotFound
			return fmt.Errorf("Error occured while trying to remove LUN mapping for volume %s! Error Code (%v) and Status=%s", name, resp.StatusCode, resp.Status)
		}
	}

	//Remove this volume from volumeInfo map
	delete(d.config.Volumes, name)

	return nil
}
