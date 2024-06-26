package model

type Process struct {
	Uuid        int    `gorm:"primaryKey;autoIncrement;column:uuid" json:"uuid"`
	Name        string `gorm:"column:name" json:"name"`
	Cmd         string `gorm:"column:args" json:"cmd"`
	Cwd         string `gorm:"column:cwd" json:"cwd"`
	AutoRestart bool   `gorm:"column:auto_restart" json:"autoRestart"`
	Push        bool   `gorm:"column:push" json:"push"`
	LogReport   bool   `gorm:"column:log_report" json:"logReport"`
	TermType    string `gorm:"column:term_type" json:"termType"`
}

func (*Process) TableName() string {
	return "process"
}
