package constants

type TerminalType string

const (
	TERMINAL_PTY TerminalType = "pty"
	TERMINAL_STD TerminalType = "std"
)

type ProcessState int32

const (
	PROCESS_STOP ProcessState = iota
	PROCESS_START
	PROCESS_WARNNING
)
