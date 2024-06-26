package model

type ProcessInfo struct {
	Name      string `json:"name"`
	Uuid      int    `json:"uuid"`
	StartTime string `json:"startTime"`
	User      string `json:"user"`
	Usage     Usage  `json:"usage"`
	State     State  `json:"state"`
	TermType  string `json:"termType"`
}

type Usage struct {
	Cpu  []float64 `json:"cpu"`
	Mem  []float64 `json:"mem"`
	Time []string  `json:"time"`
}

type State struct {
	State uint8  `json:"state"`
	Info  string `json:"info"`
}
