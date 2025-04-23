package passwordgeneratoradapter

//go:generate mockgen -source=adapter.go -package=passwordgeneratoradapter -destination=adapter_mock.go
type Adapter interface {
	Generate(pass string) (string, error)
}
