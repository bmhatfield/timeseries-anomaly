package main

import "fmt"
import "math"
import "io/ioutil"
import "encoding/json"

import "timeseries"

func sineWave(time, frequency, amplitude float64) float64 {
    return amplitude * math.Sin(frequency * time)
}

func loadJSON(filename string, destination *timeseries.TimeSeries) error {
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

func compare(series string, window int, training, testing map[string]int, scale float64) []float64 {
    deltas := make([]float64, len(series) - window)

    for offset := 0; offset < len(series) - window; offset++ {
        word := series[offset:offset+window]

        delta := math.Abs(float64(testing[word]) - (scale * float64(training[word])))

        fmt.Printf("%s found %v times in 'testing' data, %v times in 'training' data (%v scaled %0.4fx): %v\n",
            word, testing[word], (scale * float64(training[word])), training[word], scale, delta)

        deltas[offset] = delta
    }

    return deltas
}

var trainingData timeseries.TimeSeries
var anomalyData timeseries.TimeSeries

func main() {
    series := make(map[string][]float64)
    tunables := make(map[string]float64)

    var err error
    err = loadJSON("trainingdata.json", &trainingData)

    if err != nil {
        return
    }

    err = loadJSON("anomalydata.json", &anomalyData)

    if err != nil {
        return
    }

    tunables["tznWindow"] = 6
    tunables["paaWindow"] = 30
    tunables["alphabetLength"] = 4

    series["training"] = trainingData
    series["trainingpaa"] = trainingData.Paa(tunables["paaWindow"])
    series["trainingnormal"] = *trainingData.Normalize()
    series["trainingnpaa"] = trainingData.Normalize().Paa(tunables["paaWindow"])

    series["anomaly"] = anomalyData
    series["anomalypaa"] = anomalyData.Paa(tunables["paaWindow"])
    series["anomalynormal"] = *anomalyData.Normalize()
    series["anomalynpaa"] = anomalyData.Normalize().Paa(tunables["paaWindow"])

    series["breakpoints"] = timeseries.Breakpoints[tunables["alphabetLength"]]

    tr := string(trainingData.Sax(tunables["paaWindow"], tunables["alphabetLength"]))
    tx := string(anomalyData.Sax(tunables["paaWindow"], tunables["alphabetLength"]))

    trwords := preProcess(tr, tunables["tznWindow"])
    txwords := preProcess(tx, tunables["tznWindow"])

    scale := float64(len(tx) - tunables["tznWindow"] + 1) / float64(len(tr)  - tunables["tznWindow"] + 1)
    series["deltas"] = compare(tx, tunables["tznWindow"], trwords, txwords, scale)

    seriesRawJson, e := json.MarshalIndent(series, "", "  ")
    tunablesRawJson, e := json.MarshalIndent(tunables, "", "  ")

    if e == nil {
        contents := []byte(fmt.Sprintf("var series = %s;\n", seriesRawJson))
        contents = append(contents, []byte(fmt.Sprintf("var tunables = %s;\n", tunablesRawJson))...)

        ioutil.WriteFile("tsdata.js", contents, 0644)
    }
}
