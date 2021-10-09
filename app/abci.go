package gaia

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/patrickmn/go-cache"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *GaiaApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if req.Path == "/cosmos.staking.v1beta1.Query/ValidatorDelegations" {
		// To check the number of forbidden requests
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	res = app.BaseApp.Query(req)

	if req.Path == "custom/gov/tally" || req.Path == "/custom/gov/tally" {
		bz, err := sdk.SortJSON(req.Data)
		if err == nil {
			key := hex.EncodeToString(bz)
			queryGovTallyCache.Set(key, res, cache.DefaultExpiration)
		}
	}

	return res
}
