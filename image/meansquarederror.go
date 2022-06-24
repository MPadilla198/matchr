package image

import (
  im "image"
  "math"
)

// Mean Squared Average
func MSE(i1, i2 *im.Image) float64 {
  ib1 := i1.Bounds()
  ib2 := i2.Bounds()

  xbound := ib1.Max.X - ib1.Min.X
  ybound := ib1.Max.Y - ib1.Min.Y

  if xbound != ib2.Max.X - ib2.Min.X || ybound != ib2.Max.Y - ib2.Min.Y {
    // ERROR
  }

  var sum float64 = 0.0

  for dx := range xbound {
    for dy := range ybound {
      r1, g1, b1, a1 := i1.At(dx+ib1.Min.X, dy+ib1.Min.Y).RGBA()
      r2, g2, b2, a2 := i2.At(dx+ib2.Min.X, dy+ib2.Min.Y).RGBA()

      sum += (math.Pow(r2-r1, 2) + math.Pow(g2-g1, 2) + math.Pow(b2-b1, 2) + math.Pow(a2-a1, 2)) / 4.0
    }
  }

  return sum / (xbound*ybound)
}
