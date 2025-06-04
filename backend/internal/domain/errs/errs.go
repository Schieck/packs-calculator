package errs

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

func FormatValidationErrors(err error) string {
	return err.Error()
}
