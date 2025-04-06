package loggeradapter

//go:generate mockgen -source=adapter.go -package=loggeradapter -destination=adapter_mock.go
type Adapter interface {
	NewInfo(logString string)
	NewError(logString string)
}
