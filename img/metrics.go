package img

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

// TODO Change MetricFunc to be func(img1, img2 *imageMatrix) []float64
type MetricFunc func(file1, file2 string) []float64

var ImageMetrics map[string]MetricFunc

func Map(v []float64, f func(float64) float64) (vm []float64) {
	vm = make([]float64, len(v))
	for i, vi := range v {
		vm[i] = f(vi)
	}
	return
}

func init() {
	ImageMetrics = map[string]MetricFunc{"mse": Mse, "rmse": Rmse, "sam": Sam, "race": Race, "ergas": Ergas, "uqi": Uqi,
		"ssim": Ssim, "psnr": Psnr, "msssimi": Msssimi, "vif": Vif, "dlambda": Dlambda, "ds": Ds, "qnr": Qnr}
}

func Mse(file1, file2 string) []float64 {
	// log.Printf("%s and %s", file1, file2)
	img1 := newImageMatrixFromFile(file1, Luma709Model)
	img2 := newImageMatrixFromFile(file2, Luma709Model)

	r, c := img1.Dims()
	size := float64(r) * float64(c)
	result := mat.NewDense(r, c, make([]float64, r*c))

	result.Apply(func(i, j int, v float64) float64 {
		return math.Pow(math.Abs(v-img2.At(i, j)), 2)
	}, img1)

	return []float64{mat.Sum(result) / size}
}

func Rmse(file1, file2 string) []float64 {
	return Map(Mse(file1, file2), math.Sqrt)
}

func Sam(file1, file2 string) []float64 {
	return nil
}

func Race(file1, file2 string) []float64 {
	return nil
}

func Ergas(file1, file2 string) []float64 {
	return nil
}

func Uqi(file1, file2 string) []float64 {
	return nil
}

func Ssim(file1, file2 string) []float64 {
	return nil
}

func Psnr(file1, file2 string) []float64 {
	return nil
}

func Msssimi(file1, file2 string) []float64 {
	return nil
}

func Vif(file1, file2 string) []float64 {
	return nil
}

func Dlambda(file1, file2 string) []float64 {
	return nil
}

func Ds(file1, file2 string) []float64 {
	return nil
}

func Qnr(file1, file2 string) []float64 {
	return nil
}
