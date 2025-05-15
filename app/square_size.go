package app

import (
	"github.com/celestiaorg/celestia-app/v3/pkg/appconsts"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxEffectiveSquareSize returns the max effective square size.
func (app *App) MaxEffectiveSquareSize(ctx sdk.Context) int {
	height := ctx.BlockHeight()

	if height <= 1 {
		app.Logger().Info("Using default GovMaxSquareSize due to early height", "default", int(appconsts.DefaultGovMaxSquareSize))
		return int(appconsts.DefaultGovMaxSquareSize)
	}

	var govMax int
	defer func() {
		if r := recover(); r != nil {
			app.Logger().Error("GovMaxSquareSize fallback: paramStore access panic", "error", r)
			govMax = int(appconsts.DefaultGovMaxSquareSize)
		}
	}()

	if !app.BlobKeeper.HasGovMaxSquareSize(ctx) {
		app.Logger().Error("GovMaxSquareSize missing from param store, falling back to default")
		return int(appconsts.DefaultGovMaxSquareSize)
	}

	govMax = int(app.BlobKeeper.GovMaxSquareSize(ctx))
	hardMax := appconsts.SquareSizeUpperBound(app.AppVersion())
	return min(govMax, hardMax)
}
