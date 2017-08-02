package utils

// CheckError is a test helper that determines if an expected and actual error
// are the same.
func CheckError(want, got error) bool {
	if want == nil && got == nil {
		return true
	} else if want != nil && got != nil {
		return want.Error() == got.Error()
	}

	return false
}
