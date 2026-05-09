package service

import (
	"image/color"
	"log"
	"math"
	rand "math/rand/v2"
	"stonkx/internal/domain"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const (
	quotsPerSeriesGenFunc = 86391
	quotsInitialAmount    = 10
)

type series struct {
	values []float64
	alpha  float64
	pheta  float64
	sigma  float64
	mean   float64
	EMA    float64
}

func GenerateDefault() domain.SeriesEx {
	return Generate(quotsPerSeriesGenFunc, rand.ExpFloat64(), rand.Float64())
}

func Generate(size int, a, p float64) domain.SeriesEx {
	/*
		size - number of ticks generated
		0 < a < 1 - smoothing coefficeint: the more the a, the smoother the series is.
		0 < p - speed of mean-reversing
		s - standard deviation: rate of volatility (static for now)
	*/

	var ser series = series{}
	ser = initSeries(quotsInitialAmount)
	ser.alpha = a
	ser.pheta = p
	ser.sigma = math.Log(ser.mean)

	for range size {
		ser.generateNew()
	}

	ser.normalize()
	exSer := convertToEx(ser)
	return exSer
}

func initSeries(size int) series {
	res := make([]float64, size)
	start := rand.ExpFloat64()
	res[0] = start
	for i := 1; i < size; i++ {
		res[i] = res[i-1] * rand.NormFloat64()
	}

	mean := func(ser []float64) float64 {
		sum := 0.0
		for _, v := range ser {
			sum += v
		}
		return sum / float64(len(ser))
	}

	ser := series{
		values: res,
		mean:   mean(res),
		EMA:    mean(res),
	}

	ser.normalize()
	ser.mean = mean(ser.values)
	ser.EMA = ser.mean

	return ser
}

func (s *series) normalize() {
	min := func(ser []float64) float64 {
		res := ser[0]
		for _, v := range ser {
			if v < res {
				res = v
			}
		}
		return res
	}

	add := func(ser []float64, num float64) []float64 {
		res := make([]float64, len(ser))
		for i := range ser {
			res[i] = ser[i] + num
		}
		return res
	}

	if m := min(s.values); m < 0 {
		res := add(s.values, -1.01*m)
		copy(s.values, res)
	}

}

func (s *series) generateNew() {
	//f(t+1) = f(t) + θ*(EMA - f(t)) + σ*Z
	//EMA(t) = α*f(t) + (1-α)*EMA(t-1), where:
	//    α - smoothing coefficient (0 < α < 1)

	ph := s.pheta
	ema := s.EMA
	sig := s.sigma
	a := s.alpha
	lastValue := s.values[len(s.values)-1]

	newEma := a*lastValue + (1-a)*ema
	new := lastValue + ph*(newEma-lastValue) + sig*rand.NormFloat64()
	s.values = append(s.values, new)

	s.EMA = newEma
}

func (s series) plot(xSize, ySize float64, path string) {

	val := s.values
	xys := func(n int) plotter.XYs {
		pts := make(plotter.XYs, n)
		for i := range pts {

			pts[i].X = float64(i)
			pts[i].Y = val[i]
		}
		return pts
	}

	n := len(val)
	data := xys(n)

	p := plot.New()
	p.Title.Text = "Time Series"

	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = color.RGBA{G: 255, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, points)

	err = p.Save(vg.Length(xSize)*vg.Centimeter, vg.Length(ySize)*vg.Centimeter, path)
	if err != nil {
		log.Panic(err)
	}
}

func convertToEx(s series) domain.SeriesEx {
	res := domain.SeriesEx{}
	val := s.values

	convertToInt := func(num float64) int {
		return int(math.Round(num * 100.0))
	}

	candleCoefficient := func(open, close float64) float64 {
		if open > close {
			close, open = open, close
		}

		return rand.ExpFloat64() * (close / open) / 100
	}

	for i := 0; i < len(val)-1; i++ {
		q := domain.QuotEx{}

		q.Time = i
		o, c := val[i], val[i+1]

		if o > c {
			q.Low = convertToInt(c * (1 - candleCoefficient(o, c)))
			q.High = convertToInt(o * (1 + candleCoefficient(o, c)))
		} else {
			q.Low = convertToInt(o * (1 - candleCoefficient(o, c)))
			q.High = convertToInt(c * (1 + candleCoefficient(o, c)))
		}
		q.Open = convertToInt(o)
		q.Close = convertToInt(c)

		res.Quots = append(res.Quots, q)

	}

	return res
}
