package azurehelper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	accountForTests       = "azureteeaccount"
	accessKeyForTests     = "0jlvPT+Gw2+Y4ltGfuOXCkw91QQI82gsL2RjHbQPCmo7VlDneujTFnu+B7a/FC5tfVdCkVCZRti1Dpfw0Evaaw=="
	containerNameForTests = "tee-container"
	filePathForTests      = "teetest/A/A_resume.txt"
)

func init() {
	initForTests()
}

func initForTests() {
	os.Setenv(AzureStorageAccountEnvKey, accountForTests)
	os.Setenv(AzureStorageAccessKeyEnvKey, accessKeyForTests)
}

func Test_DownloadFilesFromContainer(t *testing.T) {
	fileBytes, err := DownloadFilesFromContainer(containerNameForTests, filePathForTests)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileBytes)

	tempName := "container_1.txt"
	err = ioutil.WriteFile(tempName, fileBytes, os.ModePerm)
	assert.NoError(t, err)

	err = os.Remove(tempName)
	assert.NoError(t, err)
}

func Test_DownloadFilesFromContainer_ErrContainerName(t *testing.T) {
	_, err := DownloadFilesFromContainer("containerNameForTests", filePathForTests)
	assert.Contains(t, err.Error(), "The specifed resource name contains invalid characters")

	_, err = DownloadFilesFromContainer("containernamefortests", filePathForTests)
	assert.Contains(t, err.Error(), "The specified container does not exist")

	_, err = DownloadFilesFromContainer("", filePathForTests)
	assert.Contains(t, err.Error(), "Error container name must be non-empty")
}

func Test_DownloadFilesFromContainer_ErrFilePath(t *testing.T) {
	_, err := DownloadFilesFromContainer(containerNameForTests, "filePathForTests")
	assert.Contains(t, err.Error(), "The specified blob does not exist")

	_, err = DownloadFilesFromContainer(containerNameForTests, "")
	assert.Contains(t, err.Error(), "The requested URI does not represent any resource on the server")
}

func Test_DownloadFilesFromContainer_ErrEnv(t *testing.T) {
	// os.Setenv(AzureStorageAccountEnvKey, "accessKeyForTests")
	// assert.PanicsWithValue(t, "test timed out after 30s", func() { DownloadFilesFromContainer(containerNameForTests, filePathForTests) })

	os.Setenv(AzureStorageAccessKeyEnvKey, "accessKeyForTests")
	_, err := DownloadFilesFromContainer(containerNameForTests, filePathForTests)
	assert.Contains(t, err.Error(), "Invalid credentials with error: illegal base64 data")

	os.Setenv(AzureStorageAccountEnvKey, "")
	_, err = DownloadFilesFromContainer(containerNameForTests, filePathForTests)
	assert.Contains(t, err.Error(), "Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")

	os.Setenv(AzureStorageAccountEnvKey, accountForTests)
	os.Setenv(AzureStorageAccessKeyEnvKey, "")
	_, err = DownloadFilesFromContainer(containerNameForTests, filePathForTests)
	assert.Contains(t, err.Error(), "Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
}

func Test_UploadFileToContainer(t *testing.T) {
	initForTests()
	var buffer = []byte("Hello World!")
	err := UploadFileToContainer(containerNameForTests, "upload.txt", buffer)
	assert.NoError(t, err)

	// Invalid container name
	err = UploadFileToContainer("containerNameForTests", "upload.txt", buffer)
	assert.Contains(t, err.Error(), "The specifed resource name contains invalid characters")
}
