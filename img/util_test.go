package img

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
)

type testCase struct {
	metric    string
	filePath1 string
	filePath2 string
	result    []float64
}

func newTestCase(metr, file1, file2 []byte, res []float64) testCase {
	return testCase{metric: string(metr), filePath1: string(file1), filePath2: string(file2), result: res}
}

var beginExp *regexp.Regexp
var finalExp *regexp.Regexp
var pyFloatExp *regexp.Regexp

const PYTHONICFLOATEXP string = `^[+\-]?(?:(?:0|[1-9]\d*)(?:\.\d*)?|\.\d+)(?:\d[eE][+\-]?\d+)|(?:\d[+\-]\d+j)$`
const BEGINEXP string = `^img[1-8].jpg img[1-8].jpg`

func init() {
	beginExp = regexp.MustCompile(BEGINEXP)
	finalExp = regexp.MustCompile(`^\(?` + PYTHONICFLOATEXP + `(\)|, )?`)
	pyFloatExp = regexp.MustCompile(PYTHONICFLOATEXP)
}

func importAllTestCases() (results map[string][]testCase, resultError error) {
	rdr := openTestCaseFile()
	var expectedBeginExp = regexp.MustCompile(BEGINEXP + ` [mse rmse psnr rmse\_sw uqi ssim ergas scc rase sam msssim vifp psnrb] : `)

	for errors.Is(resultError, nil) {
		if line, resultError := rdr.ReadBytes('\n'); !errors.Is(resultError, nil) {
			if matched := expectedBeginExp.Match(line); matched {
				newCase := importTestCase(line)
				r := results[newCase.metric]
				results[newCase.metric] = append(r, newCase)
			}
		}
	}

	if !errors.Is(resultError, io.EOF) {
		checkFatal(resultError)
	}
	return
}

func importTestCases(metric string) (results []testCase, resultError error) {
	rdr := openTestCaseFile()
	var expectedBeginExp *regexp.Regexp
	if len(metric) == 0 {
		expectedBeginExp = beginExp
	} else {
		expectedBeginExp = regexp.MustCompile(BEGINEXP + ` ` + metric + ` : `)
	}

	for errors.Is(resultError, nil) {
		if line, resultError := rdr.ReadBytes('\n'); !errors.Is(resultError, nil) {
			if matched := expectedBeginExp.Match(line); matched {
				results = append(results, importTestCase(line))
			}
		}
	}

	if !errors.Is(resultError, io.EOF) {
		checkFatal(resultError)
	}
	return
}

func openTestCaseFile() (rd *bufio.Reader) {
	rawData, err := os.Open("test/assets/results.txt")
	checkFatal(err)
	rd = bufio.NewReader(rawData)
	return
}

func importTestCase(line []byte) (result testCase) {
	cols := bytes.SplitN(bytes.TrimSpace(line), []byte{' '}, 6)
	if len(cols) < 2 {
		panic(cols)
	} else if len(cols) == 2 {
		result = newTestCase(nil, cols[0], cols[1], make([]float64, 0))
	} else if len(cols) > 4 {
		if !bytes.Equal(cols[3], []byte{':'}) {
			panic(cols)
		}

		floatByte := pyFloatExp.FindAll(finalExp.Find(line), 2)
		floatCache := make([]float64, len(floatByte))

		for i := 0; i < len(floatByte); i++ {
			if f, err := strconv.ParseFloat(string(floatByte[i]), 64); errors.Is(err, nil) {
				floatCache[i] = f
			} else {
				panic(cols)
			}
		}
		result = newTestCase(cols[2], cols[0], cols[1], floatCache)
	} else {
		panic(line)
	}
	return
}
