package channel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/common/configtx"
	localsigner "github.com/hyperledger/fabric/common/localmsp"
	"github.com/hyperledger/fabric/common/tools/configtxgen/encoder"
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/peer/common"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
	"gitlab.com/jaderabbit/go-rabbit/i18n"
)

// GenerateChannelTx generate channelID.tx by congtx.yaml's ProfileName
func GenerateChannelTx(profileName, channelID string) error {
	profileConfig := genesisconfig.Load(profileName, config.GetConfigPath())
	if err := doOutputChannelCreateTx(profileConfig, channelID); err != nil {
		return fmt.Errorf("Error on outputChannelCreateTx: %s, err: %s", channelID, err)
	}

	return nil
}

func doOutputChannelCreateTx(profileConfig *genesisconfig.Profile, channelID string) error {
	configtx, err := encoder.MakeChannelCreationTransaction(channelID, nil, profileConfig)
	if err != nil {
		logger.Errorf("Failed to MakeChannelCreationTransaction,err: %s", err)
		return err
	}

	path, err := getChannelPath(channelID)
	if err != nil {
		logger.Errorf("Failed to getChannelPath,err: %s", err)
		return err
	}

	path = filepath.Join(path, channelID+".tx")
	logger.Infof("Writing new channel configtx: %s", channelID, path)
	if err = ioutil.WriteFile(path, utils.MarshalOrPanic(configtx), 0644); err != nil {
		return fmt.Errorf("Error writing channel create tx: %s", err)
	}
	return nil
}

func getChannelPath(channelID string) (string, error) {
	channelPath := filepath.Join(config.DefaultConfig.RabbitDataPath, channelID)
	if err := os.MkdirAll(channelPath, os.ModePerm); err != nil {
		return "", err
	}

	return channelPath, nil
}

// CreateChannel create a chain by channelID.tx
func CreateChannel(channelID, ordererEndpoint string) error {
	if channelID == common.UndefinedParamValue {
		return errors.New("must supply channel ID")
	}

	if ordererEndpoint == common.UndefinedParamValue {
		ordererEndpoint = config.DefaultConfig.OrdererEndpoint
		logger.Warnf("The orderer endpoint is not specified, the default endpoint is used: %s", ordererEndpoint)
	}
	common.OrderingEndpoint = ordererEndpoint
	viper.Set("orderer.address", ordererEndpoint)

	cf, err := MakeCreateClientSupport(channelID)
	if err != nil {
		return err
	}
	defer cf.DeliverClient.Close()

	logger.Infof("create channel[ID:%s] by mspConfigPath: %s", channelID, viper.GetString("peer.mspConfigPath"))
	return executeCreate(cf)
}

func executeCreate(cf *ClientSupport) error {
	if err := GenerateChannelTx("SampleSingleMSPChannel", cf.ChannelID); err != nil {
		return errors.Errorf("Error generating %s.tx: %s", cf.ChannelID, err)
	}

	if err := sendCreateChainTransaction(cf); err != nil {
		logger.Errorf("Error sending create chain transaction: %s", err)
		if strings.Contains(err.Error(), "/Channel/Application at version 0, but got version 1") {
			return errors.Errorf("Error already exists for channel with the same channelID: %s", cf.ChannelID)
		}
		if strings.Contains(err.Error(), "/Channel/Application not satisfied: Failed to reach implicit threshold of 1 sub-policies, required 1 remaining") {
			return errors.Errorf("Error msp identifier, please create %s with admin msp node", cf.ChannelID)
		}
		return err
	}

	block, err := getGenesisBlock(cf)
	if err != nil {
		return err
	}

	b, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	path, err := getChannelPath(cf.ChannelID)
	if err != nil {
		err = errors.WithMessage(err, "Failed to getChannelPath")
		logger.Error(err)
		return err
	}

	return ioutil.WriteFile(filepath.Join(path, cf.ChannelID+".block"), b, os.ModePerm)
}

func sendCreateChainTransaction(cf *ClientSupport) error {
	path, err := getChannelPath(cf.ChannelID)
	if err != nil {
		logger.Errorf("Failed to getChannelPath,err: %s", err)
		return err
	}

	var chCrtEnv *cb.Envelope
	if chCrtEnv, err = createChannelFromConfigTx(filepath.Join(path, cf.ChannelID+".tx")); err != nil {
		return err
	}

	if chCrtEnv, err = signConfigTx(chCrtEnv, cf.ChannelID); err != nil {
		return err
	}

	broadcastClient, err := cf.BroadcastFactory()
	if err != nil {
		return errors.WithMessage(err, "error getting broadcast client")
	}

	defer broadcastClient.Close()
	return broadcastClient.Send(chCrtEnv)
}

func createChannelFromConfigTx(configTxFileName string) (*cb.Envelope, error) {
	cftx, err := ioutil.ReadFile(configTxFileName)
	if err != nil {
		return nil, fmt.Errorf("channel create configuration tx file not found %s", err)
	}

	return utils.UnmarshalEnvelope(cftx)
}

func signConfigTx(envConfigUpdate *cb.Envelope, channelID string) (*cb.Envelope, error) {
	payload, err := utils.ExtractPayload(envConfigUpdate)
	if err != nil {
		return nil, i18n.InvalidCreateTxErr("bad payload")
	}

	if err = checkPayload(payload, channelID); err != nil {
		return nil, err
	}

	configUpdateEnv, err := configtx.UnmarshalConfigUpdateEnvelope(payload.Data)
	if err != nil {
		return nil, i18n.InvalidCreateTxErr("Bad config update env")
	}

	signer := localsigner.NewSigner()
	sigHeader, err := signer.NewSignatureHeader()
	if err != nil {
		return nil, err
	}

	configSig := &cb.ConfigSignature{
		SignatureHeader: utils.MarshalOrPanic(sigHeader),
	}

	configSig.Signature, err = signer.Sign(util.ConcatenateBytes(configSig.SignatureHeader, configUpdateEnv.ConfigUpdate))

	configUpdateEnv.Signatures = append(configUpdateEnv.Signatures, configSig)

	return utils.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, signer, configUpdateEnv, 0, 0)
}

func checkPayload(payload *cb.Payload, channelID string) error {
	if payload.Header == nil || payload.Header.ChannelHeader == nil {
		return i18n.InvalidCreateTxErr("bad header")
	}

	ch, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return i18n.InvalidCreateTxErr("could not unmarshall channel header")
	}

	if ch.Type != int32(cb.HeaderType_CONFIG_UPDATE) {
		return i18n.InvalidCreateTxErr("bad type")
	}

	if ch.ChannelId == "" || channelID == "" {
		return i18n.InvalidCreateTxErr("empty channel id")
	}

	if ch.ChannelId != channelID {
		return i18n.InvalidCreateTxErr(fmt.Sprintf("mismatched channel ID %s != %s", ch.ChannelId, channelID))
	}

	return nil
}

func getGenesisBlock(cf *ClientSupport) (*cb.Block, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return nil, errors.New("timeout waiting for channel creation")
		default:
			block, err := cf.DeliverClient.GetSpecifiedBlock(0)
			if err != nil {
				if cf, err = MakeCreateClientSupport(cf.ChannelID); err != nil {
					return nil, errors.WithMessage(err, "failed connecting")
				}

				time.Sleep(200 * time.Millisecond)
				continue
			}

			return block, nil
		}
	}
}
