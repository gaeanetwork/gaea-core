package smartcontract

// Service interface for kinds of platforms such as Ethereum, Fabric to invoke/query smart contract.
type Service interface {
	Invoke(contractID string, arguments []string) (result []byte, err error)
	Query(contractID string, arguments []string) (result []byte, err error)
	GetPlatform() Platform
}

// Platform is smart contract service platform
type Platform int

const (
	// Fabric is hyperledger fabric chaincode
	Fabric Platform = iota
	// Ethereum is ethereum evm solidity
	Ethereum
)

func (p Platform) String() string {
	var s string
	switch p {
	case Fabric:
		s = "fabric"
	case Ethereum:
		s = "ethereum"
	default:
		s = "fabric"
	}

	return s
}
