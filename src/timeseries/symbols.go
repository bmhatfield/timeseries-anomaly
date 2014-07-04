package timeseries

var Alphabet = "abcdefghijklkmnop"

var Breakpoints = map[int][]float64{
    3: []float64{-0.43, 0.43},
    4: []float64{-0.67, 0, 0.67},
    5: []float64{-0.84, -0.25, 0.25, 0.84},
    6: []float64{-0.97, -0.43, 0, 0.43, 0.97},
    7: []float64{-1.07, -0.57, -0.18, 0.18, 0.57, 1.07},
}
