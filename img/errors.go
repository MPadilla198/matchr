package img

func checkFatal(e error) {
	if e != nil {
		panic(e)
	}
}
