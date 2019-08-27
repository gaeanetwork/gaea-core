package models

// Organization for view
type Organization struct {
	Name        string `json:"name"`
	Contract    string `json:"contract"`
	Phone       string `json:"phone"`
	PeersNumber int32  `json:"peersNum"`
}

// GetLocalOrg get local organization
func GetLocalOrg() (*Organization, error) {
	return &Organization{
		Name:        "echain",
		Contract:    "Wensi Liu",
		Phone:       "15711329686",
		PeersNumber: 1,
	}, nil
}
