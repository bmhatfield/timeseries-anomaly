package anomaly

//import "log"
import "math"

import "timeseries"

type ExtractedTimeSeries struct {
    TimeSeries timeseries.TimeSeries

    FeatureLength int
    WordLength int
    AlphabetSize int

    Words []string
    WordCounts map[string]int
}

func New(timeseries timeseries.TimeSeries) (*ExtractedTimeSeries) {
    es := &ExtractedTimeSeries{TimeSeries: timeseries}
    es.Words = make([]string, 0)
    es.WordCounts = make(map[string]int)
    return es
}

func (e *ExtractedTimeSeries) ExtractWords() {
    // Iterate over the raw timeseries, extracting normalized words and appending
    // them to the array.
    dimensionLength := e.FeatureLength * e.WordLength

    for i := 0; i < e.TimeSeries.Length() / dimensionLength; i++ {
        offset := i * dimensionLength

        ts := timeseries.TimeSeries(e.TimeSeries[offset:offset+dimensionLength])

        word := ts.SAX(e.FeatureLength, e.AlphabetSize)

        e.Words = append(e.Words, word)
        e.WordCounts[word] += 1
    }
}

func (e *ExtractedTimeSeries) CompareTo(testing *ExtractedTimeSeries) (surprise []float64) {
    wordDataPoints := e.FeatureLength * e.WordLength
    scaleFactor := float64(len(testing.TimeSeries) - wordDataPoints + 1) / float64(len(e.TimeSeries) - wordDataPoints + 1)

    for _, word := range testing.Words {
        variance := math.Abs(float64(testing.WordCounts[word]) - (scaleFactor * float64(e.WordCounts[word])))
        surprise = append(surprise, variance)
    }

    return surprise
}
