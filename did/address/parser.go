package address

import (
	"errors"
	"fmt"
)

// Parser resolves blockchian
type Parser struct {
	drivers map[string]Driver
}

var defaultParser Parser

func init() {
	defaultParser = Parser{drivers: make(map[string]Driver)}
	defaultParser.drivers[Bitcoin] = btcDriver{name: Bitcoin}
	defaultParser.drivers[Ethereum] = ethereumDriver{name: Ethereum}
}

func (p *Parser) resolve(address string) (string, error) {
	for _, driver := range p.drivers {
		blockchain, err := driver.resolve(address)
		if err == nil {
			return blockchain, err
		}
	}
	return "", errors.New("not find")
}

// Resolve find out the addres belone to which block chain
func Resolve(address string) (string, error) {
	return defaultParser.resolve(address)
}

func (p *Parser) register(netName string) (string, error) {
	d, ok := p.drivers[netName]
	if !ok {
		return "", fmt.Errorf("Not support (%s)", netName)
	}

	return d.createAddress()
}

// Register to a special net
func Register(netName string) (string, error) {
	return defaultParser.register(netName)
}
