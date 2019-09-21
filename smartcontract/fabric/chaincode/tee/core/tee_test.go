package main

/*
    ============================ Deploy tee chaincode ============================
	./peer chaincode package teepack.out -n tee_data -v 1.0 -s -S -p github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/core
	mkdir $HOME/chaincodes/tee
	mv -fv teepack.out $HOME/chaincodes/tee/teepack.out

	# install
	CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode install $HOME/chaincodes/tee/teepack.out

	# instantiate
	CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode instantiate -C syschannel -n tee_data -v 1.0 -c '{"Args":[]}' -o orderer.rabbit.com:7050

	# upgrade
	CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode upgrade -C syschannel -n tee_data -v 1.1 -c '{"Args":[]}' -o orderer.rabbit.com:7050

	============================ Test tee chaincode ============================
	./peer chaincode invoke -C syschannel -n tee_data -c '{"Args":["upload","Ciphertext","Hash","Description","Owner"]}' -o orderer.rabbit.com:7050
	./peer chaincode invoke -C syschannel -n tee_data -c '{"Args":["upload", "A.secert", "0x7db880f9b7ebf56ae271497ca98edc4743efec62f7228e77644940cc95d95263", "A_个人简历", "0x04e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01", "0x8525260dd74e2e5f1024523dc14d830f4389d2186bbeb9d95d2eed93909322e2", "[\"0x7b2252223a35303637393136363138393338363130373932353639373137363431333532323635343538313534393334363832333732393630333033363837363533313136363530313534303039323430392c2253223a36323731303232323938333638303730353438313638323138363437333337383339303033303231313736353535363635343636323834363030313939393630363233353935343230303539337d\"]"]}' -o orderer.rabbit.com:7050
	./peer chaincode invoke -C syschannel -n tee_data -c '{"Args":["queryDataByID","10fe8cc74ca2cf2d4dacf9233885bdc40c643a4c64aa1411c10ec4cb93c312a5"]}' -o orderer.rabbit.com:7050
	./peer chaincode invoke -C syschannel -n tee_data -c '{"Args":["update","10fe8cc74ca2cf2d4dacf9233885bdc40c643a4c64aa1411c10ec4cb93c312a5","Ciphertext2","Summary2","Description2"]}' -o orderer.rabbit.com:7050
	./peer chaincode query -C syschannel -n tee_data -c '{"Args":["queryDataByID","10fe8cc74ca2cf2d4dacf9233885bdc40c643a4c64aa1411c10ec4cb93c312a5"]}' -o orderer.rabbit.com:7050
	./peer chaincode query -C syschannel -n tee_data -c '{"Args":["queryHistoryByDID","10fe8cc74ca2cf2d4dacf9233885bdc40c643a4c64aa1411c10ec4cb93c312a5"]}' -o orderer.rabbit.com:7050

	./peer chaincode invoke -C syschannel -n tee_data -c '{"Args":["queryDataByOwner","Owner"]}' -o orderer.rabbit.com:7050


	============================ Error Notes ================================
	Error: could not assemble transaction, err proposal response was not successful, error code 500, msg chaincode registration failed: container exited with 2
		When the chain code is upgraded, an error is reported. When restarting the peer, it is found that it cannot be started.
		It is a config.go file referenced in my chain code. This file has an init function that loads a yaml configuration file
		in the same directory. When the chain code of the fabric is packaged and installed into the docker, the non-go file will
		not be loaded. Therefore, when the file is found during instantiation, the panic is reported incorrectly. Since the
		fabric fails to start the docker container, the docker container is automatically deleted. Therefore, the error that
		should have	been printed is killed by the fabric, so the error is caused. Keep in mind: all chain codes must be written
		in pure go files, without any non-go dependencies, such as c files, yml files, etc., otherwise the error will be
		reported incorrectly, and the reason can not be found! ! !
*/

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gaeanetwork/gaea-core/crypto/ecc"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	args := [][]byte{}
	response := stub.MockInit("1", args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)

	args1 := [][]byte{[]byte("asdfasdfasdf")}
	response1 := stub.MockInit("2", args1)
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
}

func Test_Invoke_upload(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpload), []byte(""), []byte("Hash"), []byte("Description"), []byte("Owner")})
	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpload), []byte("Ciphertext"), []byte(""), []byte("Description"), []byte("Owner")})
	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpload), []byte("Ciphertext"), []byte("Hash"), []byte(""), []byte("Owner")})
	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpload), []byte("Ciphertext"), []byte("Hash"), []byte("Description"), []byte("")})

	response := uploadData(t, stub, "1")

	data := payloadIsOk(t, response.Payload)

	bs, err1 := stub.GetState(data.ID)
	assert.NoError(t, err1)
	assert.Equal(t, bs, response.Payload)

	indexWorkWell(t, stub, ownerIDIndex, []string{data.Owner, data.ID})
}

func uploadData(t *testing.T, stub *shim.MockStub, txID string) pb.Response {
	response := stub.MockInvoke(txID, [][]byte{[]byte(tee.MethodUpload), []byte("Ciphertext"), []byte("Hash"), []byte("Description"), []byte("Owner"), []byte(""), getEmptySigs(t)})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	preArgs := []string{"Ciphertext", "Hash", "Description", PubHexForTest}
	hash, signatures := getHashAndSignatures(t, preArgs)
	response1 := stub.MockInvoke("txID", [][]byte{[]byte(tee.MethodUpload), []byte("Ciphertext"), []byte("Hash"), []byte("Description"), []byte(PubHexForTest), []byte(hex.EncodeToString(hash)), signatures})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)
	return response
}

func filedEmptyError(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInvoke("1", args)
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.Contains(t, response.Message, "argument must be a non-empty string")
}

func payloadIsOk(t *testing.T, payload []byte) tee.SharedData {
	var data tee.SharedData
	err := json.Unmarshal(payload, &data)
	assert.NoError(t, err)
	dataIsOk(t, data)
	return data
}

func dataIsOk(t *testing.T, data tee.SharedData) {
	assert.Equal(t, "Ciphertext", data.Ciphertext)
	assert.Equal(t, "Hash", data.Hash)
	assert.Equal(t, "Description", data.Description)
	assert.Equal(t, "Owner", data.Owner)
	assert.NotNil(t, data.ID)
	assert.NotZero(t, data.UpdateSecondsTimestamp)
}

func Test_Invoke_update(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpdate), []byte(""), []byte("Ciphertext"), []byte("Hash"), []byte("Description")})
	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpdate), []byte("ID"), []byte(""), []byte("Hash"), []byte("Description")})
	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpdate), []byte("ID"), []byte("Ciphertext"), []byte(""), []byte("Description")})
	filedEmptyError(t, stub, [][]byte{[]byte(tee.MethodUpdate), []byte("ID"), []byte("Ciphertext"), []byte("Hash"), []byte("")})

	uploadData(t, stub, "1")
	updateNotExistIDError(t, stub)

	updateData(t, stub)
	bs, err := stub.GetState("1")
	assert.NoError(t, err)
	assert.NotEmpty(t, bs)

	var data tee.SharedData
	err1 := json.Unmarshal(bs, &data)
	assert.NoError(t, err1)
	assert.Equal(t, "Ciphertext1", data.Ciphertext)
	assert.Equal(t, "Summary1", data.Hash)
	assert.Equal(t, "Description1", data.Description)

	// repeated updates is no difference, except for the version added 1
	updateData(t, stub)
	bs1, err2 := stub.GetState("1")
	assert.NoError(t, err2)
	assert.NotEmpty(t, bs1)

	var data1 tee.SharedData
	err4 := json.Unmarshal(bs1, &data1)
	assert.NoError(t, err4)
	assert.Equal(t, data.Ciphertext, data1.Ciphertext)
	assert.Equal(t, data.Hash, data1.Hash)
	assert.Equal(t, data.Description, data1.Description)
}

func updateNotExistIDError(t *testing.T, stub *shim.MockStub) pb.Response {
	response := stub.MockInvoke("2", [][]byte{[]byte(tee.MethodUpdate), []byte("Ciphertext"), []byte("Hash"), []byte("Description"), []byte("Owner")})
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.NotNil(t, response.Message, "Shared data does not exist, id: Ciphertext")
	return response
}

func updateData(t *testing.T, stub *shim.MockStub) {
	args := [][]byte{[]byte(tee.MethodUpdate), []byte("1"), []byte("Ciphertext1"), []byte("Summary1"), []byte("Description1"), []byte(""), getEmptySigs(t)}
	response := stub.MockInvoke("3", args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.Nil(t, response.Payload)

	preArgs := []string{"1", "Ciphertext1", "Summary1", "Description1"}
	hash, signatures := getHashAndSignatures(t, preArgs)
	args1 := [][]byte{[]byte(tee.MethodUpdate), []byte("1"), []byte("Ciphertext1"), []byte("Summary1"), []byte("Description1"), []byte(hex.EncodeToString(hash)), signatures}
	response1 := stub.MockInvoke(tee.MethodUpdate, args1)
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.Nil(t, response1.Payload)
}

func Test_Invoke_queryDataByID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	filedEmptyError(t, stub, [][]byte{[]byte("queryDataByID"), []byte("")})

	uploadData(t, stub, "1")
	queryNotExistIDError(t, stub)

	args := [][]byte{[]byte("queryDataByID"), []byte("1")}
	response := stub.MockInvoke("2", args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var data tee.SharedData
	err := json.Unmarshal(response.Payload, &data)
	assert.NoError(t, err)
	assert.Equal(t, "Ciphertext", data.Ciphertext)
	assert.Equal(t, "Hash", data.Hash)
	assert.Equal(t, "Description", data.Description)
	assert.Equal(t, "Owner", data.Owner)
	assert.NotNil(t, data.ID)
	assert.NotZero(t, data.UpdateSecondsTimestamp)
}

func queryNotExistIDError(t *testing.T, stub *shim.MockStub) {
	response := stub.MockInvoke("2", [][]byte{[]byte("queryDataByID"), []byte("Ciphertext")})
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.Equal(t, "Shared data does not exist, id: Ciphertext", response.Message)
}

func Test_Invoke_queryHistoryByDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	filedEmptyError(t, stub, [][]byte{[]byte("queryHistoryByDID"), []byte("")})
}

func Test_Invoke_queryDataByOwner(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	filedEmptyError(t, stub, [][]byte{[]byte("queryDataByOwner"), []byte("")})

	uploadData(t, stub, "1")
	queryNotExistOwnerError(t, stub)

	args := [][]byte{[]byte("queryDataByOwner"), []byte("Owner")}
	response := stub.MockInvoke("2", args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var dataList [][]byte
	err := json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 1)

	payloadIsOk(t, dataList[0])

	uploadData(t, stub, "3")
	response1 := stub.MockInvoke("4", args)
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)

	var dataList1 [][]byte
	err1 := json.Unmarshal(response1.Payload, &dataList1)
	assert.NoError(t, err1)
	assert.Len(t, dataList1, 2)

	payloadIsOk(t, dataList1[0])
	payloadIsOk(t, dataList1[1])
}

func queryNotExistOwnerError(t *testing.T, stub *shim.MockStub) pb.Response {
	response := stub.MockInvoke("2", [][]byte{[]byte("queryDataByOwner"), []byte("Ciphertext")})
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.NotNil(t, response.Message, "Shared data with owner does not exist: Ciphertext")
	return response
}

func getEmptySigs(t *testing.T) []byte {
	sigs := make([]string, 1)
	sigsBytes, err := json.Marshal(sigs)
	assert.NoError(t, err)
	return sigsBytes
}

func getHashAndSignatures(t *testing.T, args []string) ([]byte, []byte) {
	hash := chaincode.GetArgsHashBytes(args)
	privBytes, err := common.HexToBytes(PrivHexForTest)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := ecc.SignECDSA(privBytes, hash)
	assert.NoError(t, err)

	sigs := []string{signature}
	data, err := json.Marshal(sigs)
	assert.NoError(t, err)

	return hash, data
}
