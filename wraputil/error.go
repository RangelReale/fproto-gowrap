package wraputil

type ServiceErrorType int

const (
	SET_GENERIC ServiceErrorType = iota
	SET_IMPORT
	SET_EXPORT
	SET_CALL
)

type ServiceErrorHandler interface {
	HandleServiceError(errorType ServiceErrorType, err error) error
}

type ServiceErrorHandler_Default struct {
}

func (e *ServiceErrorHandler_Default) HandleServiceError(errorType ServiceErrorType, err error) error {
	return err
}
