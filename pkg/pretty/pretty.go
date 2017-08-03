package pretty

var Pad int = 2 // same default value with cmdline flag

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}
