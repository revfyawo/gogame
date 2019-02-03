package components

import "github.com/ojrac/opensimplex-go"

type Noise struct {
	Noise opensimplex.Noise
}

func (n Noise) Eval(x, y, f float64) float64 {
	return n.Noise.Eval2(x*f, y*f)
}
