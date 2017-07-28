package cli

var dbgMode bool

func DbgMode() *bool {
	return &dbgMode
}

func IsDbgMode() bool {
	return dbgMode
}
