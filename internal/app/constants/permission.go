package constants

type OprPermission string

const (
	OPERATION_START          OprPermission = "Start"
	OPERATION_STOP           OprPermission = "Stop"
	OPERATION_TERMINAL       OprPermission = "Terminal"
	OPERATION_TERMINAL_WRITE OprPermission = "Write"
	OPERATION_LOG            OprPermission = "Log"
)
