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

	var caseCount = 0
	for metricName, metricFunc := range ImageMetrics {
		var cases = tCases[metricName]
		for i := 0; i < len(cases); i++ {
			t.Run(fmt.Sprintf("Case #%d) %s(%s, %s)", caseCount, metricName, cases[i].filePath1, cases[i].filePath2), func(t *testing.T) {
				got, err := metricFunc(cases[i].filePath1, cases[i].filePath2)
				if !errors.Is(err, nil) {
					t.Fatalf("metricFunc error: %s", err)
				} else if !equal(cases[i], got) {
					t.Fail()
				}
			})
			caseCount++
		}
	}
}
