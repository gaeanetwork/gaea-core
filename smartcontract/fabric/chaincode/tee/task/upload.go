package main

import "log"

var (
	hooks = make(map[string]UploadHook)
)

// UploadHook upload result execution log to data store address
type UploadHook func(executionLog []byte) error

func uploadResults(executionLog []byte) error {
	log.Println("results length: ", len(hooks))
	for owner, upload := range hooks {
		log.Println("start to upload result of:", owner)
		if err := upload(executionLog); err != nil {
			return err
		}

		log.Println("successfully uploading.")
	}

	return nil
}
