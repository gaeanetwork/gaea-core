package keyagreement

import "crypto"

const (
	// DH is Diffie-Hellman Key Agreement as defined in PKCS #3: Diffie-Hellman Key-Agreement Standard, RSA
	// Laboratories, version 1.4, November 1993.
	DH string = "DH"
	// ECDH is Elliptic Curve Diffie-Hellman as defined in ANSI X9.63 and as described in RFC 3278: "Use of Elliptic
	// Curve Cryptography (ECC) Algorithms in Cryptographic Message Syntax (CMS)."
	ECDH string = "ECDH"
	// ECMQV is Elliptic Curve Menezes-Qu-Vanstone.
	ECMQV string = "ECMQV"
)

// KeyAgreement provides the functionality of a key agreement (or key exchange) protocol.
type KeyAgreement interface {
	GetAlgorithm() string
	GenerateSharedSecret(crypto.PrivateKey, crypto.PublicKey) ([]byte, error)
}
