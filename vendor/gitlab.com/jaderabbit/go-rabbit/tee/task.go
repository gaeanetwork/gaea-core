package tee

import "gitlab.com/jaderabbit/go-rabbit/tee/container"

// Task to create and execute the trusted execution task.
// Partners are all the data owners and the algorithm owner.
type Task struct {
	ID                     string              `json:"id"`
	DataIDs                []string            `json:"dataIDs"`
	AlgorithmID            string              `json:"algorithmID"`
	Container              container.Type      `json:"containerType"`
	ResultAddress          string              `json:"resultAddress"`
	DataNotifications      map[string]string   `json:"dataNotifications"`
	CreateSecondsTimestamp int64               `json:"createSeconds"`
	UploadSecondsTimestamp int64               `json:"uploadSeconds"`
	EvidenceHash           EvidenceHash        `json:"evidenceHash"`
	Requester              string              `json:"requester"`
	Partners               map[string]struct{} `json:"partners"`
	Executor               string              `json:"executor"`
}

// EvidenceHash to store tee input and output result hash
type EvidenceHash struct {
	DataHash         []string `json:"data"`
	AlgorithmHash    string   `json:"algorithm"`
	ResultHash       string   `json:"result"`
	ExecutionLogHash string   `json:"log"`
}
