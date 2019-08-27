package chaincode

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric/core/common/ccpackage"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

func signpackage(cfg *Config) error {
	if len(cfg.ChaincodeName) == 0 {
		return errors.New("not specified chaincode name")
	}

	ipackageFile := fmt.Sprintf("%spack.out", cfg.ChaincodeName)
	opackageFile := fmt.Sprintf("signed%spack.out", cfg.ChaincodeName)

	chaincodepath, err := GetChaincodePackagePath()
	if err != nil {
		return err
	}

	ipackageFile = filepath.Join(chaincodepath, ipackageFile)
	opackageFile = filepath.Join(chaincodepath, opackageFile)

	cfg.CommandName = "signpackage"
	cf, err := InitCmdFactory(false, false, cfg)
	if err != nil {
		return err
	}
	defer cf.Close()

	b, err := ioutil.ReadFile(ipackageFile)
	if err != nil {
		return err
	}

	env := utils.UnmarshalEnvelopeOrPanic(b)

	env, err = ccpackage.SignExistingPackage(env, cf.Signer)
	if err != nil {
		return err
	}

	b = utils.MarshalOrPanic(env)
	err = ioutil.WriteFile(opackageFile, b, 0700)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote signed package to %s successfully\n", opackageFile)

	return nil
}
