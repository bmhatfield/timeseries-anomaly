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

func loadJSON(filename string, destination *[]float64) error {
    file, err := ioutil.ReadFile(filename)

    if err != nil {
        fmt.Println("Error reading file:", filename, err)
        return err
    }

    err = json.Unmarshal(file, destination)

    if err != nil {
        fmt.Println("Unable to parse JSON:", err)
        return err
    }

    return err
}

func preProcess(series string, window int) map[string]int {
    words := make(map[string]int)

    for offset := 0; offset < len(series) - window; offset++ {
        word := series[offset:offset+window]
        words[word]++
    }

    return words
}

func compare(series string, window int, training, testing map[string]int) []float64 {
    deltas := make([]float64, len(series) - window)

    for offset := 0; offset < len(series) - window; offset++ {
        word := series[offset:offset+window]

        deltas[offset] = math.Abs(float64(testing[word]) - float64(training[word]))
    }

    return deltas
}

var trainingData []float64
var anomalyData []float64

func main() {
    series := make(map[string][]float64)

    var err error
    err = loadJSON("trainingdata.json", &trainingData)

    if err != nil {
        return
    }

    err = loadJSON("anomalydata.json", &anomalyData)

    if err != nil {
        return
    }

    tznWindow := 8
    paaWindow := 30
    alphabetLength := 4

    series["training"] = trainingData
    series["trainingpaa"] = paa(trainingData, paaWindow)
    series["trainingnormal"] = normalize(trainingData)
    series["trainingnpaa"] = paa(normalize(trainingData), paaWindow)

    series["anomaly"] = anomalyData
    series["anomalypaa"] = paa(anomalyData, paaWindow)
    series["anomalynormal"] = normalize(anomalyData)
    series["anomalynpaa"] = paa(normalize(anomalyData), paaWindow)

    series["breakpoints"] = breakpoints[alphabetLength]

    tr := string(sax(series["trainingnpaa"], alphabetLength))
    tx := string(sax(series["anomalynpaa"], alphabetLength))

    trwords := preProcess(tr, tznWindow)
    txwords := preProcess(tx, tznWindow)

    series["deltas"] = compare(tx, tznWindow, trwords, txwords)

    fmt.Println("Deltas:", series["deltas"])
    fmt.Println("Tr: ", trwords)
    fmt.Println("Tx: ", txwords)

    d, e := json.Marshal(series)

    if e == nil {
        contents := []byte("var series = ")
        contents = append(contents, d...)
        contents = append(contents, []byte(";\n")...)
        contents = append(contents, []byte(fmt.Sprintf("var paaWindow = %v;\n", paaWindow))...)
        contents = append(contents, []byte(fmt.Sprintf("var trainingLen = %v;\n", len(trainingData)))...)
        contents = append(contents, []byte(fmt.Sprintf("var anomalyLen = %v;\n", len(anomalyData)))...)

        ioutil.WriteFile("tsdata.js", contents, 0644)
    }
}
