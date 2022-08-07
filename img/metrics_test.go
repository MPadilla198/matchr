package img

import (
	"errors"
	"fmt"
	"testing"
)

func TestMetrics(t *testing.T) {
	tCases, err := importAllTestCases()
	if !errors.Is(err, nil) {
		t.Error(tCases, err)
	}

	var j = 0
	for metricName, metricFunc := range ImageMetrics {
		// log.Printf("%s", metricName)
		var cases []*testCase = tCases[metricName]
		// log.Print(len(cases))
		for i := 0; i < len(cases); i++ {
			// log.Printf("Case #%d) %s(%s, %s)", i, metricName, cases[i].filePath1, cases[i].filePath2)
			t.Run(fmt.Sprintf("Case #%d) %s(%s, %s)", i+(40*j), metricName, cases[i].filePath1, cases[i].filePath2), func(t *testing.T) {
				got, err := metricFunc(cases[i].filePath1, cases[i].filePath2)
				if !errors.Is(err, nil) {
					t.Fatalf("metricFunc error: %s", err)
				} else if !equal(cases[i], got) {
					t.Fail()
				}
			})
		}
		j++
	}
}
