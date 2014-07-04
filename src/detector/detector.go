package main

import "log"

import "io/ioutil"
import "encoding/json"

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

    log.Println("Training SAX:", TrainingData.SAX(DimensionWindow, AlphabetSize))
    log.Println("Testing SAX:", TestingData.SAX(DimensionWindow, AlphabetSize))
}
