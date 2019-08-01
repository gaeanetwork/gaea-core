package azurehelper

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// environments
const (
	AzureStorageAccountEnvKey   = "AZURE_STORAGE_ACCOUNT"
	AzureStorageAccessKeyEnvKey = "AZURE_STORAGE_ACCESS_KEY"
)

// MakeContainerURL A ContainerURL represents a URL to the Azure Storage container allowing you to manipulate its blobs.
func MakeContainerURL(containerName string) (azblob.ContainerURL, error) {
	var containerURL azblob.ContainerURL
	if len(containerName) == 0 {
		return containerURL, fmt.Errorf("Error container name must be non-empty")
	}

	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := os.Getenv(AzureStorageAccountEnvKey), os.Getenv(AzureStorageAccessKeyEnvKey)
	if len(accountName) == 0 || len(accountKey) == 0 {
		return containerURL, fmt.Errorf("Either the %s or %s environment variable is not set", AzureStorageAccountEnvKey, AzureStorageAccessKeyEnvKey)
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return containerURL, fmt.Errorf("Invalid credentials with error: %s", err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	return azblob.NewContainerURL(*URL, p), nil
}

// DownloadFilesFromContainer download all files by filePath in azure container
func DownloadFilesFromContainer(containerName, filePath string) ([]byte, error) {
	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL, err := MakeContainerURL(containerName)
	if err != nil {
		return nil, err
	}

	// Here's how to download the blob
	ctx, blobURL := context.Background(), containerURL.NewBlobURL(filePath)
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}

	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 2})

	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return nil, err
	}

	return downloadedData.Bytes(), nil
}

// UploadFileToContainer upload a file to azure container
func UploadFileToContainer(containerName, filePath string, data []byte) error {
	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL, err := MakeContainerURL(containerName)
	if err != nil {
		return err
	}

	// You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that PutBlob can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// handleErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	ctx, blobURL := context.Background(), containerURL.NewBlockBlobURL(filePath)
	_, err = azblob.UploadBufferToBlockBlob(ctx, data, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})

	return err
}
