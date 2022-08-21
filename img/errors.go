package img

import "fmt"

type utilError struct {
	name string
	err  error
}

func (e *utilError) Error() string {
	return fmt.Sprintf("error in %s:\n%s", e.name, e.err)
}

func (e *utilError) Unwrap() error {
	return e.err
}

type metricError struct {
	name string
	err  error
}

func (e *metricError) Error() string {
	return fmt.Sprintf("error in %s:\n%s", e.name, e.err)
}

func (e *metricError) Unwrap() error {
	return e.err
}

func checkFatal(e error) {
	if e != nil {
		panic(e)
	}
}
