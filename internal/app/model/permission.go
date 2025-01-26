package model

// Owned{Terminal{Write},Start,Stop}
type Permission struct {
	Id       int64  `gorm:"column:id;NOT NULL" json:"id" `
	Account  string `gorm:"column:account;NOT NULL" json:"account"`
	Pid      int32  `gorm:"column:pid;NOT NULL" json:"pid"`
	Owned    bool   `gorm:"column:owned;NOT NULL" json:"owned"`
	Start    bool   `gorm:"column:start;NOT NULL" json:"start"`
	Stop     bool   `gorm:"column:stop;NOT NULL" json:"stop"`
	Terminal bool   `gorm:"column:terminal;NOT NULL" json:"terminal"`
	Write    bool   `gorm:"column:write;NOT NULL" json:"write"`
	Log      bool   `gorm:"column:log;NOT NULL" json:"log"`
}

func (*Permission) TableName() string {
	return "permission"
}

type PermissionPo struct {
	Id       int64  `gorm:"column:id" json:"id"`
	Account  string `gorm:"column:account" json:"account"`
	Name     string `gorm:"column:name" json:"name"`
	Pid      int32  `gorm:"column:pid" json:"pid"`
	Owned    bool   `gorm:"column:owned" json:"owned"`
	Start    bool   `gorm:"column:start" json:"start"`
	Stop     bool   `gorm:"column:stop" json:"stop"`
	Terminal bool   `gorm:"column:terminal" json:"terminal"`
	Write    bool   `gorm:"column:write;NOT NULL" json:"write"`
	Log      bool   `gorm:"column:log;NOT NULL" json:"log"`
}
