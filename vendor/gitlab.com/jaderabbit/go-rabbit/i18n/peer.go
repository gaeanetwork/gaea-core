package i18n

//PeerNotFoundErr cannot find the peer
type PeerNotFoundErr string

func (e PeerNotFoundErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error finding the peer. address: %s", string(e))
}

// Code defines the error code
func (e PeerNotFoundErr) Code() string {
	return "100001"
}

//PeerNotContainErr cannot find the peer
type PeerNotContainErr string

func (e PeerNotContainErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error this peer address is not included in the channel peers. address: %s", string(e))
}

// Code defines the error code
func (e PeerNotContainErr) Code() string {
	return "100002"
}

//PeerJoinSuccessfully peer joined successfully
type PeerJoinSuccessfully string

func (e PeerJoinSuccessfully) Error() string {
	return GetDefaultPrinter().Sprintf("Peer joined successfully. address: %s", string(e))
}

// Code defines the error code
func (e PeerJoinSuccessfully) Code() string {
	return "100003"
}

//PeerAlreadyJoinedErr peer already joined this channel
type PeerAlreadyJoinedErr string

func (e PeerAlreadyJoinedErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error peer already joined this channel: %s", string(e))
}

// Code defines the error code
func (e PeerAlreadyJoinedErr) Code() string {
	return "100004"
}
