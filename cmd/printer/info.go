package printer

// Info intializes and returns a Message resource output struct
func Info(msg string) *Message {
	return &Message{Message: msg}
}
