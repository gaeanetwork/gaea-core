package tee

import "github.com/gaeanetwork/gaea-core/smartcontract"

const (
	// ChaincodeName is the name of the tee core chaincode in the blockchain.
	ChaincodeName = "tee_data"

	// MethodUpload is the name of the method that uploads a shared data to the chaincode.
	MethodUpload = "upload"
	// MethodUpdate is the name of the method that updates a shared data to the chaincode.
	MethodUpdate = "update"
	// MethodRequest is the name of the method that requests a shared data to the chaincode.
	MethodRequest = "request"
	// MethodAuthorize is the name of the method that authorizes a requester with a shared data to the chaincode.
	MethodAuthorize = "authorize"

	// MethodQueryDataByID is the name of the method that gets a shared data by id from the chaincode.
	MethodQueryDataByID = "queryDataByID"
	// MethodQueryHistoryByDID is the name of the method that gets history of the shared data by data id from the
	// chaincode.
	MethodQueryHistoryByDID = "queryHistoryByDID"
	// MethodQueryDataByOwner is the name of the method that gets all shared data of owner from the chaincode.
	MethodQueryDataByOwner = "queryDataByOwner"
	// MethodQueryRequestsByRequesterAndDID is the name of the method that gets all the requests by requester and data
	// id from the chaincode.
	MethodQueryRequestsByRequesterAndDID = "queryRequestsByRequesterAndDID"
	// MethodQueryRequestsByRequesterAndStatusAndDID is the name of the method that gets all the requests by requester
	// and request status and data id from the chaincode.
	MethodQueryRequestsByRequesterAndStatusAndDID = "queryRequestsByRequesterAndStatusAndDID"
	// MethodQueryNotificationsByOwnerAndDID is the name of the method that gets all the notifications by owner and data
	// id from the chaincode.
	MethodQueryNotificationsByOwnerAndDID = "queryNotificationsByOwnerAndDID"
	// MethodQueryNotificationsByOwnerAndRequesterAndDID is the name of the method that gets all the notifications by
	// owner and requester and data id a task from the chaincode.
	MethodQueryNotificationsByOwnerAndRequesterAndDID = "queryNotificationsByOwnerAndRequesterAndDID"
	// MethodQueryNotificationsByOwnerAndStatusAndDID is the name of the method that gets all the notifications by
	// owner and notification status and data id from the chaincode.
	MethodQueryNotificationsByOwnerAndStatusAndDID = "queryNotificationsByOwnerAndStatusAndDID"
	// MethodQueryNotificationsByOwnerAndRequesterAndStatusAndDID is the name of the method gets all the notifications
	// by owner and requester and notifcation status and data id from the chaincode.
	MethodQueryNotificationsByOwnerAndRequesterAndStatusAndDID = "queryNotificationsByOwnerAndRequesterAndStatusAndDID"

	// ImplementPlatform is implement platform
	ImplementPlatform = smartcontract.Fabric
)
