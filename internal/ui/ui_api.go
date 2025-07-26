package ui

type KeeperApp interface {
	Register(Username string, Password string, Email string) error
	Login(Username string, Password string) error
	Encrypt(Password string, Input string, Output string) error
	Decrypt(Password string, Input string, Output string) error
	Upload(Output string, Category string) error
	Download(Input string) error
	Delete(Input string) error
	Card(Number string, Date string, CVV string, Password string) error
	Deletecard(Input string) error
	Sync() error
}

type ScreenAPI interface {
}

type LoginScreenAPI interface {
	Login(Username string, Password string) error
}

type RegisterScreenAPI interface {
	Register(Username string, Password string, Email string) error
}

type EncryptScreenAPI interface {
	Encrypt(Password string, Input string, Output string) error
}

type DecryptScreenAPI interface {
	Decrypt(Password string, Input string, Output string) error
}

type DeleteScreenAPI interface {
	Delete(Input string) error
}

type CardScreenAPI interface {
	Card(Number string, Date string, CVV string, Password string) error
}

type DeletecardScreenAPI interface {
	Delete(Input string) error
}

type SyncScreenAPI interface {
	Sync() error
}
