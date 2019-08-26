package sdk

// GetDefaultServer get the user chaincode server
func GetDefaultServer() (*Server, error) {
	return GetSDKServer("systemchaincode", DefaultMspUserName)
}

// GetDefaultUserServer get the user chaincode server
func GetDefaultUserServer() (*Server, error) {
	return GetSDKServer("user", DefaultMspUserName)
}

// GetUserServer get the user chaincode server
func GetUserServer(mspUserName string) (*Server, error) {
	return GetSDKServer("user", mspUserName)
}

// GetUnionServer get the union chaincode server
func GetUnionServer(mspUserName string) (*Server, error) {
	return GetSDKServer("union", mspUserName)
}

// GetRoleServer get the role chaincode server
func GetRoleServer(mspUserName string) (*Server, error) {
	return GetSDKServer("role", mspUserName)
}

// GetRelationServer get the relation chaincode server
func GetRelationServer(mspUserName string) (*Server, error) {
	return GetSDKServer("relation", mspUserName)
}

// GetMenuServer get the menu chaincode server
func GetMenuServer(mspUserName string) (*Server, error) {
	return GetSDKServer("menu", mspUserName)
}

// GetIPListServer get the iplist chaincode server
func GetIPListServer(mspUserName string) (*Server, error) {
	return GetSDKServer("iplist", mspUserName)
}

// GetDepartmentServer get the department chaincode server
func GetDepartmentServer(mspUserName string) (*Server, error) {
	return GetSDKServer("department", mspUserName)
}

// GetBusinessServer get the business chaincode server
func GetBusinessServer(mspUserName string) (*Server, error) {
	return GetSDKServer("business", mspUserName)
}

// GetAssetServer get the asset chaincode server
func GetAssetServer(mspUserName string) (*Server, error) {
	return GetSDKServer("asset", mspUserName)
}

// GetExampleServer get the example chaincode server
func GetExampleServer(mspUserName string) (*Server, error) {
	return GetSDKServer("example", mspUserName)
}
