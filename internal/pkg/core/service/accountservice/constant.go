package accountservice

// This contains the format of the account service logs. Add new account service format strings here.
const (
	logErrRegisterFormat = "[AccountService][Register] %s, err:%s"
	logInfoRegistered    = "[AccountService][Register] Account with UserID %s registered"

	logErrVerifyAuthFormat = "[AccountService][VerifyAuth] err:%s"
)
