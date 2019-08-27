package invoker

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/gaeanetwork/gaea-core/tee/task"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

const (
	formatErrInvokeChaincode = "Error invoking chaincode %s[%s] to query data in channel: %s, dataID: %s, error: %s"
)

// Create a tee task
func Create(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 3); err != nil {
		return shim.Error(fmt.Sprintf("Error checking Args[ids, algorithmID, resultAddress], error: %v", err))
	}

	// dataIDs []string, algorithmID, resultAddress string
	ids, algorithmID, resultAddress := args[0], args[1], args[2]

	var dataIDs []string
	if err = json.Unmarshal([]byte(ids), &dataIDs); err != nil {
		return shim.Error(fmt.Sprintf("Failed to parse dataIDs to string slice, err: %s", err))
	} else if len(dataIDs) == 0 {
		return shim.Error(fmt.Sprintf("Error task data IDs must be non-empty"))
	}

	var timestamp *timestamp.Timestamp
	if timestamp, err = stub.GetTxTimestamp(); err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}

	teetask := &tee.Task{
		ID:                     stub.GetTxID(),
		DataIDs:                dataIDs,
		AlgorithmID:            algorithmID,
		ResultAddress:          filepath.Join(resultAddress, stub.GetTxID()),
		DataNotifications:      make(map[string]string),
		CreateSecondsTimestamp: timestamp.Seconds,
		UploadSecondsTimestamp: timestamp.Seconds,
		Partners:               make(map[string]struct{}),
	}

	if err = checkIfAlgorithmAndDataExistInTEE(stub, teetask); err != nil {
		return shim.Error(err.Error())
	}

	if err = requestAuthorizationsForAllData(stub, teetask); err != nil {
		return shim.Error(err.Error())
	}

	if err = saveTaskAndIndex(stub, teetask, true); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(stub.GetTxID()))
}

func checkIfAlgorithmAndDataExistInTEE(stub shim.ChaincodeStubInterface, teetask *tee.Task) error {
	// Check if the algorithm exists and save the algorithm hash and owner
	algorithm, err := querySharedDataByID(stub, teetask.AlgorithmID)
	if err != nil {
		return errors.Wrap(err, "failed to query algorithm in blockchain")
	}
	teetask.EvidenceHash.AlgorithmHash = algorithm.Hash
	teetask.Partners[algorithm.Owner] = struct{}{}

	// Check if all data exists and save all data hashes and owner
	teetask.EvidenceHash.DataHash = make([]string, len(teetask.DataIDs))
	for index, dataID := range teetask.DataIDs {
		data, err := querySharedDataByID(stub, dataID)
		if err != nil {
			return errors.Wrap(err, "failed to query data in blockchain")
		}

		teetask.EvidenceHash.DataHash[index] = data.Hash
		teetask.Partners[data.Owner] = struct{}{}
	}

	return nil
}

func querySharedDataByID(stub shim.ChaincodeStubInterface, id string) (*tee.SharedData, error) {
	args, channelID := [][]byte{[]byte(tee.MethodQueryDataByID), []byte(id)}, stub.GetChannelID()
	response := stub.InvokeChaincode(tee.ChaincodeName, args, channelID)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(formatErrInvokeChaincode, tee.ChaincodeName, tee.MethodQueryDataByID, channelID, id, response.Message)
	}

	var data tee.SharedData
	if err := json.Unmarshal(response.Payload, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func requestAuthorizationsForAllData(stub shim.ChaincodeStubInterface, teetask *tee.Task) error {
	privHexBytes, err := stub.GetState(task.KeyPrivHex)
	if err != nil {
		return fmt.Errorf("Error getting private key, May be wrong when instantiating or upgrading. error: %v", err)
	}
	pubBytes, err := stub.GetState(task.KeyPubHex)
	if err != nil {
		return fmt.Errorf("Error getting public key, May be wrong when instantiating or upgrading. error: %v", err)
	}
	privHex, requester := string(privHexBytes), string(pubBytes)

	for _, dataID := range teetask.DataIDs {
		notification, err := sendSharedDataRequest(stub, dataID, privHex, requester)
		if err != nil {
			return errors.Wrap(err, "failed to send shared data request")
		}

		teetask.DataNotifications[dataID] = notification.ID
	}

	// Clear private key
	privHexBytes, privHex = nil, ""
	teetask.Requester = requester
	return nil
}

func sendSharedDataRequest(stub shim.ChaincodeStubInterface, dataID, privHex, requester string) (*tee.Notification, error) {
	data, err := querySharedDataByID(stub, dataID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query data in blockchain")
	}

	args, err := constructSharedDataRequestArgs(data, privHex, requester)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct shared data request")
	}

	channelID := stub.GetChannelID()
	response := stub.InvokeChaincode(tee.ChaincodeName, args, channelID)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(formatErrInvokeChaincode, tee.ChaincodeName, tee.MethodRequest, channelID, dataID, response.Message)
	}

	var notification tee.Notification
	if err := json.Unmarshal(response.Payload, &notification); err != nil {
		return nil, err
	}

	return &notification, nil
}

func constructSharedDataRequestArgs(data *tee.SharedData, privHex, requester string) ([][]byte, error) {
	privBytes, err := hex.DecodeString(privHex)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode hex private key hex string")
	}

	args := []string{data.ID, requester}
	if len(data.Signatures) > 0 && data.Signatures[0] != "" {
		hash, sigs, err := chaincode.GetArgsHashAndSignatures(privBytes, args)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get arguments hash and signatures, args: %v", args)
		}

		return [][]byte{[]byte("request"), []byte(args[0]), []byte(args[1]), []byte(common.BytesToHex(hash)), sigs}, nil
	}

	return [][]byte{[]byte("request"), []byte(args[0]), []byte(args[1])}, nil
}

func saveTaskAndIndex(stub shim.ChaincodeStubInterface, teetask *tee.Task, saveIndex bool) error {
	taskBytes, err := json.Marshal(teetask)
	if err != nil {
		return errors.Wrap(err, "failed to parse tee task to bytes")
	}

	if err = stub.PutState(teetask.ID, taskBytes); err != nil {
		return errors.Wrap(err, "failed to put tee task to state")
	}

	if saveIndex {
		taskIDIndexKey, err := stub.CreateCompositeKey(task.TaskIDIndex, []string{task.ChaincodeName, teetask.ID})
		if err != nil {
			return errors.Wrap(err, "failed to create composite key to save tee task index")
		}

		return stub.PutState(taskIDIndexKey, chaincode.CompositeValue)
	}

	return nil
}
