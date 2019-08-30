package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/tee"
	"gitlab.com/jaderabbit/go-rabbit/tee/container"
	"gitlab.com/jaderabbit/go-rabbit/tee/task"
)

func execute(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 3); err != nil {
		return shim.Error(err.Error())
	}
	var sigs []string
	if len(args) > 4 {
		if err = json.Unmarshal([]byte(args[4]), &sigs); err != nil {
			return shim.Error("Failed to unmarshal signatures for args[4], error: " + err.Error())
		} else if len(sigs) == 0 {
			return shim.Error("Signature slice length is 0")
		}
	}

	taskID, executor, ctypeStr := args[0], args[1], args[2]
	response := chaincode.GetDataByID(stub, taskID)
	if response.Status != shim.OK {
		return shim.Error(fmt.Sprintf("Failed to query data by id: %s", response.Message))
	}

	var teetask tee.Task
	if err = json.Unmarshal(response.Payload, &teetask); err != nil {
		return shim.Error(err.Error())
	}

	if _, ok := teetask.Partners[executor]; !ok {
		return shim.Error(fmt.Sprintf("The executor does not exists in task.Partners. executor: %s, partners: %v", executor, teetask.Partners))
	}

	if err = os.MkdirAll(teetask.ResultAddress, 0755); err != nil {
		return shim.Error(fmt.Sprintf("Failed to mkdir result address, address: %s, error: %v", teetask.ResultAddress, err))
	}

	path := filepath.Join(teetask.ResultAddress, task.AlgorithmName)
	algorithm, err := ioutil.ReadFile(path)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to read algorithm file, path: %s, error: %v", path, err))
	}

	ctype, err := strconv.Atoi(ctypeStr)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to prase container type, container type: %s, error: %v", ctypeStr, err))
	}
	teetask.Container = container.Type(ctype)

	// If you choose to sign, you need to verify
	if len(sigs) > 0 && sigs[0] != "" {
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, executor); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}
	}

	if err = doExecute(stub, []byte(algorithm), &teetask); err != nil {
		return shim.Error(err.Error())
	}

	if err = saveTaskAndIndex(stub, &teetask, false); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(stub.GetTxID()))
}

func doExecute(stub shim.ChaincodeStubInterface, algorithm []byte, teetask *tee.Task) error {
	var executionLog bytes.Buffer
	executionLog.WriteString("开始下载数据...\n")
	dataList, err := downloadData(stub, teetask)
	if err != nil {
		return err
	}
	executionLog.WriteString("下载数据成功！\n")

	container := container.GetContainer(teetask.Container)
	executionLog.WriteString("正在创建可信计算环境...\n")
	if err = container.Create(); err != nil {
		return err
	}
	defer container.Destroy()
	executionLog.WriteString("创建可信计算环境成功！\n")

	executionLog.WriteString("正在装载算法和数据...\n")
	if err = container.Upload(algorithm, dataList); err != nil {
		return err
	}
	executionLog.WriteString("装载算法和数据成功！\n")

	executionLog.WriteString("正在校验程序和数据的完整性...\n")
	if err = container.Verify(teetask.EvidenceHash.AlgorithmHash, teetask.EvidenceHash.DataHash); err != nil {
		return err
	}
	executionLog.WriteString("程序和数据的完整性校验成功！\n")

	executionLog.WriteString("正在执行计算...\n")
	result, err := container.Execute()
	if err != nil {
		return fmt.Errorf("Error executing the container, err: %s, result: %s", err, string(result))
	}
	executionLog.WriteString("执行计算成功！结果是：\n")
	executionLog.Write(result)
	executionLog.WriteString("开始销毁可信计算环境及所有数据...\n")
	executionLog.WriteString("可信计算环境销毁成功！\n")

	resultHashBytes, logHashBytes := sha256.Sum256(result), sha256.Sum256(executionLog.Bytes())
	teetask.EvidenceHash.ResultHash = hex.EncodeToString(resultHashBytes[:])
	teetask.EvidenceHash.ExecutionLogHash = hex.EncodeToString(logHashBytes[:])

	return uploadResults(executionLog.Bytes())
}

// Check that the notification status is authorized, and if so, return the notification.
func checkNotificationAuthorized(stub shim.ChaincodeStubInterface, notificationID string) (*tee.Notification, error) {
	args, channelID := [][]byte{[]byte(tee.MethodQueryDataByID), []byte(notificationID)}, stub.GetChannelID()
	response := stub.InvokeChaincode(tee.ChaincodeName, args, channelID)
	if response.Status != shim.OK {
		return nil, fmt.Errorf("error invoking chaincode %s[%s] to query notification in channel: %s, error: %v",
			tee.ChaincodeName, tee.MethodQueryDataByID, channelID, response.Message+string(response.Payload))
	}

	var notification tee.Notification
	if err := json.Unmarshal(response.Payload, &notification); err != nil {
		return nil, err
	}

	if notification.Status != tee.Authorized {
		return nil, fmt.Errorf("Failed to check request data status, notificationID: %s, status: %s, message: %s", notificationID, notification.Status, notification.RefusedReason)
	}

	return &notification, nil
}
