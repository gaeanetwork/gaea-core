package i18n

//InvalidCreateTxErr invalid channel create transaction
type InvalidCreateTxErr string

func (e InvalidCreateTxErr) Error() string {
	return GetDefaultPrinter().Sprintf("Invalid channel create transaction : %s", string(e))
}

// Code defines the error code
func (e InvalidCreateTxErr) Code() string {
	return "300001"
}

//GBFileNotFoundErr genesis block file not found
type GBFileNotFoundErr string

func (e GBFileNotFoundErr) Error() string {
	return GetDefaultPrinter().Sprintf("genesis block file not found %s", string(e))
}

// Code defines the error code
func (e GBFileNotFoundErr) Code() string {
	return "300002"
}

//ProposalFailedErr proposal failed
type ProposalFailedErr string

func (e ProposalFailedErr) Error() string {
	return GetDefaultPrinter().Sprintf("proposal failed (err: %s)", string(e))
}

// Code defines the error code
func (e ProposalFailedErr) Code() string {
	return "300003"
}

// ChannelNotExistsErr channel not exists
type ChannelNotExistsErr string

func (e ChannelNotExistsErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error does not exists this channel. id: %s", string(e))
}

// Code defines the error code
func (e ChannelNotExistsErr) Code() string {
	return "300004"
}

// ChannelKeyEmptyErr channel key is empty
type ChannelKeyEmptyErr string

func (e ChannelKeyEmptyErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error channel key[%s] is empty", string(e))
}

// Code defines the error code
func (e ChannelKeyEmptyErr) Code() string {
	return "300005"
}
