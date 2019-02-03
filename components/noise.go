package components

import "github.com/ojrac/opensimplex-go"

type Noise struct {
	Noise opensimplex.Noise
}

func (n Noise) Eval(x, y, f float64) float64 {
	return n.Noise.Eval2(x*f, y*f)
}

func (n Noise) EvalOctaves(x, y, f float64, octaves uint) float64 {
	noise := 0.
	var one uint = 1
	var i uint
	for i = 0; i < octaves; i++ {
		amp := float64(one << (octaves - i))
		freq := float64(one<<i) * f
		noise += n.Noise.Eval2(x*freq, y*freq) * amp
	}
	noise /= float64(one << octaves)
	return noise
}
