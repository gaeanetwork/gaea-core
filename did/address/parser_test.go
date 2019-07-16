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
}
