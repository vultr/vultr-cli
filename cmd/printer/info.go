package printer

// Info intializes and returns a generic resource output struct
func Info(msg string) *Generic {
	return &Generic{Message: msg}
}
