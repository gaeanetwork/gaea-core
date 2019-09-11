package user

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/stretchr/testify/assert"
)

func Test_checkLen(t *testing.T) {
	// username
	err := checkUsernameLen("")
	assert.Error(t, err, "Invalid username length")

	username := common.GetRandomStringByLen(16)
	err = checkUsernameLen(username)
	assert.NoError(t, err)

	username = common.GetRandomStringByLen(32)
	err = checkUsernameLen(username)
	assert.NoError(t, err)

	username = common.GetRandomStringByLen(64)
	err = checkUsernameLen(username)
	assert.Error(t, err, "Invalid username length")

	// password
	err = checkPasswordLen("")
	assert.Error(t, err, "Invalid password length")

	password := common.GetRandomStringByLen(16)
	err = checkPasswordLen(password)
	assert.NoError(t, err)

	password = common.GetRandomStringByLen(32)
	err = checkPasswordLen(password)
	assert.NoError(t, err)

	password = common.GetRandomStringByLen(64)
	err = checkPasswordLen(password)
	assert.Error(t, err, "Invalid password length")
}
