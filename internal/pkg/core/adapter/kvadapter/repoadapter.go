package kvadapter

// RepoAdapter defines an interface for a key value store
//go:generate mockgen -source=repoadapter.go -package=kvadapter -destination=repoadapter_mock.go
type RepoAdapter interface {
	SetValueUntilChannelClose(key string, data string, ttl int, isOpen *bool)
	Delete(key string) (err error)
	GetValue(key string) (res string, err error)
}
