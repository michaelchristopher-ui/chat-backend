package wsconnadapter

//go:generate mockgen -source=adapter.go -package=wsconnadapter -destination=adapter_mock.go
type Adapter interface {
	WriteJSON(v interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
}
