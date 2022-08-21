package img

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
)

type MetricFunc func(file1, file2 string) ([]float64, error)

var ImageMetrics map[string]MetricFunc

func init() {
	ImageMetrics = map[string]MetricFunc{"mse": Mse, "rmse": Rmse, "sam": Sam, "race": Race, "ergas": Ergas, "uqi": Uqi,
		"ssim": Ssim, "psnr": Psnr, "msssimi": Msssimi, "vif": Vif, "dlambda": Dlambda, "ds": Ds, "qnr": Qnr}
}

func Mse(file1, file2 string) ([]float64, error) {
	log.Printf("%s and %s", file1, file2)
	img1 := newImageMatrixFromFile(file1).lumin()
	img2 := newImageMatrixFromFile(file2).lumin()

	r, c := img1.Dims()
	result := mat.NewDense(r, c, nil)

	result.Apply(func(i, j int, v float64) float64 {
		return math.Pow(math.Abs(v-img2.At(i, j)), 2)
	}, img1)

	return []float64{mat.Sum(result)}, nil
}

func Rmse(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Sam(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Race(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Ergas(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Uqi(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Ssim(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Psnr(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Msssimi(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Vif(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Dlambda(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Ds(file1, file2 string) ([]float64, error) {
	return nil, nil
}

func Qnr(file1, file2 string) ([]float64, error) {
	return nil, nil
}
