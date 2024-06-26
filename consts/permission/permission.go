package permission

type OprPermission string

const (
	START_OPERATION    OprPermission = "Start"
	STOP_OPERATION     OprPermission = "Stop"
	TERMINAL_OPERATION OprPermission = "Terminal"
)
