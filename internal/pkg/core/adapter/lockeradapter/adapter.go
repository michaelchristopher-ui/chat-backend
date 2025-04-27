package lockeradapter

import _ "github.com/golang/mock/mockgen/model" //When generating mocks in reflect mode, this workaround should be used.

//go:generate mockgen -package=lockeradapter -destination=adapter_mock.go sync Locker
//Instead of using our own adapter, we'll be using the built in sync.Locker interface
