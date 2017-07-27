package check

var devMode bool
var dbgMode bool

func DevMode() *bool {
	return &devMode
}

func IsDevMode() bool {
	return devMode
}

func DbgMode() *bool {
	return &dbgMode
}

func IsDbgMode() bool {
	return dbgMode
}
