package chaincode

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/core/common/ccpackage"
	"github.com/hyperledger/fabric/msp"
	mspmgmt "github.com/hyperledger/fabric/msp/mgmt"
	pcommon "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
)

const packageCmdName = "package"
const packageDesc = "Package the specified chaincode into a deployment spec."

func chaincodePackage(cfg *Config) error {
	if len(cfg.ChaincodeName) == 0 {
		return errors.New("not specified chaincode name")
	}

	if len(cfg.ChaincodeVersion) == 0 {
		return errors.New("not specified chaincode version")
	}

	if len(cfg.ChaincodePath) == 0 {
		return errors.New("not specified path")
	}

	packageName := fmt.Sprintf("%spack.out", cfg.ChaincodeName)
	chaincodepath, err := GetChaincodePackagePath()
	if err != nil {
		return err
	}

	packagePath := filepath.Join(chaincodepath, packageName)

	cfg.CommandName = "package"
	cf, err := InitCmdFactory(false, false, cfg)
	if err != nil {
		return err
	}
	defer cf.Close()

	spec, err := cfg.CreateChaincodeSpec()
	if err != nil {
		return err
	}

	cds, err := getChaincodeDeploymentSpec(spec, true)
	if err != nil {
		return fmt.Errorf("error getting chaincode code %s: %s", cfg.ChaincodeName, err)
	}

	bytesToWrite, err := getChaincodeInstallPackage(cds, cf)
	if err != nil {
		return err
	}

	logger.Debugf("Packaged chaincode into deployment spec of size <%d>, with args = %s", len(bytesToWrite), packagePath)
	err = ioutil.WriteFile(packagePath, bytesToWrite, 0700)
	if err != nil {
		logger.Errorf("failed writing deployment spec to file [%s]: [%s]", packagePath, err)
		return err
	}

	return nil
}

//getChaincodeInstallPackage returns either a raw ChaincodeDeploymentSpec or
//a Envelope with ChaincodeDeploymentSpec and (optional) signature
func getChaincodeInstallPackage(cds *pb.ChaincodeDeploymentSpec, cf *ChaincodeCmdFactory) ([]byte, error) {
	//this can be raw ChaincodeDeploymentSpec or Envelope with signatures
	var objToWrite proto.Message

	//start with default cds
	objToWrite = cds

	var err error

	var owner msp.SigningIdentity

	if cf.Signer == nil {
		return nil, fmt.Errorf("Error getting signer")
	}
	owner = cf.Signer

	mspid, err := mspmgmt.GetLocalMSP().GetIdentifier()
	if err != nil {
		return nil, err
	}
	ip := "AND('" + mspid + ".admin')"

	sp, err := getInstantiationPolicy(ip)
	if err != nil {
		return nil, err
	}

	//we get the Envelope of type CHAINCODE_PACKAGE
	objToWrite, err = ccpackage.OwnerCreateSignedCCDepSpec(cds, sp, owner)
	if err != nil {
		return nil, err
	}

	//convert the proto object to bytes
	bytesToWrite, err := proto.Marshal(objToWrite)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling chaincode package : %s", err)
	}

	return bytesToWrite, nil
}

func getInstantiationPolicy(policy string) (*pcommon.SignaturePolicyEnvelope, error) {
	p, err := cauthdsl.FromString(policy)
	if err != nil {
		return nil, fmt.Errorf("Invalid policy %s, err %s", policy, err)
	}
	return p, nil
}

func GetChaincodePackagePath() (string, error) {
	chaincodepath := filepath.Join(config.GetConfigPath(), "chaincodepack")
	if !common.FileOrFolderExists(chaincodepath) {
		err := os.MkdirAll(chaincodepath, 0777)
		if err != nil {
			return "", fmt.Errorf("failed to create path, path:%s, err:%s", chaincodepath, err.Error())
		}
	}
	return chaincodepath, nil
}
