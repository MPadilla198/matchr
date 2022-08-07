package img

type ImgMetric func(file1, file2 string) ([]float64, error)

var ImageMetrics map[string]ImgMetric

func init() {
	ImageMetrics = map[string]ImgMetric{"mse": Mse, "rmse": Rmse, "sam": Sam, "race": Race, "ergas": Ergas, "uqi": Uqi,
		"ssim": Ssim, "psnr": Psnr, "msssimi": Msssimi, "vif": Vif, "dlambda": Dlambda, "ds": Ds, "qnr": Qnr}
}

func Mse(file1, file2 string) ([]float64, error) {
	return nil, nil
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
