package util

type ErrorType int

const (
	VALIDATION_ERROR ErrorType = iota
	APPLICATION_ERROR
)

// GetStatusCodeByErrorType gets Http Status code by ErrorType enum.
func GetStatusCodeByErrorType(errorType ErrorType) int {
	statusCodeMap := map[ErrorType]int{VALIDATION_ERROR: 400, APPLICATION_ERROR: 500}
	if statusCode := statusCodeMap[errorType]; statusCode != 0 {
		return statusCode
	} else {
		return 500
	}
}
