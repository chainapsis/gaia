package gaia_test

import (
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

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

func TestGuessParametersFromPath(t *testing.T) {
	hash := crypto.Sha256([]byte{1,2,3})
	hexStr := hex.EncodeToString(hash)
	path := "/test/tt/cosmos1test/cosmospub1/invalidbech1test/invalidbech1/123/aaa/" + hexStr

	res := app.GuessParametersFromPath(path)

	require.Equal(t, "/test/tt/{acc_address}/{acc_pub_address}/invalidbech1test/invalidbech1/{int_maybe_id}/aaa/{hex_maybe_hash}", res)
}
