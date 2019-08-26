package i18n

// InvokeResponseFailedErr invoke succeeded but the response was not 200
type InvokeResponseFailedErr string

func (e InvokeResponseFailedErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error invoke succeeded but the response status was not 200 (response: %s)", string(e))
}

// Code defines the error code
func (e InvokeResponseFailedErr) Code() string {
	return "400001"
}

// BusinessNotFoundErr cannot find the business
type BusinessNotFoundErr string

func (e BusinessNotFoundErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error cannot find the business: %s", string(e))
}

// Code defines the error code
func (e BusinessNotFoundErr) Code() string {
	return "400002"
}

// BusinessEmptyErr emtpy business
type BusinessEmptyErr string

func (e BusinessEmptyErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error emtpy business: %s", string(e))
}

// Code defines the error code
func (e BusinessEmptyErr) Code() string {
	return "400003"
}

// BusinessVersionTooLowErr business version is too low
type BusinessVersionTooLowErr struct {
	SrcV  string
	DestV string
}

func (e BusinessVersionTooLowErr) Error() string {
	return GetDefaultPrinter().Sprintf("Error business version too low: src[%s] --> dest[%s]", e.SrcV, e.DestV)
}

// Code defines the error code
func (e BusinessVersionTooLowErr) Code() string {
	return "400004"
}
