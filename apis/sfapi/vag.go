// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// CreateVolumeAccessGroup tbd
func (c *Client) CreateVolumeAccessGroup(r *CreateVolumeAccessGroupRequest) (vagID int64, err error) {
	var result CreateVolumeAccessGroupResult
	response, err := c.Request("CreateVolumeAccessGroup", r, NewReqID())
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Fatal(err)
		return 0, err
	}
	vagID = result.Result.VagID
	return

}

// ListVolumeAccessGroups tbd
func (c *Client) ListVolumeAccessGroups(r *ListVolumeAccessGroupsRequest) (vags []VolumeAccessGroup, err error) {
	response, err := c.Request("ListVolumeAccessGroups", r, NewReqID())
	if err != nil {
		log.Error(err)
		return
	}
	var result ListVolumesAccessGroupsResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		log.Fatal(err)
		return nil, err
	}
	vags = result.Result.Vags
	return
}

// AddInitiatorsToVolumeAccessGroup tbd
func (c *Client) AddInitiatorsToVolumeAccessGroup(r *AddInitiatorsToVolumeAccessGroupRequest) error {
	response, err := c.Request("AddInitiatorsToVolumeAccessGroup", r, NewReqID())
	if err != nil {
		log.Error(string(response))
		log.Error(err)
		return err
	}
	return nil
}
