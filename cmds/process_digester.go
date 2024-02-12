package cmds

import (
	"context"
	"github.com/ProtoconNet/mitum-credential/digest"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/logging"
	"sync"
)

var blockSessionPool = sync.Pool{
	New: func() interface{} {

		bs := currencydigest.NewBlockSession()
		bs.SetPrepareFuncs(currencydigest.PrepareFuncs)
		bs.SetPrepareFuncs(digest.PrepareFuncs)
		bs.SetHandlerFuncs(currencydigest.HandlerFuncs)
		bs.SetHandlerFuncs(digest.HandlerFuncs)
		bs.SetCommitFuncs(currencydigest.CommitFuncs)
		bs.SetCommitFuncs(digest.CommitFuncs)

		return bs
	},
}

func ProcessDigester(ctx context.Context) (context.Context, error) {
	var log *logging.Logging
	if err := util.LoadFromContextOK(ctx, launch.LoggingContextKey, &log); err != nil {
		return ctx, err
	}

	var st *currencydigest.Database
	if err := util.LoadFromContext(ctx, currencycmds.ContextValueDigestDatabase, &st); err != nil {
		return ctx, err
	}

	if st == nil {
		return ctx, nil
	}

	var design launch.NodeDesign
	if err := util.LoadFromContext(ctx, launch.DesignContextKey, &design); err != nil {
		return ctx, err
	}
	root := launch.LocalFSDataDirectory(design.Storage.Base)

	var newReaders func(context.Context, string, *isaac.BlockItemReadersArgs) (*isaac.BlockItemReaders, error)
	var fromRemotes isaac.RemotesBlockItemReadFunc

	if err := util.LoadFromContextOK(ctx,
		launch.NewBlockItemReadersFuncContextKey, &newReaders,
		launch.RemotesBlockItemReaderFuncContextKey, &fromRemotes,
	); err != nil {
		return ctx, err
	}

	var sourceReaders *isaac.BlockItemReaders

	switch i, err := newReaders(ctx, root, nil); {
	case err != nil:
		return ctx, err
	default:
		sourceReaders = i
	}

	di := currencydigest.NewDigester(st, root, sourceReaders, fromRemotes, design.NetworkID, nil, &blockSessionPool)
	_ = di.SetLogging(log)

	ctx = context.WithValue(ctx, currencycmds.ContextValueBlockSession, &blockSessionPool)

	return context.WithValue(ctx, currencycmds.ContextValueDigester, di), nil
}
