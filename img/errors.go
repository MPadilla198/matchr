package img

type testError struct{}

func checkFatal(e error) {
	if e != nil {
		panic(e)
	}
}
