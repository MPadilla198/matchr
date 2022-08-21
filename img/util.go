package img

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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

// https://en.wikipedia.org/wiki/HSL_and_HSV
// https://alienryderflex.com/hsp.html
var (
	HSLModel      color.Model = color.ModelFunc(hslModel)
	HSVModel      color.Model = color.ModelFunc(hsvModel)
	HSPModel      color.Model = color.ModelFunc(hspModel)
	Luma601Model  color.Model = color.ModelFunc(luma601Model)
	Luma240Model  color.Model = color.ModelFunc(luma240Model)
	Luma709Model  color.Model = color.ModelFunc(luma709Model)
	Luma2020Model color.Model = color.ModelFunc(luma2020Model)
)

func hslModel(c color.Color) color.Color {
	if _, ok := c.(HSL); ok {
		return c
	}
	return nil
}

type HSL struct {
}

func (c HSL) RGBA() (r, g, b, a uint32) {
	return
}

func hsvModel(c color.Color) color.Color {
	return nil
}

type HSV struct {
}

func (c HSV) RGBA() (r, g, b, a uint32) {
	return
}

func hspModel(c color.Color) color.Color {
	return nil
}

type HSP struct {
}

func (c HSP) RGBA() (r, g, b, a uint32) {
	return
}

func luma601Model(c color.Color) color.Color {
	return nil
}

type Luma601 struct {
}

func (c Luma601) RGBA() (r, g, b, a uint32) {
	return
}

func luma240Model(c color.Color) color.Color {
	return nil
}

type Luma240 struct {
}

func (c Luma240) RGBA() (r, g, b, a uint32) {
	return
}

func luma709Model(c color.Color) color.Color {
	return nil
}

type Luma709 struct {
}

func (c Luma709) RGBA() (r, g, b, a uint32) {
	return
}

func luma2020Model(c color.Color) color.Color {
	return nil
}

type Luma2020 struct {
}

func (c Luma2020) RGBA() (r, g, b, a uint32) {
	return
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

	// ONLY TO BE USED IN THIS FILE
	_r   mat.Dense
	_g   mat.Dense
	_b   mat.Dense
	_a   mat.Dense
	_gra mat.Dense // Average gray
	_lum mat.Dense // luminance
}

func (im *imageMatrix) Dims() (r, c int) {
	r, c, _, _ = _size(im, &im._r)
	return
}

func newImageMatrixFromImage(given image.Image) *imageMatrix {
	im := new(imageMatrix)
	im.Image = given
	return im
}

func newImageMatrixFromFile(filePath string) *imageMatrix {
	log.Printf("\"%s\"", filepath.Ext(strings.TrimSpace(filePath)))
	return _newImageMatrixFromFile(filePath, map[string]decodeFunc{
		".jpg": jpeg.Decode,
		".png": png.Decode,
	}[filepath.Ext(strings.TrimSpace(filePath))])
}

func _newImageMatrixFromFile(filePath string, decoder decodeFunc) *imageMatrix {
	rawData, err := os.Open(filePath)
	checkFatal(err)

	decodedImg, err := decoder(bufio.NewReader(rawData))
	checkFatal(err)
	log.Printf("Color Model: %s", decodedImg.ColorModel())
	min := decodedImg.Bounds().Min
	max := decodedImg.Bounds().Max
	r := max.X - min.X
	c := max.Y - min.Y
	return &imageMatrix{decodedImg,
		*mat.NewDense(r, c, nil),
		*mat.NewDense(r, c, nil),
		*mat.NewDense(r, c, nil),
		*mat.NewDense(r, c, nil),
		*mat.NewDense(r, c, nil),
		*mat.NewDense(r, c, nil)}
}

func _size(im *imageMatrix, dense *mat.Dense) (row, col, x, y int) {
	row, col = dense.Dims() // (row, col) = (x, y)
	min := im.Image.Bounds().Min
	max := im.Image.Bounds().Max
	x = max.X - min.X
	y = max.Y - min.Y
	return
}

func _sameSize(im *imageMatrix, dense *mat.Dense) bool {
	r, c, x, y := _size(im, dense)
	return r == x && c == y
}

// DATA POPULATING FUNCS FOR PACKAGE

func _validate(im *imageMatrix, given *mat.Dense, populateFunc func()) {
	if !_sameSize(im, given) {
		populateFunc()
		if !_sameSize(im, given) {
			panic(fmt.Sprintf(`imageMatrix: %+v\nDense: %+v\nNOT THE SAME SIZE`, im, given))
		}
	}
}

func (im *imageMatrix) red(given *mat.Dense) {
	_validate(im, given, im._populateRGBA)
	given = &im._r
}

func (im *imageMatrix) green(given *mat.Dense) {
	_validate(im, given, im._populateRGBA)
	given = &im._g
}

func (im *imageMatrix) blue(given *mat.Dense) {
	_validate(im, given, im._populateRGBA)
	given = &im._b
}

func (im *imageMatrix) alpha(given *mat.Dense) {
	_validate(im, given, im._populateRGBA)
	given = &im._a
}

func (im *imageMatrix) _populateRGBA() {
	rowlen, collen := im.Dims()
	for r := 0; r < rowlen; r++ {
		for c := 0; c < collen; c++ {
			red, green, blue, alpha := im.Image.At(r, c).RGBA()
			im._r.Set(r, c, float64(red))
			im._g.Set(r, c, float64(green))
			im._b.Set(r, c, float64(blue))
			im._a.Set(r, c, float64(alpha))
		}
	}
}

func (im *imageMatrix) gray(given *mat.Dense) { // Average gray
	_validate(im, given, im._populateGrayScale)
	given = &im._gra
}

func (im *imageMatrix) lumin() *mat.Dense { // luminance
	//_validate(im, given, im._populateGrayScale)
	return &im._lum
}

func (im *imageMatrix) _populateGrayScale() {
	rowlen, collen := im.Dims()
	for r := 0; r < rowlen; r++ {
		for c := 0; c < collen; c++ {
			red, green, blue, alpha := im.Image.At(r, c).RGBA()
			im._gra.Set(r, c, (float64(red)+float64(green)+float64(blue)+float64(alpha))/4.0)
			im._lum.Set(r, c, (float64(red)*0.2126+float64(green)*0.7152+float64(blue)*0.0722)*(float64(alpha)/256.0))
		}
	}
}

/*************************************************
*
* TESTING UTILITIES
*
***************************************************/

const ERRMARGIN float64 = 0.0000000000001

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
		// log.Printf("%s", line)
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
