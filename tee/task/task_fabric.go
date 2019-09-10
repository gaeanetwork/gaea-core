package task

const (
	// ChaincodeName is the name of the tee task chaincode in the blockchain.
	ChaincodeName = "tee_exec"

	// MethodGet is the name of the method that gets a task from the chaincode.
	MethodGet = "get"
	// MethodGetAll is the name of the method that gets all the tasks from the chaincode.
	MethodGetAll = "getAll"
	// MethodCreate is the name of the method that creates a task to the chaincode.
	MethodCreate = "create"
	// MethodExecute is the name of the method that executes a task to the chaincode.
	MethodExecute = "exectue"

	// AlgorithmName is the name of the algorithm in the docker container.
	AlgorithmName = "main"

	// AzureSplitSep is used to compose azure container name and filepath to the DataStoreAddress if the DataStoreType is Azure.
	AzureSplitSep = "~"

	// TaskIDIndex is to query all the task.
	TaskIDIndex = "task~id"

	// KeyPrivHex is to save the privHex to status.
	KeyPrivHex = "task~privHex"
	// KeyPubHex is to save the pubHex to status.
	KeyPubHex = "task~pubHex"

	// DefaultResultAddress is the default address where the result of the tee task is calculated in the container.
	DefaultResultAddress = "/tmp/teedata"
)
