package gaia

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// Querying gov tally calculates the all voting of the validators and delegators.
	// It can take a time 20sec ~ 30sec.
	// Immediate response is less important compared to how long this query takes.
	// Therefore, to reduce resource usage, cache for 5 minutes.
	queryGovTallyCache = cache.New(5*time.Minute, 10*time.Minute)
)
