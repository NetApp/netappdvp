// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/netapp/netappdvp/utils"
)

// ListVolumesForAccount tbd
func (c *Client) ListVolumesForAccount(listReq *ListVolumesForAccountRequest) (volumes []Volume, err error) {
	response, err := c.Request("ListVolumesForAccount", listReq, NewReqID())
	if err != nil {
		log.Errorf("error detected in ListVolumesForAccount API response: %+v", err)
		return nil, errors.New("device API error")
	}
	var result ListVolumesResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Errorf("error detected unmarshalling ListVolumesForAccount API response: %+v", err)
		return nil, errors.New("json-decode error")
	}
	volumes = result.Result.Volumes
	return volumes, err
}

// GetVolumeByID tbd
func (c *Client) GetVolumeByID(volID int64) (v Volume, err error) {
	var req ListActiveVolumesRequest
	req.StartVolumeID = volID
	req.Limit = 1
	volumes, err := c.ListActiveVolumes(&req)
	if err != nil {
		return v, err
	}
	if len(volumes) < 1 {
		return Volume{}, fmt.Errorf("Failed to find volume with ID: %d", volID)
	}
	return volumes[0], nil
}

// GetVolumeByDockerName tbd
func (c *Client) GetVolumeByDockerName(n string, acctID int64) (v Volume, err error) {
	vols, err := c.GetVolumesByDockerName(n, acctID)
	if err == nil && len(vols) == 1 {
		return vols[0], nil
	}

	if len(vols) > 1 {
		err = fmt.Errorf("Found more than one Volume with DockerName: %s for Account: %d", n, acctID)
	} else if len(vols) < 1 {
		err = fmt.Errorf("Failed to find any Volumes with DockerName: %s for Account: %d", n, acctID)
	}
	return v, err
}

// GetVolumesByDockerName tbd
func (c *Client) GetVolumesByDockerName(dockerName string, acctID int64) (v []Volume, err error) {
	log.Debug("attempting GetVolumesByDockerName call")
	var listReq ListVolumesForAccountRequest
	var foundVolumes []Volume
	listReq.AccountID = acctID
	volumes, err := c.ListVolumesForAccount(&listReq)
	if err != nil {
		log.Errorf("error detected in ListVolumesForAccount API response: %+v", err)
		return foundVolumes, errors.New("device API error")
	}
	for _, vol := range volumes {
		attrs, _ := vol.Attributes.(map[string]interface{})
		log.Debugf("looking for docker-name %+v in %+v\n", dockerName, attrs)
		if attrs["docker-name"] == dockerName && vol.Status == "active" {
			foundVolumes = append(foundVolumes, vol)
		} else if vol.Name == strings.Replace(dockerName, "_", "-", -1) {
			foundVolumes = append(foundVolumes, vol)
		}
	}
	if len(foundVolumes) > 1 {
		log.Warningf("found more than one volume with the docker-name: %s\n%+v", dockerName, foundVolumes)
	}
	if len(foundVolumes) == 0 {
		log.Warningf("failed to find volume with docker name: %s", dockerName)
		return foundVolumes, errors.New("volume not found")
	}
	return foundVolumes, nil
}

// GetVolumeByName tbd
func (c *Client) GetVolumeByName(n string, acctID int64) (v Volume, err error) {
	vols, err := c.GetVolumesByName(n, acctID)
	if err == nil && len(vols) == 1 {
		return vols[0], nil
	}

	if len(vols) > 1 {
		err = fmt.Errorf("Found more than one Volume with Name: %s for Account: %d", n, acctID)
	} else if len(vols) < 1 {
		err = fmt.Errorf("Failed to find any Volumes with Name: %s for Account: %d", n, acctID)
	}
	return v, err
}

// GetVolumesByName tbd
func (c *Client) GetVolumesByName(sfName string, acctID int64) (v []Volume, err error) {
	var listReq ListVolumesForAccountRequest
	var foundVolumes []Volume
	listReq.AccountID = acctID
	volumes, err := c.ListVolumesForAccount(&listReq)
	if err != nil {
		log.Errorf("error retrieving volumes in GetVolumesByname: %+v ", err)
		return foundVolumes, errors.New("device API error")
	}
	for _, vol := range volumes {
		if vol.Name == sfName && vol.Status == "active" {
			foundVolumes = append(foundVolumes, vol)
		}
	}
	if len(foundVolumes) > 1 {
		log.Warningf("discovered more than one volume with the name: %s\n%+v", sfName, foundVolumes)
	}
	if len(foundVolumes) == 0 {
		log.Errorf("no volumes found in list matching name %s and account %d", sfName, acctID)
		return foundVolumes, errors.New("volume not found")
	}
	return foundVolumes, nil
}

// ListActiveVolumes tbd
func (c *Client) ListActiveVolumes(listVolReq *ListActiveVolumesRequest) (volumes []Volume, err error) {
	response, err := c.Request("ListActiveVolumes", listVolReq, NewReqID())
	if err != nil {
		log.Errorf("error response from ListActiveVolumes request: %+v ", err)
		return nil, errors.New("device API error")
	}
	var result ListVolumesResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Errorf("error detected unmarshalling ListActiveVolumes API response: %+v", err)
		return nil, errors.New("json-decode error")
	}
	volumes = result.Result.Volumes
	return volumes, err
}

func (c *Client) CloneVolume(req *CloneVolumeRequest) (vol Volume, err error) {
	response, err := c.Request("CloneVolume", req, NewReqID())
	var result CloneVolumeResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Errorf("error detected unmarshalling CloneVolume API response: %+v", err)
		return Volume{}, errors.New("json-decode error")
	}

	wait := 0
	multiplier := 1
	for wait < 10 {
		wait += wait
		vol, err = c.GetVolumeByID(result.Result.VolumeID)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(multiplier))
		multiplier *= wait
	}
	return
}

// CreateVolume tbd
func (c *Client) CreateVolume(createReq *CreateVolumeRequest) (vol Volume, err error) {
	response, err := c.Request("CreateVolume", createReq, NewReqID())
	if err != nil {
		log.Errorf("error response from CreateVolume request: %+v ", err)
		return Volume{}, errors.New("device API error")
	}
	var result CreateVolumeResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Errorf("error detected unmarshalling CreateVolume API response: %+v", err)
		return Volume{}, errors.New("json-decode error")
	}

	vol, err = c.GetVolumeByID(result.Result.VolumeID)
	return
}

// AddVolumeToAccessGroup tbd
func (c *Client) AddVolumeToAccessGroup(groupID int64, volIDs []int64) (err error) {
	req := &AddVolumesToVolumeAccessGroupRequest{
		VolumeAccessGroupID: groupID,
		Volumes:             volIDs,
	}
	_, err = c.Request("AddVolumesToVolumeAccessGroup", req, NewReqID())
	if err != nil {
		log.Errorf("error response from Add to VAG request: %+v ", err)
		return errors.New("device API error")
	}
	return err
}

// DeleteRange tbd
func (c *Client) DeleteRange(startID, endID int64) {
	idx := startID
	for idx < endID {
		c.DeleteVolume(idx)
	}
	return
}

// DeleteVolume tbd
func (c *Client) DeleteVolume(volumeID int64) (err error) {
	// TODO(jdg): Add options like purge=True|False, range, ALL etc
	var req DeleteVolumeRequest
	req.VolumeID = volumeID
	_, err = c.Request("DeleteVolume", req, NewReqID())
	if err != nil {
		// TODO: distinguish what the error was?
		log.Errorf("error response from DeleteVolume request: %+v ", err)
		return errors.New("device API error")
	}
	return
}

// DetachVolume tbd
func (c *Client) DetachVolume(v Volume) (err error) {
	if c.SVIP == "" {
		log.Errorf("error response from DetachVolume request: %+v ", err)
		return errors.New("detach volume error")
	}
	tgt := &utils.IscsiTargetInfo{
		IP:     c.SVIP,
		Portal: c.SVIP,
		Iqn:    v.Iqn,
	}
	err = utils.IscsiDisableDelete(tgt)
	return
}

// AttachVolume tbd
func (c *Client) AttachVolume(v *Volume, iface string) (path, device string, err error) {
	var req GetAccountByIDRequest
	path = "/dev/disk/by-path/ip-" + c.SVIP + "-iscsi-" + v.Iqn + "-lun-0"

	if c.SVIP == "" {
		err = errors.New("unable to perform iSCSI actions without setting SVIP")
		log.Errorf("unable to attach volume, SVIP is NOT set")
		return path, device, err
	}

	if utils.IscsiSupported() == false {
		err := errors.New("unable to attach, open-iscsi tools not found on host")
		log.Errorf("unable to attach volume, open-iscsi utils not found")
		return path, device, err
	}

	req.AccountID = v.AccountID
	a, err := c.GetAccountByID(&req)
	if err != nil {
		log.Errorf("failed to get account %v, error: %+v ", v.AccountID, err)
		return path, device, errors.New("volume attach failure")
	}

	// Make sure it's not already attached
	if utils.WaitForPathToExist(path, 1) {
		log.Debugf("get device file from path: %s", path)
		device = strings.TrimSpace(utils.GetDeviceFileFromIscsiPath(path))
		return path, device, nil
	}

	err = utils.LoginWithChap(v.Iqn, c.SVIP, a.Username, a.InitiatorSecret, iface)
	if err != nil {
		log.Errorf("failed to login with CHAP credentials: %+v ", err)
		return path, device, err
	}
	if utils.WaitForPathToExist(path, 5) {
		device = strings.TrimSpace(utils.GetDeviceFileFromIscsiPath(path))
		return path, device, nil
	}
	return path, device, nil
}
