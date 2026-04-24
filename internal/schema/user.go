package schema

type UserAuth struct { // for user authentication
	ID    uint64
	Token string
}

type UserRegister struct { // for user registration
	ID       uint64
	UserName string
	Password string
}

type UserLevel struct {
	Level uint8
	Exp   uint64
}

type UserInfo struct {
	UserRegister
}
