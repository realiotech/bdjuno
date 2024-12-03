package types

import (
	"fmt"

	"github.com/forbole/juno/v6/node/remote"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v6/node/local"
	mintkeeper "github.com/realiotech/realio-network/x/mint/keeper"
	minttypes "github.com/realiotech/realio-network/x/mint/types"

	nodeconfig "github.com/forbole/juno/v6/node/config"

	banksource "github.com/forbole/callisto/v4/modules/bank/source"
	localbanksource "github.com/forbole/callisto/v4/modules/bank/source/local"
	remotebanksource "github.com/forbole/callisto/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/callisto/v4/modules/distribution/source"
	localdistrsource "github.com/forbole/callisto/v4/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/callisto/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/callisto/v4/modules/gov/source"
	localgovsource "github.com/forbole/callisto/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/callisto/v4/modules/gov/source/remote"
	mintsource "github.com/forbole/callisto/v4/modules/mint/source"
	localmintsource "github.com/forbole/callisto/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/callisto/v4/modules/mint/source/remote"
	slashingsource "github.com/forbole/callisto/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/callisto/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/callisto/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/callisto/v4/modules/staking/source"
	localstakingsource "github.com/forbole/callisto/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/callisto/v4/modules/staking/source/remote"
	"github.com/forbole/callisto/v4/utils/simapp"
	realioApp "github.com/realiotech/realio-network/app"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, cdc codec.Codec) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, cdc)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, cdc codec.Codec) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, cdc)
	if err != nil {
		return nil, err
	}

	app := simapp.NewSimApp(cdc)

	realioApp := realioApp.New(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{},
	)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrkeeper.NewQuerier(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govkeeper.NewQueryServer(&app.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, mintkeeper.NewQueryServerImpl(realioApp.MintKeeper)),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
