package source

import (
	"cosmossdk.io/math"
	minttypes "github.com/realiotech/realio-network/x/mint/types"
)

type Source interface {
	GetInflation(height int64) (math.LegacyDec, error)
	Params(height int64) (minttypes.Params, error)
}
