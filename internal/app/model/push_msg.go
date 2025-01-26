package model

type Push struct {
	Id     int64  `gorm:"column:id;NOT NULL" json:"id"`
	Method string `gorm:"column:method;NOT NULL" json:"method"`
	Url    string `gorm:"column:url;NOT NULL" json:"url"`
	Body   string `gorm:"column:body;NOT NULL" json:"body"`
	Remark string `gorm:"column:remark;NOT NULL" json:"remark"`
	Enable bool   `gorm:"column:enable;NOT NULL" json:"enable"`
}

func (*Push) TableName() string {
	return "push"
}
