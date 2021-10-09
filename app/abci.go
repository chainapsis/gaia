package gaia

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/patrickmn/go-cache"
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

	if req.Path == "custom/gov/tally" || req.Path == "/custom/gov/tally" {
		bz, err := sdk.SortJSON(req.Data)
		if err == nil {
			key := hex.EncodeToString(bz)
			cached, found := queryGovTallyCache.Get(key)
			if cached != nil && found {
				res, ok := cached.(abci.ResponseQuery)
				if ok {
					// Check the count of the cached response
					app.SimpleMetrics.Measure("custom/gov/tally+cached", 0)

					return res
				}
			}
		}
	}

	start := time.Now()
	res = app.BaseApp.Query(req)
	elapsed := time.Since(start)

	app.SimpleMetrics.Measure(req.Path, elapsed)

	if req.Path == "custom/gov/tally" || req.Path == "/custom/gov/tally" {
		bz, err := sdk.SortJSON(req.Data)
		if err == nil {
			key := hex.EncodeToString(bz)
			queryGovTallyCache.Set(key, res, cache.DefaultExpiration)
		}
	}

	return res
}
