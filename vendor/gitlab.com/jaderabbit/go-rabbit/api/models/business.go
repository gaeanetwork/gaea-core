package models

import (
	"mime/multipart"

	"gitlab.com/jaderabbit/go-rabbit/i18n"
)

var (
	busiInfos = []*BusiInfo{&BusiInfo{
		Name:        SystemEChainName,
		Status:      "NOTOPEN",
		Description: "The basic business",
		Version:     "1.0",
	}}
)

// BusiInfo for view
type BusiInfo struct {
	Name        string                  `form:"name" json:"name"`
	Description string                  `form:"description" json:"description"`
	Version     string                  `form:"version" json:"version"`
	ChannelID   string                  `form:"channelID" json:"channelID"`
	Status      string                  `json:"status"`
	Chaincode   []*multipart.FileHeader `json:"chaincode"`
}

// GetBusiness finding all business listings
func GetBusiness(name string) ([]*BusiInfo, error) {
	if name == "" {
		return busiInfos, nil
	}

	for i, busiInfo := range busiInfos {
		if busiInfo.Name == name {
			return busiInfos[i : i+1], nil
		}
	}

	return nil, i18n.BusinessNotFoundErr(name)
}

// AddBusiness add a business, except channelID and status is empty
func AddBusiness(busiInfo *BusiInfo) error {
	if busiInfo == nil || busiInfo.Name == "" {
		return i18n.BusinessEmptyErr("nil")
	}

	busiInfo.Status = "NOTOPEN"
	busiInfos = append(busiInfos, busiInfo)
	return nil
}

// StartBusiness start a business, only channelID will be changed
func StartBusiness(name string, channelID string) error {
	if name == "" {
		return i18n.BusinessEmptyErr("nil")
	}

	for _, binfo := range busiInfos {
		if name == binfo.Name {
			binfo.ChannelID = channelID
			return nil
		}
	}

	return i18n.BusinessNotFoundErr(name)
}

// UpgradeBusiness upgrade a business, only chaincode and version will be changed
func UpgradeBusiness(busiInfo *BusiInfo) error {
	if busiInfo == nil || busiInfo.Name == "" {
		return i18n.BusinessEmptyErr("nil")
	}

	for _, binfo := range busiInfos {
		if busiInfo.Name == binfo.Name {
			binfo.Chaincode = busiInfo.Chaincode
			if busiInfo.Version <= binfo.Version {
				return i18n.BusinessVersionTooLowErr{SrcV: binfo.Version, DestV: busiInfo.Version}
			}
			binfo.Version = busiInfo.Version
			return nil
		}
	}

	return i18n.BusinessNotFoundErr(busiInfo.Name)
}
