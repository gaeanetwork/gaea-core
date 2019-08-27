package models

import (
	"gitlab.com/jaderabbit/go-rabbit/chaincode/asset"
)

// Asset input parameter
type Asset struct {
	Key         string                  `json:"key"`
	MapData     map[string]*asset.Field `json:"map_data,omitempty"`
	IsPublic    bool                    `json:"is_public,omitempty"`
	RemoveFiled []string                `json:"remove_field,omitempty"`
}
