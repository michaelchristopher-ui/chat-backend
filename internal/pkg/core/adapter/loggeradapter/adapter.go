package loggeradapter

type Adapter interface {
	NewInfo(logString string)
	NewError(logString string)
}
