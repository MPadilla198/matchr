package img

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

// TODO Change MetricFunc to be func(img1, img2 *imageMatrix) []float64
type MetricFunc func(file1, file2 *imageMatrix) []float64

var ImageMetrics map[string]MetricFunc

func Map(v []float64, f func(float64) float64) (vm []float64) {
	vm = make([]float64, len(v))
	for i, vi := range v {
		vm[i] = f(vi)
	}
	return
}

func init() {
	ImageMetrics = map[string]MetricFunc{"mse": Mse, "rmse": Rmse, "sam": Sam, "rase": Rase, "ergas": Ergas, "uqi": Uqi,
		"ssim": Ssim, "psnr": Psnr, "msssimi": Msssimi, "vif": Vif, "dlambda": Dlambda, "ds": Ds, "qnr": Qnr}
}

func _validateInput(img1, img2 *imageMatrix) (ok bool) {
	r1, c1 := img1.Dims()
	r2, c2 := img2.Dims()
	if r1 == r2 && c1 == c2 {
		ok = true
	}
	return
}

func Mse(img1, img2 *imageMatrix) (result []float64) {
	if ok := _validateInput(img1, img2); !ok {
		panic("Image validation Failed.")
	}

	r, c := img1.Dims()
	resultMat := mat.NewDense(r, c, make([]float64, r*c))

	mat.NewDense(r, c, make([]float64, r*c)).Apply(func(i, j int, v float64) float64 {
		return math.Pow(math.Abs(v-img2.At(i, j)), 2)
	}, img1)

	result = []float64{mat.Sum(resultMat) / float64(r) * float64(c)}
	return
}

func Rmse(img1, img2 *imageMatrix) []float64 {
	return Map(Mse(img1, img2), math.Sqrt)
}

func Sam(img1, img2 *imageMatrix) []float64 {
	if ok := _validateInput(img1, img2); !ok {
		panic("Image validation Failed.")
	}
	r, c := img1.Dims()
	result := mat.NewDense(r, c, make([]float64, r*c))

	result.Apply(func(i, j int, _ float64) float64 {
		v1 := img1.ColorAt(i, j, true)
		v2 := img2.ColorAt(i, j, true)
		return math.Acos(mat.Dot(v1, v2) / math.Sqrt(mat.Dot(v1, v1)*mat.Dot(v2, v2)))
	}, img1)

	return []float64{mat.Sum(result) / (float64(r) * float64(c))}
}

func Rase(img1, img2 *imageMatrix) []float64 {
	if ok := _validateInput(img1, img2); !ok {
		panic("Image validation Failed.")
	}
	r, c := img1.Dims()
	mr := (MeanRadiance(img1) + MeanRadiance(img2)) / 2.0
	size := float64(r) * float64(c)
	result := mat.NewDense(r, c, make([]float64, r*c))

	result.Apply(func(i, j int, _ float64) float64 {
		v1 := img1.ColorAt(i, j, false)
		v2 := img2.ColorAt(i, j, false)
		return math.Acos(mat.Dot(v1, v2) / math.Sqrt(mat.Dot(v1, v1)*mat.Dot(v2, v2)))
	}, img1)

	return []float64{(100.00 / mr) * math.Sqrt((1/size)*mat.Sum(result))}
}

func MeanRadiance(img *imageMatrix) float64 {
	r, c := img.Dims()
	return mat.Sum(img) / (float64(r) * float64(c))
}

func Ergas(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Uqi(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Ssim(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Psnr(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Msssimi(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Vif(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Dlambda(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Ds(file1, file2 *imageMatrix) []float64 {
	return nil
}

func Qnr(file1, file2 *imageMatrix) []float64 {
	return nil
}
