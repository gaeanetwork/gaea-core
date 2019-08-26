package src

import (
	"fmt"
	"time"

	"gitlab.com/jaderabbit/go-rabbit/chaincode/system"
	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"gitlab.com/jaderabbit/go-rabbit/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertActionLog(action types.Action, userID string, parameter ...string) {
	actionLogConn, err := mongodb.GetActionLogConnection()
	if err != nil {
		logger.Errorf("failed to get the action log connection,  err:%s", err.Error())
		return
	}

	actionLog, err := types.NewActionLog(action, userID)
	if err != nil {
		logger.Errorf("failed to new the action log,  err:%s", err.Error())
		return
	}

	if len(parameter) > 0 {
		actionLog.Mark = parameter[0]
	}

	if len(parameter) > 1 {
		actionLog.SessionID = parameter[1]
	}

	if err := actionLogConn.Insert(actionLog); err != nil {
		logger.Errorf("failed to insert the action log,  err:%s", err.Error())
	}
}

func QueryActionLogCount(action types.Action, userID string, beginTime, endTime int64) (int64, error) {
	actionLogConn, err := mongodb.GetActionLogConnection()
	if err != nil {
		return 0, fmt.Errorf("failed to get the action log connection,  err:%s", err.Error())
	}

	query := bson.M{}

	showName := action.String()
	if len(showName) > 0 {
		query["action"] = action
	}

	if len(userID) > 0 {
		query["user_id"] = userID
	}

	if beginTime > 0 && endTime > 0 {
		query["timestamp"] = bson.M{"$gte": beginTime, "$lte": endTime}
	} else if beginTime > 0 {
		query["timestamp"] = bson.M{"$gte": beginTime}
	} else if endTime > 0 {
		query["timestamp"] = bson.M{"$lte": endTime}
	}

	return actionLogConn.QueryCount(query)
}

func QueryActionLogList(action types.Action, userID string, beginTime, endTime, pageIndex, pageSize int64) (*system.QueryResult, error) {
	actionLogConn, err := mongodb.GetActionLogConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get the action log connection,  err:%s", err.Error())
	}

	query := bson.M{}

	showName := action.String()
	if len(showName) > 0 {
		query["action"] = action
	}

	if len(userID) > 0 {
		query["user_id"] = userID
	}

	if beginTime > 0 && endTime > 0 {
		query["timestamp"] = bson.M{"$gte": beginTime, "$lte": endTime}
	} else if beginTime > 0 {
		query["timestamp"] = bson.M{"$gte": beginTime}
	} else if endTime > 0 {
		query["timestamp"] = bson.M{"$lte": endTime}
	}
	skip := (pageIndex - 1) * pageSize

	arrayLogs, err := actionLogConn.GetByFilter(query, skip, pageSize)
	if err != nil {
		return nil, err
	}

	count, err := actionLogConn.QueryCount(query)
	if err != nil {
		return nil, err
	}

	result := &system.QueryResult{
		TotalNumber: int(count),
		Data:        arrayLogs,
		PageIndex:   int(pageIndex),
		PageSize:    int(pageSize),
	}
	return result, nil
}

func QueryActionLog(id string) (*types.ActionLog, error) {
	actionLogConn, err := mongodb.GetActionLogConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get the action log connection,  err:%s", err.Error())
	}

	return actionLogConn.Get(id)
}

func QueryLasetActionLog(sessionID string) (*types.ActionLog, error) {
	actionLogConn, err := mongodb.GetActionLogConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get the action log connection,  err:%s", err.Error())
	}

	dt := time.Now().Add(-1 * time.Hour)
	query := bson.M{"action": types.Login.String(), "session_id": sessionID, "timestamp": bson.M{"$gte": dt.Unix()}}

	arrayLogs, err := actionLogConn.GetByFilter(query, 0, 1)
	if err != nil {
		return nil, err
	}

	if len(arrayLogs) == 0 {
		return nil, nil
	}

	return arrayLogs[0], nil
}
