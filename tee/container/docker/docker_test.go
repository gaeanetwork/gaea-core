package docker

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func Test_Create(t *testing.T) {
	c := New()
	err := c.Create()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.address)
}

func Test_Upload(t *testing.T) {
	c := New()
	algorithm, err := ioutil.ReadFile("/home/rabbit/teetest/client/resume")
	assert.NoError(t, err)
	A, err := ioutil.ReadFile("/home/rabbit/teetest/A/A_resume.txt")
	assert.NoError(t, err)
	B, err := ioutil.ReadFile("/home/rabbit/teetest/B/B_resume.txt")
	assert.NoError(t, err)
	C, err := ioutil.ReadFile("/home/rabbit/teetest/C/C_resume.txt")
	assert.NoError(t, err)
	dataList := [][]byte{A, B, C}
	err = c.Upload(algorithm, dataList)
	assert.NoError(t, err)

	os.Remove("main")
	os.Remove("0")
	os.Remove("1")
	os.Remove("2")
}
