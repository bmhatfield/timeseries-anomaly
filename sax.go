package main

import "fmt"
import "math"
import "io/ioutil"
import "encoding/json"

var f = 0.25
var a = 1.0

var breakpoints = map[int][]float64{
    3: []float64{-0.43, 0.43},
    4: []float64{-0.67, 0, 0.67},
    5: []float64{-0.84, -0.25, 0.25, 0.84},
    6: []float64{-0.97, -0.43, 0, 0.43, 0.97},
    7: []float64{-1.07, -0.57, -0.18, 0.18, 0.57, 1.07},
}

var alphabet = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

func sum(points []float64) (s float64) {
    for _, p := range points {
        s += p
    }

    return
}

func mean(points []float64) (float64) {
    return sum(points) / float64(len(points))
}

func stdev(points []float64) (float64) {
    m := mean(points)

    variances := make([]float64, len(points))

    for i, p := range points {
        variances[i] = math.Pow(p - m, 2)
    }

    return math.Sqrt(mean(variances))
}

func normalize(points []float64) ([]float64) {
    m := mean(points)
    s := stdev(points)

    normalized := make([]float64, len(points))

    for i, p := range points {
        normalized[i] = (p - m) / s
    }

    return normalized
}

func paa(points []float64, window int) (ps []float64) {
    ps = make([]float64, len(points)/window)

    for w := 0; w < len(points) / window; w++ {
        offset := w * window

        ps[w] = mean(points[offset:offset + window - 1])
    }

    return ps
}

func sw(time, frequency, amplitude float64) float64 {
    return amplitude * math.Sin(frequency * time)
}

func sax(paaSeries []float64, alphaLen int) []byte {
    w := make([]byte, len(paaSeries))

    for i, s := range paaSeries {
        for b, breakpoint := range breakpoints[alphaLen] {
            if s <= breakpoint {
                w[i] = alphabet[b]
                break
            } else if s > breakpoints[alphaLen][alphaLen - 2] {
                w[i] = alphabet[alphaLen - 1]
                break
            }
        }
    }

    return w
}

func main() {
    paaWindow := 5
    alphabetLength := 4
    rawLength := 100

    series := make(map[string][]float64)

    data := make([]float64, rawLength)

    for t := 0; t < rawLength; t++ {
        data[t] = sw(float64(t), f, a)
    }

    series["sin"] = data
    series["paa"] = paa(data, 5)
    series["normal"] = normalize(data)
    series["npaa"] = paa(normalize(data), paaWindow)
    series["breakpoints"] = breakpoints[alphabetLength]

    fmt.Println("word: ", string(sax(series["npaa"], alphabetLength)))

    d, e := json.Marshal(series)

    if e == nil {
        contents := []byte("var series = ")
        contents = append(contents, d...)
        contents = append(contents, []byte(";\n")...)
        contents = append(contents, []byte(fmt.Sprintf("var paaWindow = %v;\n", paaWindow))...)
        contents = append(contents, []byte(fmt.Sprintf("var rawLength = %v;\n", rawLength))...)

        ioutil.WriteFile("tsdata.js", contents, 0644)
    }
}
