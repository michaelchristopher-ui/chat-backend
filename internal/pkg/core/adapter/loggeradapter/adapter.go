package loggeradapter

//go:generate mockgen -source=adapter.go -package=loggeradapter -destination=adapter_mock.go
type Adapter interface {
	NewInfo(logString string, params ...any)
	NewError(logString string, params ...any)
}
