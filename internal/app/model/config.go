package model

type Config struct {
	Id    int     `gorm:"column:id;primary_key"`
	Key   string  `gorm:"column:key"`
	Value *string `gorm:"column:value"`
}

func (n *Config) TableName() string {
	return "config"
}

type SystemConfigurationVo struct {
	Key      string `json:"key"`
	Value    any    `json:"value"`
	Default  string `json:"default"`
	Describe string `json:"describe"`
}
