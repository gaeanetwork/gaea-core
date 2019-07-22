package smartcontract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Platform_String(t *testing.T) {
	var platform = Fabric
	assert.Equal(t, "fabric", platform.String())

	platform = Ethereum
	assert.Equal(t, "ethereum", platform.String())

	platform = Platform(867)
	assert.Equal(t, "fabric", platform.String())
}
