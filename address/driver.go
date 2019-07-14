package address

// Block Chain Names
const (
	UnKnown = "unknown"
	Bitcoin = "btc"
)

// Driver interface
type Driver interface {
	resolve(address string) (string, error)
}

type btcDriver struct {
	name string
}

func (btc btcDriver) resolve(address string) (string, error) {
	return Bitcoin, nil
}
