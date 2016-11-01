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
		log.Error(err)
		return nil, err
	}
	if err := json.Unmarshal([]byte(response), &clusterCapResult); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &clusterCapResult.Result.ClusterCapacity, err
}
