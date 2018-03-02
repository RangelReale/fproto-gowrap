package wraputil

type ServiceErrorType int

const (
	SET_GENERIC ServiceErrorType = iota
	SET_IMPORT
	SET_EXPORT
	SET_CALL
)

// Error customization.
type ServiceErrorHandler interface {
	HandleServiceError(errorType ServiceErrorType, err error) error
}

// Default error customization only output the same error.
type ServiceErrorHandler_Default struct {
}

func (e *ServiceErrorHandler_Default) HandleServiceError(errorType ServiceErrorType, err error) error {
	return err
}
