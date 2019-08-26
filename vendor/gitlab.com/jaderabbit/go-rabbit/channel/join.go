package channel

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric/core/scc/cscc"
	"github.com/hyperledger/fabric/peer/common"
	pcommon "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	putils "github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"gitlab.com/jaderabbit/go-rabbit/i18n"
	"gitlab.com/jaderabbit/go-rabbit/src"
)

// JoinResult the peer join result
type JoinResult struct {
	Err         error  `json:"err"`
	PeerAddress string `json:"peerAddress"`
}

// JoinChain for peer
func JoinChain(channelID string, peerAddresses ...string) error {
	if channelID == common.UndefinedParamValue {
		return errors.New("must supply channel ID")
	}

	count := len(peerAddresses)
	if count == 0 {
		return errors.New("this channel must contain at least one peer")
	}

	joinResultC := make(chan JoinResult, count)
	for _, peerAddress := range peerAddresses {
		go func(peerAddress string) {
			cf, err := MakeGeneralClientSupport(channelID, peerAddress)
			if err != nil {
				joinResultC <- JoinResult{Err: err, PeerAddress: peerAddress}
				return
			}
			defer cf.EndorserClient.Close()

			logger.Infof("peer[address:%s] join channel[ID:%s] by mspConfigPath: %s", peerAddress, channelID, viper.GetString("peer.mspConfigPath"))
			if err = executeJoin(cf, channelID); err != nil {
				if strings.Contains(err.Error(), "LedgerID already exists") {
					joinResultC <- JoinResult{Err: i18n.PeerAlreadyJoinedErr(channelID), PeerAddress: peerAddress}
					return
				}
			}

			src.InsertActionLog(types.PeerJoinChannel,
				"system",
				fmt.Sprintf("peer[address:%s] join channel[ID:%s] by mspConfigPath: %s",
					peerAddress,
					channelID,
					viper.GetString("peer.mspConfigPath")))

			joinResultC <- JoinResult{Err: i18n.PeerJoinSuccessfully(peerAddress), PeerAddress: peerAddress}
		}(peerAddress)
	}
	succeed := true
	results := make([]JoinResult, 0)
	for joinResult := range joinResultC {
		count--
		results = append(results, joinResult)
		if _, ok := joinResult.Err.(i18n.PeerJoinSuccessfully); !ok {
			succeed = false
		}
		if count <= 0 {
			close(joinResultC)
		}
	}

	if !succeed {
		logger.Errorf("Error peers join failed. error: %s", results)
		return fmt.Errorf("%s", results)
	}

	return nil
}

func executeJoin(cf *ClientSupport, channelID string) (err error) {
	spec, err := getJoinCCSpec(channelID)
	if err != nil {
		return err
	}

	// Build the ChaincodeInvocationSpec message
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := cf.Signer.Serialize()
	if err != nil {
		return fmt.Errorf("Error serializing identity for %s: %s", cf.Signer.GetIdentifier(), err)
	}

	var prop *pb.Proposal
	prop, _, err = putils.CreateProposalFromCIS(pcommon.HeaderType_CONFIG, "", invocation, creator)
	if err != nil {
		return fmt.Errorf("Error creating proposal for join %s", err)
	}

	var signedProp *pb.SignedProposal
	signedProp, err = putils.GetSignedProposal(prop, cf.Signer)
	if err != nil {
		return fmt.Errorf("Error creating signed proposal %s", err)
	}

	var proposalResp *pb.ProposalResponse
	proposalResp, err = cf.EndorserClient.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return i18n.ProposalFailedErr(err.Error())
	}

	if proposalResp == nil {
		return i18n.ProposalFailedErr("nil proposal response")
	}

	if proposalResp.Response.Status != 0 && proposalResp.Response.Status != 200 {
		return i18n.ProposalFailedErr(fmt.Sprintf("bad proposal response %d: %s", proposalResp.Response.Status, proposalResp.Response.Message))
	}
	logger.Infof("Successfully submitted proposal to join channel: %s", channelID)
	return nil
}

func getJoinCCSpec(channelID string) (*pb.ChaincodeSpec, error) {
	path, err := getChannelPath(channelID)
	if err != nil {
		logger.Errorf("Failed to getChannelPath,err: %s", err)
		return nil, err
	}

	gb, err := ioutil.ReadFile(filepath.Join(path, channelID+".block"))
	if err != nil {
		return nil, i18n.GBFileNotFoundErr(err.Error())
	}
	// Build the spec
	input := &pb.ChaincodeInput{Args: [][]byte{[]byte(cscc.JoinChain), gb}}

	spec := &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]),
		ChaincodeId: &pb.ChaincodeID{Name: "cscc"},
		Input:       input,
	}

	return spec, nil
}
