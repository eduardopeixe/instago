package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	AlertMsgGeneric = "something went wrong. Please try again"
)

// Alert is used to render Bootstrap Alert messages
type Alert struct {
	Level   string
	Message string
}

// Data is the top level struct that views expect data to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}
