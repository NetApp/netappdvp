// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// Get cluster capacity stats
func (c *Client) GetClusterCapacity() (capacity *ClusterCapacity, err error) {
	var (
		clusterCapReq    GetClusterCapacityRequest
		clusterCapResult GetClusterCapacityResult
	)

	response, err := c.Request("GetClusterCapacity", clusterCapReq, NewReqID())
	if err != nil {
		log.Errorf("error detected in GetClusterCapacity API response: %+v", err)
		return nil, errors.New("device API error")
	}
	if err := json.Unmarshal([]byte(response), &clusterCapResult); err != nil {
		log.Errorf("error detected unmsarshalling json response: %+v", err)
		return nil, errors.New("json decode error")
	}
	return &clusterCapResult.Result.ClusterCapacity, err
}
