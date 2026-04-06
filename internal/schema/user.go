package schema

type UserAuth struct {
	ID       uint64
	Password string
}

type UserLevel struct {
	Level uint8
	Exp   uint64
}

type Equipment interface {
	GetID() uint64
}

type UserInfo struct {
	UserAuth
	UserLevel
	EquipmentList []Equipment
}
