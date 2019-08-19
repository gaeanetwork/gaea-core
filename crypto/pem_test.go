package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gaeanetwork/gaea-core/did/crypto"
	"github.com/stretchr/testify/assert"
)

func Test_ExportPrivateKeytoPem(t *testing.T) {
	filePath := "Test_ExportPrivateKeytoPem"
	defer os.RemoveAll(filePath)

	priv, err := crypto.GenerateKey()
	assert.NoError(t, err)

	err1 := ExportPrivateKeytoPem(filepath.Join(filePath, PRIVFILE), crypto.FromECDSA(priv), true)
	assert.NoError(t, err1)

	// a messy file name
	err2 := ExportPrivateKeytoPem(filepath.Join(filePath, "阿什◇♂√█▉*♂☎☻☀▪√㏂⊙…‥ 顿飞*^%%`1432~^&**$#@738`"), crypto.FromECDSA(priv), true)
	assert.NoError(t, err2)

	// mkdir permission denied
	err3 := ExportPrivateKeytoPem(filepath.Join("/data", PRIVFILE), crypto.FromECDSA(priv), true)
	assert.Contains(t, err3.Error(), "permission denied")

	// create file permission denied
	err4 := ExportPrivateKeytoPem(filepath.Join("/", PRIVFILE), crypto.FromECDSA(priv), true)
	assert.Contains(t, err4.Error(), "permission denied")
}

func Test_ExportPublicKeytoPem(t *testing.T) {
	filePath := "ExportPublicKeytoPem"
	defer os.RemoveAll(filePath)

	priv, err := crypto.GenerateKey()
	assert.NoError(t, err)

	err1 := ExportPublicKeytoPem(filepath.Join(filePath, PRIVFILE), crypto.FromECDSAPub(&priv.PublicKey))
	assert.NoError(t, err1)
}
