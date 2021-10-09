package gaia_test

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	app "github.com/cosmos/gaia/v5/app"
)

func TestSimpleMetric(t *testing.T) {
	metric := app.NewSimpleMetrics()

	metric.Measure("test", time.Second)
	metric.Measure("test", time.Second*3)
	metric.Measure("test", time.Second*5)
	metric.Measure("test", time.Second*3)
	metric.Measure("test", time.Second)

	res := metric.CalcAverageResponse("test")

	require.Equal(t, 2.6, res.Average)
	require.True(t, math.Abs(2.24-res.Variance) < 0.00001)
	require.True(t, math.Abs(1.49666295471-res.StandardDeviation) < 0.00001)
	require.Equal(t, int64(5), res.NumItems)
}
