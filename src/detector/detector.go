package main

import "fmt"
import "log"

import "io/ioutil"
import "encoding/json"

import "anomaly"
import "timeseries"

func loadJSON(filename string, destination *timeseries.TimeSeries) {
    file, err := ioutil.ReadFile(filename)

    if err != nil {
        log.Fatalf("Error reading file '%s': %v", filename, err)
    }

    err = json.Unmarshal(file, destination)

    if err != nil {
        log.Fatalf("Unable to parse JSON: %v", err)
    }
}

func main() {
    var TrainingData timeseries.TimeSeries
    loadJSON("trainingdata.json", &TrainingData)

    var TestingData timeseries.TimeSeries
    loadJSON("anomalydata.json", &TestingData)

    DimensionWindow := 30
    AlphabetSize := 4
    WordLength := 4

    extractedTrainingData := anomaly.New(TrainingData)
    extractedTrainingData.FeatureLength = DimensionWindow
    extractedTrainingData.WordLength = WordLength
    extractedTrainingData.AlphabetSize = AlphabetSize
    extractedTrainingData.ExtractWords()

    extractedTestingData := anomaly.New(TestingData)
    extractedTestingData.FeatureLength = DimensionWindow
    extractedTestingData.WordLength = WordLength
    extractedTestingData.AlphabetSize = AlphabetSize
    extractedTestingData.ExtractWords()

    data := make(map[string]interface{})

    data["paaWindow"] = DimensionWindow
    data["wordLength"] = WordLength

    data["trainingnpaa"] = TrainingData.Normalized().PAA(DimensionWindow)
    data["anomalynpaa"] = TestingData.Normalized().PAA(DimensionWindow)
    data["surprise"] = extractedTrainingData.CompareTo(extractedTestingData)

    seriesRawJson, e := json.MarshalIndent(data, "", "  ")

    if e == nil {
        ioutil.WriteFile("tsdata.js", []byte(fmt.Sprintf("var series = %s;\n", seriesRawJson)), 0644)
    }
}
