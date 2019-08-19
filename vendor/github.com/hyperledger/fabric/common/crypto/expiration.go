/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/msp"
)

var (
	rwMutex            sync.RWMutex
	expiredDuration    = 5 * time.Second
	cachedExpiresAtMap = make(map[[32]byte]cachedExpiresAt)
)

type cachedExpiresAt struct {
	cachedAt  time.Time
	expiresAt time.Time
}

// ExpiresAt returns when the given identity expires, or a zero time.Time
// in case we cannot determine that
func ExpiresAt(identityBytes []byte) time.Time {
	cached := sha256.Sum256(identityBytes)
	rwMutex.RLock()
	cachedObj, exists := cachedExpiresAtMap[cached]
	rwMutex.RUnlock()
	if exists && time.Since(cachedObj.cachedAt) < expiredDuration {
		return cachedObj.expiresAt
	}

	sId := &msp.SerializedIdentity{}
	// If protobuf parsing failed, we make no decisions about the expiration time
	if err := proto.Unmarshal(identityBytes, sId); err != nil {
		return time.Time{}
	}
	bl, _ := pem.Decode(sId.IdBytes)
	if bl == nil {
		// If the identity isn't a PEM block, we make no decisions about the expiration time
		return time.Time{}
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return time.Time{}
	}
	rwMutex.Lock()
	cachedExpiresAtMap[cached] = cachedExpiresAt{time.Now(), cert.NotAfter}
	rwMutex.Unlock()

	return cert.NotAfter
}
