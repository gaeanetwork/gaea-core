package system

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/common/flogging"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
	"gitlab.com/jaderabbit/go-rabbit/core/types"
)

var (
	ManagerRole = &Role{
		ID: "64000714-9bb7-47c7-a725-f34316fe3776",
	}

	logger = flogging.MustGetLogger("chaincode.model")
)

// Status the user/role/org status, including normal and deactivated
type Status int32

const (
	// Inactivated register user status, need to be checked the account by root
	Inactivated Status = iota + 1

	// OnLine the status is online
	OnLine

	// Delete the user/role/org was deleted
	Delete
)

var userStatus = [...]string{
	"inactivated",
	"online",
	"delete",
}

func (s Status) String() string {
	if Inactivated <= s && s <= Delete {
		return userStatus[s-1]
	}
	return ""
}

// Check check whether value is valid
func (s Status) Check() bool {
	if s < Inactivated {
		return false
	}

	if s > Delete {
		return false
	}

	return true
}

// User the user in the EBaaS
type User struct {
	ID             string              `json:"id,omitempty"`
	UserName       string              `json:"user_name,omitempty"`
	Password       string              `json:"password,omitempty"`
	Status         Status              `json:"status,omitempty"`
	IsRoot         bool                `json:"is_root,omitempty"`
	IsAdmin        bool                `json:"is_admin,omitempty"`
	RealName       string              `json:"real_name,omitempty"`
	PhoneNumber    string              `json:"phone_number,omitempty"`
	RoleID         string              `json:"role,omitempty"`
	OrgID          string              `json:"org,omitempty"`
	DepartmentID   string              `json:"department,omitempty"`
	CreateTime     int64               `json:"create_time,omitempty"`
	LastUpdateTime int64               `json:"last_update_time,omitempty"`
	TxID           string              `json:"txID,omitempty"`
	HostIP         string              `json:"hostip,omitempty"`
	MapMSPUser     map[string]*MSPUser `json:"mapmspuser,omitempty"`
}

type MSPUser struct {
	MSPID     string `json:"mspid,omitempty"`
	IsDefault bool   `json:"default"`
}

// AddMSPUser add msp user and check the mspUserName and MSPID
func (u *User) AddMSPUser(mspUserName, MSPID string, isDefault bool) error {
	if len(mspUserName) == 0 {
		return errors.New("not specified mspUserName")
	}

	if len(MSPID) == 0 {
		return errors.New("not specified MSPID")
	}

	if u.MapMSPUser == nil {
		u.MapMSPUser = make(map[string]*MSPUser)
	}

	u.MapMSPUser[mspUserName] = &MSPUser{
		MSPID:     MSPID,
		IsDefault: isDefault,
	}
	return nil
}

// GetDefaultMSP get the default msp info of the user,
// the return args, first is msp user name, second is msp id
func (u *User) GetDefaultMSP() (string, string, error) {
	if u.MapMSPUser == nil {
		return "", "", fmt.Errorf("user(Name:%s) has not msp user", u.UserName)
	}

	for mspUserName, mspInfo := range u.MapMSPUser {
		if !mspInfo.IsDefault {
			continue
		}
		return mspUserName, mspInfo.MSPID, nil
	}
	return "", "", fmt.Errorf("user(Name:%s) has not msp user", u.UserName)
}

// Role the role of user in the EBaaS
type Role struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Status         Status `json:"status,omitempty"`
	OrgID          string `json:"org,omitempty"`
	CreateTime     int64  `json:"create_time,omitempty"`
	LastUpdateTime int64  `json:"last_update_time,omitempty"`
	TxID           string `json:"txID,string"`
}

// Organization the organization of user in the EBaaS
type Organization struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Flag           string `json:"flag,omitempty"`
	Status         Status `json:"status,omitempty"`
	CreateTime     int64  `json:"create_time,omitempty"`
	LastUpdateTime int64  `json:"last_update_time,omitempty"`
	TxID           string `json:"txID,omitempty"`
}

// Department the department of organization in the EBaas
type Department struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Manager        string `json:"manager,omitempty"`
	OrgID          string `json:"org,omitempty"`
	ParentID       string `json:"parent,omitempty"`
	Status         Status `json:"status,omitempty"`
	CreateTime     int64  `json:"create_time,omitempty"`
	LastUpdateTime int64  `json:"last_update_time,omitempty"`
	TxID           string `json:"txID,omitempty"`
}

// Menu the menu of the user, role, department
type Menu struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Parent         string  `json:"parent,omitempty"`
	Link           string  `json:"link,omitempty"`
	Img            string  `json:"img,omitempty"`
	Status         Status  `json:"status,omitempty"`
	OrgID          string  `json:"org,omitempty"`
	CreateTime     int64   `json:"create_time,omitempty"`
	LastUpdateTime int64   `json:"last_update_time,omitempty"`
	TxID           string  `json:"txID,omitempty"`
	Index          float64 `json:"index,omitempty"`
}

type MenuList []*Menu

func (ml MenuList) Len() int {
	return len(ml)
}

func (ml MenuList) Swap(i, j int) {
	ml[i], ml[j] = ml[j], ml[i]
}

func (ml MenuList) Less(i, j int) bool {
	return ml[i].Index < ml[j].Index
}

// RelationType the relation type
type RelationType uint32

const (
	// UserAndMenu the relation of user and menu
	UserAndMenu RelationType = iota + 1
	// RoleAndMenu the relation of role and menu
	RoleAndMenu
	// DepartmentAndMenu the relation of department and menu
	DepartmentAndMenu
)

// Check check whether value is valid
func (r RelationType) Check() bool {
	if r < UserAndMenu {
		return false
	}

	if r > DepartmentAndMenu {
		return false
	}

	return true
}

// Relation the relation of menu and user, role, department
type Relation struct {
	LeftID         string       `json:"left,omitempty"`
	RightID        []string     `json:"right,omitempty"`
	Type           RelationType `json:"type,omitempty"`
	TxID           string       `json:"txID,omitempty"`
	OrgID          string       `json:"org,omitempty"`
	CreateTime     int64        `json:"create_time,omitempty"`
	LastUpdateTime int64        `json:"last_update_time,omitempty"`
}

// QueryResult the Query result, total number, data, page size,page index
type QueryResult struct {
	TotalNumber int         `json:"totalNumber,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	PageSize    int         `json:"pagesize,omitempty"`
	PageIndex   int         `json:"pageindex,omitempty"`
}

// IPType the IP type
type IPType uint32

const (
	// BlackList the ip black list
	BlackList IPType = iota + 1
	// WhiteList the ip white list
	WhiteList
)

// Check check whether value is valid
func (ip IPType) Check() bool {
	if ip < BlackList {
		return false
	}

	if ip > WhiteList {
		return false
	}

	return true
}

// IPList ip between start ip and end ip
type IPList struct {
	ID             string `json:"id,omitempty"`
	Start          int64  `json:"start,omitempty"`
	End            int64  `json:"end,omitempty"`
	StartIP        string `json:"startip,omitempty"`
	EndIP          string `json:"endip,omitempty"`
	Mark           string `json:"mark,omitempty"`
	Type           IPType `json:"type,omitempty"`
	TxID           string `json:"txID,omitempty"`
	OrgID          string `json:"org,omitempty"`
	CreateTime     int64  `json:"create_time,omitempty"`
	LastUpdateTime int64  `json:"last_update_time,omitempty"`
}

// BusinessStatus the business status
type BusinessStatus string

const (
	// CreateBusiness business status- create
	CreateBusiness BusinessStatus = "create"

	// OnLineBusiness online status
	OnLineBusiness BusinessStatus = "online"

	// OffLineBusiness offline status
	OffLineBusiness BusinessStatus = "offline"
)

// Business for chaincode manage
type Business struct {
	ID             string                          `json:"id,omitempty"`
	Name           string                          `json:"name,omitempty"`
	Description    string                          `json:"description,omitempty"`
	Version        string                          `json:"version,omitempty"`
	ChannelID      string                          `json:"channelID,omitempty"`
	Status         BusinessStatus                  `json:"status,omitempty"`
	Chaincode      map[string]*types.ChaincodeInfo `json:"chaincode,omitempty"`
	OrgID          string                          `json:"org,omitempty"`
	CreateTime     int64                           `json:"create_time,omitempty"`
	LastUpdateTime int64                           `json:"last_update_time,omitempty"`
	TxID           string                          `json:"txID,omitempty"`
}

func ConvertBusinessStatus(str string) (BusinessStatus, error) {
	switch str {
	case "create":
		return CreateBusiness, nil
	case "online":
		return OnLineBusiness, nil
	case "offline":
		return OffLineBusiness, nil
	default:
		return "", fmt.Errorf("%s is not business status", str)
	}
}

type Union struct {
	ID               string                     `json:"id,omitempty"`
	Name             string                     `json:"name,omitempty"`
	Description      string                     `json:"description,omitempty"`
	ContactName      string                     `json:"contactname,omitempty"`
	ContactNumber    string                     `json:"contactnumber,omitempty"`
	OrdererType      string                     `json:"orderertype,omitempty"`
	Participants     map[string]*ParticipantOrg `json:"participants,omitempty"`
	Domain           string                     `json:"domain,omitempty"`
	OrdererHostNames []string                   `json:"orderer_host_names,omitempty"`
	Status           UnionStatus                `json:"status,omitempty"`
	OrgID            string                     `json:"orgid,omitempty"`
	CreateMSPID      string                     `json:"createmspid,omitempty"`
	CreateTime       int64                      `json:"create_time,omitempty"`
	LastUpdateTime   int64                      `json:"last_update_time,omitempty"`
	TxID             string                     `json:"txID,omitempty"`
}

type ParticipantOrg struct {
	Status         ParticipantStatus `json:"status,omitempty"`
	Peers          map[string]*Peer  `json:"peers,omitempty"`
	UserCount      int               `json:"user_count,omitempty"`
	MSPID          string            `json:"mspid,omitempty"`
	CreateTime     int64             `json:"create_time,omitempty"`
	LastUpdateTime int64             `json:"last_update_time,omitempty"`
	TxID           string            `json:"txID,omitempty"`
}

type Peer struct {
	NetworkID   string `json:"networkid,omitempty"`
	Address     string `json:"address,omitempty"`
	LocalMspID  string `json:"localmspid,omitempty"`
	IsOrdererOn bool   `json:"isordereron,omitempty"`
}

type ParticipantStatus int

const (
	// PartWaitingForJoin waiting for participant join
	PartWaitingForJoin ParticipantStatus = 1 + iota

	// PartJoin participant join
	PartJoin

	// PartRefuse participant refused to join unin
	PartRefuse
)

var participantStatus = [...]string{
	"waiting",
	"join",
	"refause",
}

func (ps ParticipantStatus) String() (string, error) {
	if PartWaitingForJoin <= ps && ps <= PartRefuse {
		return participantStatus[ps-1], nil
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(ps))
	return "", fmt.Errorf("non-ParticipantStatus(%s)", string(buf[n:]))
}

type UnionStatus int

const (
	// UnionInDeployment the union is deploying
	UnionInDeployment UnionStatus = 1 + iota

	// UnionRunning the union is running
	UnionRunning

	// UnionStop the union stop running
	UnionStop
)

var unionStatus = [...]string{
	"deploy",
	"running",
	"stop",
}

func (us UnionStatus) String() (string, error) {
	if UnionInDeployment <= us && us <= UnionStop {
		return unionStatus[us-1], nil
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(us))
	return "", fmt.Errorf("non-ParticipantStatus(%s)", string(buf[n:]))
}

func InitModel() {
	role := &Role{}
	if err := config.GetRabbitViper().UnmarshalKey("defaultRole", role); err != nil {
		logger.Panicf("Could not Unmarshal %s YAML config, err: %v", "defaultRole", err)
	}
	ManagerRole.Name = role.Name
	ManagerRole.Status = OnLine
	ManagerRole.OrgID = "-"
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
