package gaia

import (
	"encoding/json"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *GaiaApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if req.Path == "/cosmos.staking.v1beta1.Query/ValidatorDelegations" {
		// To check the number of forbidden requests
		app.SimpleMetrics.Measure("/cosmos.staking.v1beta1.Query/ValidatorDelegations+forbidden", 0)
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	if req.Path == "custom/staking/validatorDelegations" || req.Path == "/custom/staking/validatorDelegations" {
		// To check the number of forbidden requests
		app.SimpleMetrics.Measure("custom/staking/validatorDelegations+forbidden", 0)
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	if req.Path == "custom/staking/validatorUnbondingDelegations" || req.Path == "/custom/staking/validatorUnbondingDelegations" {
		// To check the number of forbidden requests
		app.SimpleMetrics.Measure("custom/staking/validatorUnbondingDelegations+forbidden", 0)
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	if req.Path == "/cosmos.tx.v1beta1.Service/GetTxsEvent" {
		// To check the number of forbidden requests
		app.SimpleMetrics.Measure("/cosmos.tx.v1beta1.Service/GetTxsEvent+forbidden", 0)
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	if req.Path == "/cosmos.gov.v1beta1.Query/TallyResult" {
		// To check the number of forbidden requests
		app.SimpleMetrics.Measure("/cosmos.gov.v1beta1.Query/TallyResult+forbidden", 0)
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

	cacheKey := GetCacheKey(req)
	if len(cacheKey) > 0 {
		cached, found := GetCachedValue(cacheKey)
		if found {
			// Check the count of the cached response
			app.SimpleMetrics.Measure(req.Path+"+cached", 0)
			return cached
		}
	}

	start := time.Now()
	res = app.BaseApp.Query(req)
	elapsed := time.Since(start)

	app.SimpleMetrics.Measure(req.Path, elapsed)

	if len(cacheKey) > 0 {
		SetCache(req.Path, cacheKey, res)
	}

	return res
}
