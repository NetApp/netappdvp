// Copyright 2016 NetApp, Inc. All Rights Reserved.

package sfapi

// APIError wrapper
type APIError struct {
	ID    int `json:"id"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"error"`
}

// QoS settings
type QoS struct {
	MinIOPS   int64 `json:"minIOPS,omitempty"`
	MaxIOPS   int64 `json:"maxIOPS,omitempty"`
	BurstIOPS int64 `json:"burstIOPS,omitempty"`
	BurstTime int64 `json:"-"`
}

// VolumePair settings
type VolumePair struct {
	ClusterPairID    int64  `json:"clusterPairID"`
	RemoteVolumeID   int64  `json:"remoteVolumeID"`
	RemoteSliceID    int64  `json:"remoteSliceID"`
	RemoteVolumeName string `json:"remoteVolumeName"`
	VolumePairUUID   string `json:"volumePairUUID"`
}

// Volume settings
type Volume struct {
	VolumeID           int64        `json:"volumeID"`
	Name               string       `json:"name"`
	AccountID          int64        `json:"accountID"`
	CreateTime         string       `json:"createTime"`
	Status             string       `json:"status"`
	Access             string       `json:"access"`
	Enable512e         bool         `json:"enable512e"`
	Iqn                string       `json:"iqn"`
	ScsiEUIDeviceID    string       `json:"scsiEUIDeviceID"`
	ScsiNAADeviceID    string       `json:"scsiNAADeviceID"`
	Qos                QoS          `json:"qos"`
	VolumeAccessGroups []int64      `json:"volumeAccessGroups"`
	VolumePairs        []VolumePair `json:"volumePairs"`
	DeleteTime         string       `json:"deleteTime"`
	PurgeTime          string       `json:"purgeTime"`
	SliceCount         int64        `json:"sliceCount"`
	TotalSize          int64        `json:"totalSize"`
	BlockSize          int64        `json:"blockSize"`
	VirtualVolumeID    string       `json:"virtualVolumeID"`
	Attributes         interface{}  `json:"attributes"`
}

type Snapshot struct {
	SnapshotID int64       `json:"snapshotID"`
	VolumeID   int64       `json:"volumeID"`
	Name       string      `json:"name"`
	Checksum   string      `json:"checksum"`
	Status     string      `json:"status"`
	TotalSize  int64       `json:"totalSize"`
	GroupID    int64       `json:"groupID"`
	CreateTime string      `json:"createTime"`
	Attributes interface{} `json:"attributes"`
}

// ListVolumesForAccountRequest tbd
type ListVolumesForAccountRequest struct {
	AccountID int64 `json:"accountID"`
}

// ListActiveVolumesRequest tbd
type ListActiveVolumesRequest struct {
	StartVolumeID int64 `json:"startVolumeID"`
	Limit         int64 `json:"limit"`
}

// ListVolumesResult tbd
type ListVolumesResult struct {
	ID     int `json:"id"`
	Result struct {
		Volumes []Volume `json:"volumes"`
	} `json:"result"`
}

// CreateVolumeRequest tbd
type CreateVolumeRequest struct {
	Name       string      `json:"name"`
	AccountID  int64       `json:"accountID"`
	TotalSize  int64       `json:"totalSize"`
	Enable512e bool        `json:"enable512e"`
	Qos        QoS         `json:"qos,omitempty"`
	Attributes interface{} `json:"attributes"`
}

// CreateVolumeResult tbd
type CreateVolumeResult struct {
	ID     int `json:"id"`
	Result struct {
		VolumeID int64 `json:"volumeID"`
	} `json:"result"`
}

// DeleteVolumeRequest tbd
type DeleteVolumeRequest struct {
	VolumeID int64 `json:"volumeID"`
}

type CloneVolumeRequest struct {
	VolumeID     int64       `json:"volumeID"`
	Name         string      `json:"name"`
	NewAccountID int64       `json:"newAccountID"`
	NewSize      int64       `json:"newSize"`
	Access       string      `json:"access"`
	SnapshotID   int64       `json:"snapshotID"`
	Attributes   interface{} `json:"attributes"`
}

type CloneVolumeResult struct {
	Id     int `json:"id"`
	Result struct {
		CloneID     int64 `json:"cloneID"`
		VolumeID    int64 `json:"volumeID"`
		AsyncHandle int64 `json:"asyncHandle"`
	} `json:"result"`
}

type CreateSnapshotRequest struct {
	VolumeID                int64       `json:"volumeID"`
	SnapshotID              int64       `json:"snapshotID"`
	Name                    string      `json:"name"`
	EnableRemoteReplication bool        `json:"enableRemoteReplication"`
	Retention               string      `json:"retention"`
	Attributes              interface{} `json:"attributes"`
}

type CreateSnapshotResult struct {
	Id     int `json:"id"`
	Result struct {
		SnapshotID int64  `json:"snapshotID"`
		Checksum   string `json:"checksum"`
	} `json:"result"`
}

type ListSnapshotsRequest struct {
	VolumeID int64 `json:"volumeID"`
}

type ListSnapshotsResult struct {
	ID     int `json:"id"`
	Result struct {
		Snapshots []Snapshot `json:"snapshots"`
	} `json:"result"`
}

type RollbackToSnapshotRequest struct {
	VolumeID         int64       `json:"volumeID"`
	SnapshotID       int64       `json:"snapshotID"`
	SaveCurrentState bool        `json:"saveCurrentState"`
	Name             string      `json:"name"`
	Attributes       interface{} `json:"attributes"`
}

type RollbackToSnapshotResult struct {
	ID     int `json:"id"`
	Result struct {
		Checksum   string `json:"checksum"`
		SnapshotID int64  `json:"snapshotID"`
	} `json:"result"`
}

type DeleteSnapshotRequest struct {
	SnapshotID int64 `json:"snapshotID"`
}

// AddVolumesToVolumeAccessGroupRequest tbd
type AddVolumesToVolumeAccessGroupRequest struct {
	VolumeAccessGroupID int64   `json:"volumeAccessGroupID"`
	Volumes             []int64 `json:"volumes"`
}

// CreateVolumeAccessGroupRequest tbd
type CreateVolumeAccessGroupRequest struct {
	Name       string   `json:"name"`
	Volumes    []int64  `json:"volumes,omitempty"`
	Initiators []string `json:"initiators,omitempty"`
}

// CreateVolumeAccessGroupResult tbd
type CreateVolumeAccessGroupResult struct {
	ID     int `json:"id"`
	Result struct {
		VagID int64 `json:"volumeAccessGroupID"`
	} `json:"result"`
}

// AddInitiatorsToVolumeAccessGroupRequest tbd
type AddInitiatorsToVolumeAccessGroupRequest struct {
	Initiators []string `json:"initiators"`
	VAGID      int64    `json:"volumeAccessGroupID"`
}

// ListVolumeAccessGroupsRequest tbd
type ListVolumeAccessGroupsRequest struct {
	StartVAGID int64 `json:"startVolumeAccessGroupID,omitempty"`
	Limit      int64 `json:"limit,omitempty"`
}

// ListVolumesAccessGroupsResult tbd
type ListVolumesAccessGroupsResult struct {
	ID     int `json:"id"`
	Result struct {
		Vags []VolumeAccessGroup `json:"volumeAccessGroups"`
	} `json:"result"`
}

// EmptyResponse tbd
type EmptyResponse struct {
	ID     int `json:"id"`
	Result struct {
	} `json:"result"`
}

// VolumeAccessGroup tbd
type VolumeAccessGroup struct {
	Initiators     []string    `json:"initiators"`
	Attributes     interface{} `json:"attributes"`
	DeletedVolumes []int64     `json:"deletedVolumes"`
	Name           string      `json:"name"`
	VAGID          int64       `json:"volumeAccessGroupID"`
	Volumes        []int64     `json:"volumes"`
}

// GetAccountByNameRequest tbd
type GetAccountByNameRequest struct {
	Name string `json:"username"`
}

// GetAccountByIDRequest tbd
type GetAccountByIDRequest struct {
	AccountID int64 `json:"accountID"`
}

// GetAccountResult tbd
type GetAccountResult struct {
	ID     int `json:"id"`
	Result struct {
		Account Account `json:"account"`
	} `json:"result"`
}

// Account tbd
type Account struct {
	AccountID       int64       `json:"accountID,omitempty"`
	Username        string      `json:"username,omitempty"`
	Status          string      `json:"status,omitempty"`
	Volumes         []int64     `json:"volumes,omitempty"`
	InitiatorSecret string      `json:"initiatorSecret,omitempty"`
	TargetSecret    string      `json:"targetSecret,omitempty"`
	Attributes      interface{} `json:"attributes,omitempty"`
}

// AddAccountRequest tbd
type AddAccountRequest struct {
	Username        string      `json:"username"`
	InitiatorSecret string      `json:"initiatorSecret,omitempty"`
	TargetSecret    string      `json:"targetSecret,omitempty"`
	Attributes      interface{} `json:"attributes,omitempty"`
}

// AddAccountResult tbd
type AddAccountResult struct {
	ID     int `json:"id"`
	Result struct {
		AccountID int64 `json:"accountID"`
	} `json:"result"`
}
