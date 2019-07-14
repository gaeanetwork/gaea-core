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
}
