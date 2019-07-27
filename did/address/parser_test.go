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

func TestVerifyBTCSign(t *testing.T) {
	publicKey := "02a673638cb9587cb68ea08dbef685c" +
		"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5"
	signatureSerialize := "30450220090ebfb3690a0ff115bb1b38b" +
		"8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e707" +
		"1cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479"

	d := btcDriver{name: Bitcoin}
	res, err := d.verifySign(signatureSerialize, publicKey)
	if err != nil {
		t.Errorf("verify sign failed, error: %s\n", err)
	}
	if res == false {
		t.Error("verfiy error")
	}
}

func TestVerifyETHSign(t *testing.T) {
	publicKey := "0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652"
	signatureSerialize := "0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301"

	d := ethereumDriver{name: Ethereum}
	res, err := d.verifySign(signatureSerialize, publicKey)
	if err != nil {
		t.Errorf("verify sign failed, error: %s\n", err)
	}
	if res == false {
		t.Error("verfiy error")
	}
}
