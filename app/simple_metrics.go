package gaia

import (
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GuessParametersFromPath(path string)string {
	config := sdk.GetConfig()
	// It is not possible to perfectly select parameters from the path. guess and pick some.
	split := strings.Split(path, "/")
	for i := range split {
		frag := split[i]

		if len(frag) == 0 {
			continue
		}

		canBech32 := strings.Index(frag, "1")
		if canBech32 >= 0 {
			prefix := frag[0:canBech32]
			switch prefix {
			case config.GetBech32AccountAddrPrefix():
				split[i] = "{acc_address}"
				continue
			case config.GetBech32AccountPubPrefix():
				split[i] = "{acc_pub_address}"
				continue
			case config.GetBech32ConsensusAddrPrefix():
				split[i] = "{cons_address}"
				continue
			case config.GetBech32ConsensusPubPrefix():
				split[i] = "{cons_pub_address}"
				continue
			case config.GetBech32ValidatorAddrPrefix():
				split[i] = "{val_address}"
				continue
			case config.GetBech32ValidatorPubPrefix():
				split[i] = "{val_pub_address}"
				continue
			}
		}

		if _, err := strconv.ParseInt(frag,10,64); err == nil {
			split[i] = "{int_maybe_id}"
			continue
		}

		if len(frag) > 40 {
			// A general word can be treated as hex as long as there are only characters within the hex range.
			// In general, in the case of hex, there is a high probability that it is a hash value.
			// Therefore, it is expected that only values over 20 bytes will come (probably 32 bytes).
			if _, err := hex.DecodeString(frag); err == nil {
				split[i] = "{hex_maybe_hash}"
				continue
			}
		}
	}

	return strings.Join(split, "/")
}

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
