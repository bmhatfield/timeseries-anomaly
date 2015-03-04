package timeseries

import "log"
import "math"

// Custom TimeSeries type wrapping the []float64 type.
type TimeSeries []float64

// Convenience method that dereferences the series pointer and
// calls len() on the object.
func (ts *TimeSeries) Length() (int) {
    return len(*ts)
}

// Sums the time series values
func (ts *TimeSeries) Sum() (sum float64) {
    for _, value :=  range *ts {
        sum += value
    }

    return sum
}

// Averages the timeseries values
func (ts *TimeSeries) Mean() (float64) {
    return ts.Sum() / float64(ts.Length())
}

// Computes the Standard Deviation of the timeseries values
func (ts *TimeSeries) Stdev() (float64) {
    mean := ts.Mean() // cache series mean

    deviations := make(TimeSeries, ts.Length())

    for index, value := range *ts {
        deviations[index] = math.Pow(value - mean, 2)
    }

    return math.Sqrt(deviations.Mean())
}

// Computes a z-normal (standard score) version of the timeseries
func (ts *TimeSeries) Normalized() (*TimeSeries) {
    mean := ts.Mean() // cache series mean
    stdev := ts.Stdev() // cache series standard deviation

    znormal := make(TimeSeries, ts.Length())

    for index, value := range *ts {
        znormal[index] = (value - mean) / stdev
    }

    return &znormal
}

// Computes a peicewise-aggregate-approximation of the timeseries.
// This method returns a dimension-reduced version of the timeseries of
// the order len(TimeSeries) / window. Each value in the returned series
// is representative of n values (n=window).
func (ts *TimeSeries) PAA(window int) (*TimeSeries) {
    paa := make(TimeSeries, ts.Length() / window) // dimension-reduced container

    series := *ts // dereference pointer to allow for slicing

    for i := 0; i < ts.Length() / window; i++ {
        offset := i * window

        subseq := TimeSeries(series[offset:offset + window])

        paa[i] = subseq.Mean()
    }

    return &paa
}

// Computes a dimensionality-reduced symbolic approximation of the timeseries.
// Dimensionality reduction is performed by PAA.
func (ts *TimeSeries) SAX(window, alphabetSize int) (sax string) {
    topBreakpoint := alphabetSize - 2
    breakpoints := Breakpoints[alphabetSize]

    normalPAA := *ts.Normalized().PAA(window)

    for _, value := range normalPAA {
        if value > breakpoints[topBreakpoint] {
            sax += string(Alphabet[alphabetSize - 1])
            continue
        } else {
            for index, breakpoint := range breakpoints {
                if value <= breakpoint {
                    sax += string(Alphabet[index])
                    break
                }
            }
        }
    }

    if len(sax) != len(normalPAA) {
        // This should never be the case; symbolized logic should always lead to
        // a symbolized string with a length matching the number of values to symbolize.
        log.Fatalf("Severe Symbolication Issue: SAX len of %v != PAA len of %v", len(sax), len(normalPAA))
    }

    return sax
}
