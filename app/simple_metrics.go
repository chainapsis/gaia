package gaia

import (
	"math"
	"time"
)

const ZScoreP90 = 1.29
const ZScoreP95 = 1.645
const ZScoreP99 = 2.33

type CalcAverageData struct {
	// Based on sec
	SumValues float64
	// Based on sec
	SumSquareValues float64

	NumItems int64
}

type AverageResponse struct {
	Average           float64
	Variance          float64
	StandardDeviation float64
	NumItems          int64

	P1  float64
	P5  float64
	P10 float64

	P90 float64
	P95 float64
	P99 float64
}

type SimpleMetrics struct {
	typeMap map[string]*CalcAverageData
}

func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{
		typeMap: make(map[string]*CalcAverageData),
	}
}

func (metric *SimpleMetrics) Measure(typ string, duration time.Duration) {
	_, ok := metric.typeMap[typ]
	if !ok {
		metric.typeMap[typ] = &CalcAverageData{
			SumValues:       0,
			SumSquareValues: 0,
			NumItems:        0,
		}
	}

	data := metric.typeMap[typ]
	data.NumItems++
	floatSec := float64(duration.Nanoseconds()) / 1_000_000_000
	data.SumValues += floatSec
	data.SumSquareValues += math.Pow(floatSec, 2)
}

func (metric *SimpleMetrics) CalcAverageResponse(typ string) AverageResponse {
	data, ok := metric.typeMap[typ]
	if !ok {
		return AverageResponse{}
	}

	sumValues := data.SumValues
	sumSquareValues := data.SumSquareValues
	numItems := float64(data.NumItems)

	average := sumValues / numItems
	squareDeviation := sumSquareValues + (math.Pow(average, 2) * numItems) + (-2 * average * sumValues)
	variance := squareDeviation / numItems
	standardDeviation := math.Sqrt(variance)
	return AverageResponse{
		Average:           average,
		Variance:          variance,
		StandardDeviation: standardDeviation,
		NumItems:          data.NumItems,

		P1:  (-ZScoreP99 * standardDeviation) + average,
		P5:  (-ZScoreP95 * standardDeviation) + average,
		P10: (-ZScoreP90 * standardDeviation) + average,

		P90: ZScoreP90*standardDeviation + average,
		P95: ZScoreP95*standardDeviation + average,
		P99: ZScoreP99*standardDeviation + average,
	}
}

func (metric *SimpleMetrics) CalcAllAverageResponses() map[string]AverageResponse {
	resp := make(map[string]AverageResponse)
	for typ := range metric.typeMap {
		resp[typ] = metric.CalcAverageResponse(typ)
	}
	return resp
}
