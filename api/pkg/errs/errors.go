package errs

// Err implements the Error interface with error marshaling.
type Err struct {
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Details map[string]string `json:"details"`
}

func New(message, code string) *Err {
	return &Err{Message: message, Code: code}
}

func (e *Err) Error() string {
	return e.Message
}

// IsExpected finds Err{} inside passed error.
func IsExpected(err error) bool {
	_, ok := err.(*Err)
	return ok
}

// GetCode returns code of given error or empty string if error is not custom
func GetCode(err error) string {
	v, ok := err.(*Err)
	if !ok {
		return ""
	}
	return v.Code
}
