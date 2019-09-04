package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var reviewersKey = "chaincodeMapReviewersKey"

// Status the chaincode status, review, approved and offline
type Status uint32

const (
	// Review the chaincode is in the process of being reviewed
	Review Status = iota + 1

	// Approved the chaincode has been approved.
	Approved

	// Refused the chaincode has been refused
	Refused

	// Offline The chain code is offline
	Offline
)

// RabbitChaincode used as an attribute in chaincode struct
type RabbitChaincode struct{}

// ReviewChaincode review chaincode
type ReviewChaincode struct {
	Reviewers   map[string]*reviewer `json:"reviewers,omitempty"`
	CreatorMsp  string               `json:"creator,omitempty"`
	CheckStatus Status               `json:"checkstatus,omitempty"`
	Timestamp   *timestamp.Timestamp `json:"timestamp,omitempty"`
	Mark        string               `json:"mark,omitempty"`
}

type reviewer struct {
	Timestamp   *timestamp.Timestamp `json:"timestamp,omitempty"`
	CheckStatus Status               `json:"checkstatus,omitempty"`
	Mark        string               `json:"mark,omitempty"`
}

// Init this function is executed in the chain code initialization function,
// and the modified chain code status is the review state; the fist parameter is init,
// and the second parameter is msp collection, for example: org1msp,org2msp,
// the third parameter is the mark, description of the instantiate or upgrade.
func (rcc *RabbitChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	reviewcc := &ReviewChaincode{
		Reviewers: make(map[string]*reviewer),
	}

	if len(args) > 0 && len(args[0]) > 0 {
		arrayArgs := strings.Split(args[0], ",")
		for _, arg := range arrayArgs {
			if len(arg) == 0 {
				continue
			}
			reviewcc.Reviewers[arg] = &reviewer{}
		}
	}

	if len(args) > 1 && len(args[1]) > 0 {
		reviewcc.Mark = args[1]
	}

	mspID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get msp id, err:%s", err.Error()))
	}

	// when run the uint test, set the mspID to msp1, Because the chaincode mock does not have an implementation method GetCreator
	// mspID := "msp1"

	ttamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	reviewcc.CreatorMsp = mspID
	reviewcc.Timestamp = ttamp
	reviewcc.CheckStatus = Review
	if len(reviewcc.Reviewers) == 0 {
		reviewcc.CheckStatus = Approved
	}

	byteMap, err := json.Marshal(reviewcc)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to init chaincode status when json marshal mapRev, err:%s", err.Error()))
	}

	if err = stub.PutState(reviewersKey, byteMap); err != nil {
		return shim.Error(fmt.Sprintf("failed to init chaincode status when put state reviewersKey, err:%s", err.Error()))
	}

	return shim.Success(nil)
}

// Invoke this function invoke in the Invoke function of the chaincode, check the status, review the chaincode,
// and offline the chaincode.
func (rcc *RabbitChaincode) Invoke(stub shim.ChaincodeStubInterface) (pb.Response, bool) {
	function, _ := stub.GetFunctionAndParameters()

	switch function {
	case "review":
		return rcc.review(stub, Approved), false
	case "refused":
		return rcc.review(stub, Refused), false
	case "offline":
		return rcc.offline(stub), false
	case "reonline":
		return rcc.reOnline(stub), false
	case "status":
		return rcc.getStatus(stub), false
	case "reviewhistory":
		return rcc.getReviewChaincodeHistory(stub), false
	}

	if err := rcc.checkStatus(stub); err != nil {
		return shim.Error(err.Error()), false
	}

	return shim.Success(nil), true
}

// GetStatus get the chaincode status
func (rcc *RabbitChaincode) getStatus(stub shim.ChaincodeStubInterface) pb.Response {
	mapRev, err := rcc.getReviewChaincode(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	byteMapRev, err := json.Marshal(mapRev)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(byteMapRev)
}

func (rcc *RabbitChaincode) getReviewChaincode(stub shim.ChaincodeStubInterface) (*ReviewChaincode, error) {
	byteRev, err := stub.GetState(reviewersKey)
	if err != nil {
		return nil, fmt.Errorf("failed to Review chaincode when get state, key:%s, err:%s", reviewersKey, err.Error())
	}

	mapRev := &ReviewChaincode{}

	if err = json.Unmarshal(byteRev, mapRev); err != nil {
		return nil, fmt.Errorf("failed to Review chaincode when json unmarshal mapRev, err:%s", err.Error())
	}
	return mapRev, nil
}

func (rcc *RabbitChaincode) getReviewChaincodeHistory(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetHistoryForKey(reviewersKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iterator.Close()

	arrayReviewChaincode := []*ReviewChaincode{}

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		reviewCC := &ReviewChaincode{}
		if err = json.Unmarshal(response.Value, reviewCC); err != nil {
			return shim.Error(err.Error())
		}

		arrayReviewChaincode = append(arrayReviewChaincode, reviewCC)
	}

	historyByte, err := json.Marshal(&arrayReviewChaincode)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(historyByte)
}

// CheckStatus check the chaincode status
func (rcc *RabbitChaincode) checkStatus(stub shim.ChaincodeStubInterface) error {
	mapRev, err := rcc.getReviewChaincode(stub)
	if err != nil {
		return err
	}

	if mapRev.CheckStatus == Review {
		return fmt.Errorf("the chaincode is in the process of being reviewed")
	}

	if mapRev.CheckStatus == Offline {
		return fmt.Errorf("chaincode is offline")
	}

	return nil
}

// Review this function is executed in the chain code initialization function,
// and the modified chain code status is the review state
func (rcc *RabbitChaincode) review(stub shim.ChaincodeStubInterface, chaincodeStatus Status) pb.Response {
	mapRev, err := rcc.getReviewChaincode(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if mapRev.CheckStatus == Approved {
		return shim.Error("the chaincode has been approved.")
	}

	ttamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	mspID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to Review chaincode when get msp id, err:%s", err.Error()))
	}

	// when run the uint test, set the mspID to msp1, Because the chaincode mock does not have an implementation method GetCreator
	//mspID := "msp1"

	revStatus, ok := mapRev.Reviewers[mspID]
	if !ok {
		return shim.Error("failed to Review chaincode, do not have permission to review chaincode")
	}

	if revStatus.CheckStatus == Approved && chaincodeStatus == Approved {
		return shim.Error("Repeat review chaincode")
	}

	if revStatus.CheckStatus == Refused && chaincodeStatus == Refused {
		return shim.Error("Repeat refause chaincode")
	}

	_, args := stub.GetFunctionAndParameters()

	revStatus.CheckStatus = chaincodeStatus
	revStatus.Timestamp = ttamp
	if len(args) > 0 {
		revStatus.Mark = args[0]
	}

	mapRev.Reviewers[mspID] = revStatus

	isAllReview := true
	for _, revStatus := range mapRev.Reviewers {
		if revStatus.CheckStatus != Approved {
			isAllReview = false
			break
		}
	}

	if isAllReview {
		mapRev.CheckStatus = Approved
		mapRev.Timestamp = ttamp
	}

	byteRev, err := json.Marshal(mapRev)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to Review chaincode when json Marshal mapRev, err:%s", err.Error()))
	}

	if err = stub.PutState(reviewersKey, byteRev); err != nil {
		return shim.Error(fmt.Sprintf("failed to Review chaincode when put state reviewersKey, err:%s", err.Error()))
	}

	return shim.Success(nil)
}

// Offline when chaincode status is offline, the function of the chaincode is unavailable.
// if making the chain code available again, upgrade the chaincode and then approved it or use the reOnline function
func (rcc *RabbitChaincode) offline(stub shim.ChaincodeStubInterface) pb.Response {
	ttamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	mspID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to offline chaincode when get msp id, err:%s", err.Error()))
	}

	// when run the uint test, set the mspID to msp1, Because the chaincode mock does not have an implementation method GetCreator
	// mspID := "msp1"

	mapRev, err := rcc.getReviewChaincode(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if mapRev.CheckStatus == Offline {
		return shim.Error("failed to offline chaincode, chaincode status has benn offline")
	}

	if mapRev.CreatorMsp != mspID {
		return shim.Error(fmt.Sprintf("failed to offline chaincode, chaincode , err:%s", err.Error()))
	}

	_, args := stub.GetFunctionAndParameters()
	mapRev.CheckStatus = Offline
	mapRev.Timestamp = ttamp
	if len(args) > 0 {
		mapRev.Mark = args[0]
	}

	byteRev, err := json.Marshal(mapRev)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to offline chaincode when json Marshal mapRev, err:%s", err.Error()))
	}

	if err = stub.PutState(reviewersKey, byteRev); err != nil {
		return shim.Error(fmt.Sprintf("failed to offline chaincode when put state reviewersKey, err:%s", err.Error()))
	}

	return shim.Success(nil)
}

// reOnline when chaincode status is offline, this function make the chaincode to Review or Approved status.
func (rcc *RabbitChaincode) reOnline(stub shim.ChaincodeStubInterface) pb.Response {
	ttamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	mspID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to reOnline chaincode when get msp id, err:%s", err.Error()))
	}

	// when run the uint test, set the mspID to msp1, Because the chaincode mock does not have an implementation method GetCreator
	// mspID := "msp1"

	mapRev, err := rcc.getReviewChaincode(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if mapRev.CheckStatus != Offline {
		return shim.Error("failed to reOnline chaincode, chaincode status is not offline")
	}

	if mapRev.CreatorMsp != mspID {
		return shim.Error(fmt.Sprintf("failed to reOnline chaincode, err:%s", err.Error()))
	}

	_, args := stub.GetFunctionAndParameters()
	mapRev.Timestamp = ttamp
	if len(args) > 0 {
		mapRev.Mark = args[0]
	}

	isAllReview := true
	for _, revStatus := range mapRev.Reviewers {
		if revStatus.CheckStatus != Approved {
			isAllReview = false
			break
		}
	}

	mapRev.CheckStatus = Review
	if isAllReview {
		mapRev.CheckStatus = Approved
	}

	byteRev, err := json.Marshal(mapRev)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to reOnline chaincode when json Marshal mapRev, err:%s", err.Error()))
	}

	if err = stub.PutState(reviewersKey, byteRev); err != nil {
		return shim.Error(fmt.Sprintf("failed to reOnline chaincode when put state reviewersKey, err:%s", err.Error()))
	}

	return shim.Success(nil)
}
