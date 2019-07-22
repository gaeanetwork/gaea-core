package address

import "testing"

func TestParse(t *testing.T) {
	chain, err := Resolve("1N2f642sbgCMbNtXFajz9XDACDFnFzdXzV")
	if err != nil {
		t.Errorf("not find %s\n", err)
	}
	if chain != Bitcoin {
		t.Errorf("find error chain %s\n", chain)
	}

	chain, err = Resolve("24602722816b6cad0e143ce9fabf31f6026ec622")
	if err != nil {
		t.Errorf("not find %s\n", err)
	}
	if chain != Ethereum {
		t.Errorf("find error chain %s\n", chain)
	}

	chain, err = Resolve("0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359")
	if err != nil {
		t.Errorf("not find %s\n", err)
	}
	if chain != Ethereum {
		t.Errorf("find error chain %s\n", chain)
	}
}

func TestGenETHAddress(t *testing.T) {
	d := ethereumDriver{name: Ethereum}
	addr, err := d.createAddress()
	if err != nil {
		t.Errorf("get error %s\n", err)
	}
	chain, err := d.resolve(addr)
	if err != nil {
		t.Errorf("resolve error: %s\n", err)
	}
	if chain != Ethereum {
		t.Errorf("find error chain: %s, %s\n", chain, addr)
	}
}

func TestGenBTCAddress(t *testing.T) {
	d := btcDriver{name: Bitcoin}

	addr, err := d.createAddress()
	if err != nil {
		t.Errorf("get error %s\n", err)
	}
	chain, err := d.resolve(addr)
	if err != nil {
		t.Errorf("resolve error: %s\n", err)
	}
	if chain != Bitcoin {
		t.Errorf("find error chain: %s, %s\n", chain, addr)
	}
}
