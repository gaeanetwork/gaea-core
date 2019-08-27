package ecc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/gaeanetwork/gaea-core/crypto/keyagreement"
	"github.com/pkg/errors"
)

// ECDH is Elliptic Curve Diffie-Hellman as defined in ANSI X9.63 and as described in RFC 3278: "Use of Elliptic
// Curve Cryptography (ECC) Algorithms in Cryptographic Message Syntax (CMS)."
//
// see more detail: https://www.ietf.org/rfc/rfc3278.txt
type ECDH struct{}

// GetAlgorithm returns the algorithm name of this KeyAgreement object.
func (ecdh *ECDH) GetAlgorithm() string {
	return keyagreement.ECDH
}

// GenerateSharedSecret creates the shared secret and returns it as a sha256 hashed object.
func (ecdh *ECDH) GenerateSharedSecret(priv crypto.PrivateKey, pub crypto.PublicKey) ([]byte, error) {
	privateKey, ok := priv.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("priv only support ecdsa.PrivateKey point type")
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("pub only support ecdsa.PublicKey point type")
	}
	curve := publicKey.Curve.Params()
	x, _ := curve.ScalarMult(publicKey.X, publicKey.Y, privateKey.D.Bytes())

	fmt.Println("x:", x.Bytes())
	sharedKey := sha256.Sum256(x.Bytes())
	return sharedKey[:], nil
}
