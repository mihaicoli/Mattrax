package pkg

// ErrorHandler is used to report errors thrown by the subpackages.
// These errors are already handled so just log or report them to any instrumentation platform
var ErrorHandler func(errDescription string, err error)
