package cmd

import (
	"crypto/tls"
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/client"
	"github.com/hyperledger/fabric/msp"
)

// ChaincodeCmdFactory holds the clients used by ChaincodeCmd
type ChaincodeCmdFactory struct {
	EndorserClients []*client.EndorserClient
	DeliverClients  []*client.DeliverClient
	Certificate     tls.Certificate
	Signer          msp.SigningIdentity
	BroadcastClient *client.BroadcastClient
}

// Close the EndorserClients and DeliverClients of ChaincodeCmdFactory
func (cf *ChaincodeCmdFactory) Close() error {
	eLen, dLen := len(cf.EndorserClients), len(cf.DeliverClients)
	if eLen != dLen {
		return fmt.Errorf("EndorserClients length[%d] is not equals to DeliverClients[%d]", eLen, dLen)
	}

	for index := 0; index < eLen; index++ {
		cf.EndorserClients[index].Close()
		cf.DeliverClients[index].Close()
	}

	if cf.BroadcastClient != nil {
		cf.BroadcastClient.Close()
	}

	return nil
}
