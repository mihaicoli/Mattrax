package pkg

// ErrorHandler is used to report errors thrown by the subpackages.
// These errors are already handled so just log or report them to any instrumentation platform
var ErrorHandler func(errDescription string, err error) // TODO: Remove this. Replace with normal, good old error handling

// AdvancedError is a custom implementation of the error type with support for the context required for a SOAP fault & normal logging
type AdvancedError struct {
	Err                                 error
	InternalDescription                 string
	FaultCauser, FaultType, FaultReason string
}

// Error returns the error string. This satisfies the error interface.
func (e AdvancedError) Error() string {
	return e.Err.Error()
}
