package v20240610preview

// Copyright (c) Microsoft Corporation
// Licensed under the Apache License 2.0.

// NodePoolAutoscaling represents a node pool autoscaling configuration.
// Visibility for the entire struct is "read".
type NodePoolAutoscaling struct {
	MinReplicas int32 `json:"minReplicas,omitempty"`
	MaxReplicas int32 `json:"maxReplicas,omitempty"`
}

// NodePoolProfile represents a worker node pool configuration.
// Visibility for the entire struct is "read".
type NodePoolProfile struct {
	Name                   string              `json:"name,omitempty"`
	Version                string              `json:"version,omitempty"`
	Labels                 []string            `json:"labels,omitempty"`
	Taints                 []string            `json:"taints,omitempty"`
	DiskSize               int32               `json:"diskSize,omitempty"`
	EphemeralOSDisk        bool                `json:"ephemeralOsDisk,omitempty"`
	Replicas               int32               `json:"replicas,omitempty"`
	SubnetID               string              `json:"subnetId,omitempty"`
	EncryptionAtHost       bool                `json:"encryptionAtHost,omitempty"`
	AutoRepair             bool                `json:"autoRepair,omitempty"`
	DiscEncryptionSetID    string              `json:"discEncryptionSetId,omitempty"`
	TuningConfigs          []string            `json:"tuningConfigs,omitempty"`
	AvailabilityZone       string              `json:"availabilityZone,omitempty"`
	DiscStorageAccountType string              `json:"discStorageAccountType,omitempty"`
	VMSize                 string              `json:"vmSize,omitempty"`
	Autoscaling            NodePoolAutoscaling `json:"autoscaling,omitempty"`
}
