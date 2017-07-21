package pretty

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}

/*
func WriteToFile(f string, contents []byte) error {
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("output written to %s", f))
	return nil
}
*/
