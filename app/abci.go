package gaia

import (
	"encoding/json"
	"strings"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *GaiaApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if req.Path == "/cosmos.staking.v1beta1.Query/ValidatorDelegations" {
		// To check the number of forbidden requests
		app.SimpleMetrics.Measure("/cosmos.staking.v1beta1.Query/ValidatorDelegations", 0)
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	if req.Path == "/simple-metric" || req.Path == "simple-metric" {
		metricRes := app.SimpleMetrics.CalcAllAverageResponses()
		jsonRes, err := json.Marshal(metricRes)
		if err != nil {
			panic(err)
		}
		return abci.ResponseQuery{
			Value: jsonRes,
		}
	}

	// Only measure the grpc request.
	// Custom request is hard to measure because it has the parameter datas in the path.
	if strings.HasPrefix(req.Path, "custom") || strings.HasPrefix(req.Path, "/custom") {
		return app.BaseApp.Query(req)
	}

	start := time.Now()
	res = app.BaseApp.Query(req)
	elapsed := time.Since(start)

	app.SimpleMetrics.Measure(req.Path, elapsed)

	return res
}
