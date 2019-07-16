package address

import (
	"errors"
)

// Parser resolves blockchian
type Parser struct {
	drivers map[string]Driver
}

var defaultParser Parser

func init() {
	defaultParser = Parser{drivers: make(map[string]Driver)}
	defaultParser.drivers[Bitcoin] = btcDriver{name: Bitcoin}
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
