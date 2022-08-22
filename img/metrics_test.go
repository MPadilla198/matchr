package img

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

var ASSETPATHBYTES []byte // const

func init() {
	ASSETPATHBYTES = []byte(`test/assets/`)
}

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
				fileBuilder1 := new(strings.Builder)
				_, err := fileBuilder1.Write(ASSETPATHBYTES)
				if !errors.Is(err, nil) {
					panic("Test String Pool: Error in New().")
				}
				fileBuilder2 := new(strings.Builder)
				_, err = fileBuilder2.Write(ASSETPATHBYTES)
				if !errors.Is(err, nil) {
					panic("Test String Pool: Error in New().")
				}
				_, err = fileBuilder1.Write(cases[i].filePath1)
				if !errors.Is(err, nil) {
					t.Errorf("metricFunc error: %s", err)
					return
				}
				_, err = fileBuilder2.Write(cases[i].filePath2)
				if !errors.Is(err, nil) {
					t.Errorf("metricFunc error: %s", err)
					return
				}
				img1 := newImageMatrixFromFile(fileBuilder1.String(), Luma709Model)
				img2 := newImageMatrixFromFile(fileBuilder2.String(), Luma709Model)
				got := metricFunc(img1, img2)
				if !equal(cases[i], got) {
					t.Logf("expected: %v", cases[i].result)
					t.Logf("got: %v", got)
					t.Fail()
					return
				}
			})
			caseCount++
		}
	}
}
