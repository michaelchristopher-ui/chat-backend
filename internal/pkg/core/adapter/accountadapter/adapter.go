package accountadapter

// Adapter defines an interface for the account feature
type Adapter interface {
	VerifyAuth(VerifyAuthReq) error
	Register(RegisterReq) error
}

//VerifyAuthReq defines the struct input of the VerifyAuth function
type VerifyAuthReq struct {
	UserID   string
	Password string
}

//RegisterReq defines the struct input of the Register function
type RegisterReq struct {
	UserID   string
	Password string
}
