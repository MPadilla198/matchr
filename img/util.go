package img

import (
	"bufio"
	"bytes"
	"errors"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/*************************************************
*
* COLOR UTILITIES
*
***************************************************/

type Color interface {
	color.Color
	EEEA() (e1, e2, e3 float64)
}

type Model interface {
	Convert(c color.Color) Color
}

type model struct {
	f func(color.Color) Color
}

func (m model) Convert(c color.Color) Color {
	return m.f(c)
}

func ModelFunc(f func(color.Color) Color) Model {
	return model{f: f}
}

var (
	// HSLModel https://en.wikipedia.org/wiki/HSL_and_HSV
	HSLModel = ModelFunc(hslModel)
	HSVModel = ModelFunc(hsvModel)
	// HSPModel https://alienryderflex.com/hsp.html
	HSPModel = ModelFunc(hspModel)
	// Luma601Model https://en.wikipedia.org/wiki/HSL_and_HSV#Lightness
	Luma601Model  = ModelFunc(luma601Model)
	Luma240Model  = ModelFunc(luma240Model)
	Luma709Model  = ModelFunc(luma709Model)
	Luma2020Model = ModelFunc(luma2020Model)
)

func hslModel(c color.Color) Color {
	if hslC, ok := c.(HSL); ok {
		return hslC
	} else if hsv, ok := c.(HSV); ok {
		l := hsv.v * (1 - (hsv.s / 2))
		var sl float32
		if l == 0 || l == 1 {
			sl = 0
		} else {
			sl = (hsv.v - l) / float32(math.Min(float64(l), float64(1.0-l)))
		}
		return HSL{hsv.h, sl, l}
	} // else if hsp, ok := c.(HSP); ok {}

	// TODO Convert to HSL from RGBA
	return nil
}

type HSL struct {
	h, s, l float32
}

func (c HSL) EEEA() (e1, e2, e3 float64) {
	e1 = float64(c.h)
	e2 = float64(c.s)
	e3 = float64(c.l)
	return
}

func (c HSL) RGBA() (r, g, b, a uint32) {
	return
}

func hsvModel(c color.Color) Color {
	if hsvC, ok := c.(HSV); ok {
		return hsvC
	} else if hsl, ok := c.(HSL); ok {
		v := hsl.l + hsl.s*(float32(math.Min(float64(hsl.l), float64(1.0-hsl.l))))
		var sv float32
		if v == 0 {
			sv = 0
		} else {
			sv = 2 * (1.0 - (hsl.l - v))
		}
		return HSL{hsl.h, sv, v}
	} // else if hsp, ok := c.(HSP); ok {}
	return nil
}

type HSV struct {
	h, s, v float32
}

func (c HSV) EEEA() (e1, e2, e3 float64) {
	e1 = float64(c.h)
	e2 = float64(c.s)
	e3 = float64(c.v)
	return
}

func (c HSV) RGBA() (r, g, b, a uint32) {
	return
}

func hspModel(c color.Color) Color {
	return nil
}

type HSP struct {
	h, s, p float32
}

func (c HSP) EEEA() (e1, e2, e3 float64) {
	e1 = float64(c.h)
	e2 = float64(c.s)
	e3 = float64(c.p)
	return
}

func (c HSP) RGBA() (r, g, b, a uint32) {
	return
}

func luma601Model(c color.Color) Color {
	r, g, b, a := c.RGBA()
	rfl := float64(r) / float64(a)
	gfl := float64(g) / float64(a)
	bfl := float64(b) / float64(a)
	return Luma601{L: (0.2989 + rfl) + (0.5870 * gfl) + (0.1140 * bfl)}
}

type Luma601 struct {
	L float64
}

func (c Luma601) EEEA() (e1, e2, e3 float64) {
	e1 = c.L
	e2 = c.L
	e3 = c.L
	return
}

func (c Luma601) RGBA() (r, g, b, a uint32) {
	a = 0xffff
	l := uint32((float64(a) * c.L) + 0.5)
	return l, l, l, a
}

func luma240Model(c color.Color) Color {
	r, g, b, a := c.RGBA()
	rfl := float64(r) / float64(a)
	gfl := float64(g) / float64(a)
	bfl := float64(b) / float64(a)
	return Luma240{L: (0.212 + rfl) + (0.701 * gfl) + (0.087 * bfl)}
}

type Luma240 struct {
	L float64
}

func (c Luma240) EEEA() (e1, e2, e3 float64) {
	e1 = c.L
	e2 = c.L
	e3 = c.L
	return
}

func (c Luma240) RGBA() (r, g, b, a uint32) {
	a = 0xffff
	l := uint32((float64(a) * c.L) + 0.5)
	return l, l, l, a
}

func luma709Model(c color.Color) Color {
	r, g, b, a := c.RGBA()
	rfl := float64(r) / float64(a)
	gfl := float64(g) / float64(a)
	bfl := float64(b) / float64(a)
	return Luma709{L: (0.2126 + rfl) + (0.7152 * gfl) + (0.0722 * bfl)}
}

type Luma709 struct {
	L float64
}

func (c Luma709) EEEA() (e1, e2, e3 float64) {
	e1 = c.L
	e2 = c.L
	e3 = c.L
	return
}

func (c Luma709) RGBA() (r, g, b, a uint32) {
	a = 0xffff
	l := uint32((float64(a) * c.L) + 0.5)
	return l, l, l, a
}

func luma2020Model(c color.Color) Color {
	r, g, b, a := c.RGBA()
	rfl := float64(r) / float64(a)
	gfl := float64(g) / float64(a)
	bfl := float64(b) / float64(a)
	return Luma2020{L: (0.2627 + rfl) + (0.6780 * gfl) + (0.0593 * bfl)}
}

type Luma2020 struct {
	L float64
}

func (c Luma2020) EEEA() (e1, e2, e3 float64) {
	e1 = c.L
	e2 = c.L
	e3 = c.L
	return
}

func (c Luma2020) RGBA() (r, g, b, a uint32) {
	a = 0xffff
	l := uint32((float64(a) * c.L) + 0.5)
	return l, l, l, a
}

/*************************************************
*
* IMAGE MATRIX UTILITIES
*
***************************************************/

const ASSETPATH = "test/assets/"

type decodeFunc func(r io.Reader) (image.Image, error)

// TODO REWORK imageMatrix to use less dense
type imageMatrix struct {
	image.Image

	c Model
}

func (im *imageMatrix) Dims() (r, c int) {
	min := im.Image.Bounds().Min
	max := im.Image.Bounds().Max
	r = max.Y - min.Y
	c = max.X - min.X
	return
}

func (im *imageMatrix) At(i, j int) (l float64) {
	min := im.Image.Bounds().Min
	_, _, l = im.c.Convert(im.Image.At(min.X+j, min.Y+i)).EEEA()
	return
}

func (im *imageMatrix) T() mat.Matrix {
	return im
}

func (im *imageMatrix) Compile() *mat.Dense {
	return mat.DenseCopyOf(im)
}

func newImageMatrix(given image.Image, cm Model) *imageMatrix {
	return &imageMatrix{given, cm}
}

func newImageMatrixFromFile(filePath string, cm Model) *imageMatrix {
	// log.Printf("\"%s\"", filepath.Ext(strings.TrimSpace(filePath)))
	return _newImageMatrixFromFile(filePath, map[string]decodeFunc{
		".jpg": jpeg.Decode,
		".png": png.Decode,
	}[filepath.Ext(strings.TrimSpace(filePath))], cm)
}

func _newImageMatrixFromFile(filePath string, decoder decodeFunc, cm Model) *imageMatrix {
	rawData, err := os.Open(filePath)
	checkFatal(err)
	decodedImg, err := decoder(bufio.NewReader(rawData))
	checkFatal(err)

	return &imageMatrix{decodedImg, cm}
}

/*************************************************
*
* TESTING UTILITIES
*
***************************************************/

const ERRMARGIN float64 = 0.000000000000000000001

type testCase struct {
	metric    string
	filePath1 []byte
	filePath2 []byte
	result    []float64
}

func equal(reference *testCase, given []float64) (isEqual bool) {
	if reference == nil || given == nil || len(given) != len(reference.result) {
		return false
	}
	// innocent until proven guilty
	isEqual = true
	for i := 0; isEqual && i < len(given); i++ {
		isEqual = isEqual && (given[i]-reference.result[i] <= ERRMARGIN)
	}
	return
}

func newTestCase(metr, file1, file2 []byte, res []float64) *testCase {
	return &testCase{metric: string(metr), filePath1: file1, filePath2: file2, result: res}
}

var _beginExp *regexp.Regexp
var _finalExp *regexp.Regexp
var _pyFloatExp *regexp.Regexp

const PYTHONICFLOATEXP string = `[+\-]?(?:(?:0|[1-9]\d*)(?:\.\d*)?|\.\d+)(?:(?:\d[eE][+\-]?\d+)|(?:\d[+\-]\d+j))?`
const BEGINEXP string = `^img[1-8].jpg img[1-8].jpg`
const DATAFILEPATH string = ASSETPATH + "results.txt"

func init() {
	_beginExp = regexp.MustCompile(BEGINEXP)
	_finalExp = regexp.MustCompile(`^\(?` + PYTHONICFLOATEXP + `(\)|, )?`)
	_pyFloatExp = regexp.MustCompile(PYTHONICFLOATEXP)
}

func importAllTestCases() (results map[string][]*testCase, resultError error) {
	rdr := openTestCaseFile()
	var expectedBeginExp = regexp.MustCompile(BEGINEXP + ` (mse)|(rmse)|(psnr)|(rmse\_sw)|(uqi)|(ssim)|(ergas)|(scc)|(rase)|(sam)|(msssim)|(vifp)|(psnrb) : `)

	results = make(map[string][]*testCase, 279)
	line, resultError := rdr.ReadBytes('\n')
	for errors.Is(resultError, nil) {
		// log.Printf("%S", line)
		if matched := expectedBeginExp.Match(line); matched {
			// log.Printf("Matched")
			newCase := importTestCase(line)
			r := results[newCase.metric]
			results[newCase.metric] = append(r, newCase)
		}
		line, resultError = rdr.ReadBytes('\n')
	}

	if errors.Is(resultError, io.EOF) {
		resultError = nil
	}
	return
}

func importMetricTestCases(metric string) (results []*testCase, resultError error) {
	rdr := openTestCaseFile()
	var expectedBeginExp *regexp.Regexp
	if len(metric) == 0 {
		expectedBeginExp = _beginExp
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
	rawData, err := os.Open(DATAFILEPATH)
	checkFatal(err)
	rd = bufio.NewReader(rawData)
	return
}

func importTestCase(line []byte) (result *testCase) {
	cols := bytes.SplitN(bytes.TrimSpace(line), []byte{' '}, 6)
	log.Printf("%s", cols)
	if len(cols) == 2 {
		result = newTestCase(nil, cols[0], cols[1], make([]float64, 0))
	} else if len(cols) == 5 {
		if !bytes.Equal(cols[3], []byte{':'}) {
			panic(cols)
		}

		if string(cols[4]) == "inf" {
			result = newTestCase(cols[2], cols[0], cols[1], []float64{math.Inf(1)})
		} else if string(cols[4]) == "nan" {
			result = newTestCase(cols[2], cols[0], cols[1], []float64{math.NaN()})
		} else if f, err := strconv.ParseFloat(string(_pyFloatExp.Find(cols[4])), 64); errors.Is(err, nil) {
			result = newTestCase(cols[2], cols[0], cols[1], []float64{f})
		} else {
			panic(err)
		}
	} else if len(cols) == 6 {
		if !bytes.Equal(cols[3], []byte{':'}) {
			panic(cols)
		}

		floatByte := _pyFloatExp.FindAll(_finalExp.Find(line), 2)
		log.Printf("float Byte: %v", floatByte)
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
