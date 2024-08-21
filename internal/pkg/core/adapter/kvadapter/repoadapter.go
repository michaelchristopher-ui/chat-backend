package kvadapter

// RepoAdapter defines an interface for a key value store
type RepoAdapter interface {
	SetValueUntilChannelClose(key string, data string, ttl int, isOpen *bool)
	Delete(key string) (err error)
	GetValue(key string) (res string, err error)
}
