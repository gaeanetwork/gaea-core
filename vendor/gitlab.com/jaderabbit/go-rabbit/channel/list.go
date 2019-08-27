package channel

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/scc/cscc"
	"github.com/hyperledger/fabric/peer/common"
	common2 "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

// List this peer joined the channel
func List() ([]*pb.ChannelInfo, error) {
	cf, err := MakeGeneralClientSupport(common.UndefinedParamValue, common.UndefinedParamValue)
	if err != nil {
		return nil, err
	}
	defer cf.EndorserClient.Close()

	return getChannels(cf)
}

func getChannels(cf *ClientSupport) ([]*pb.ChannelInfo, error) {
	var err error

	invocation := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]),
			ChaincodeId: &pb.ChaincodeID{Name: "cscc"},
			Input:       &pb.ChaincodeInput{Args: [][]byte{[]byte(cscc.GetChannels)}},
		},
	}

	var prop *pb.Proposal
	c, _ := cf.Signer.Serialize()
	prop, _, err = utils.CreateProposalFromCIS(common2.HeaderType_ENDORSER_TRANSACTION, "", invocation, c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot create proposal, due to %s", err))
	}

	var signedProp *pb.SignedProposal
	signedProp, err = utils.GetSignedProposal(prop, cf.Signer)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot create signed proposal, due to %s", err))
	}

	proposalResp, err := cf.EndorserClient.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed sending proposal, got %s", err))
	}

	if proposalResp.Response == nil || proposalResp.Response.Status != 200 {
		return nil, errors.New(fmt.Sprintf("Received bad response, status %d: %s", proposalResp.Response.Status, proposalResp.Response.Message))
	}

	var channelQueryResponse pb.ChannelQueryResponse
	err = proto.Unmarshal(proposalResp.Response.Payload, &channelQueryResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read channels list response, %s", err))
	}

	return channelQueryResponse.Channels, nil
}
