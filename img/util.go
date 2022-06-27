package img

import (
  "image"
  "gonum.org/v1/gonum/mat"
)

type image_matrix struct {
  r mat.Matrix
  g mat.Matrix
  b mat.Matrix
  a mat.Matrix
}

// matricize_image returns an image_matrix, corresponding to the r, g, b, and a
// color channels of the image.
func matricize_image(img image.Image) image_matrix {
  bounds := img.Bounds()
  x := bounds.Max.X-bounds.Min.X
  y := bounds.Max.Y-bounds.Min.Y
  size := x+y

  r := make([]float64, size, size)
  g := make([]float64, size, size)
  b := make([]float64, size, size)
  a := make([]float64, size, size)

  for j := 0; j < y; j++ {
    for i := 0; i < x; i++ {
      curr_index := (j * x) + i
      cr, cg, cb, ca := img.At(i+bounds.Min.X, j+bounds.Min.Y).RGBA()

      r[curr_index] = float64(cr)
      g[curr_index] = float64(cg)
      b[curr_index] = float64(cb)
      a[curr_index] = float64(ca)
    }
  }

  return image_matrix{r: mat.NewDense(x,y,r), g: mat.NewDense(x,y,g), b: mat.NewDense(x,y,b), a: mat.NewDense(x,y,a),}
}
