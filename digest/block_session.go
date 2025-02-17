package digest

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-credential/state"
	crcystate "github.com/ProtoconNet/mitum-currency/v3/state"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/fixedtree"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bulkWriteLimit = 500

type BlockSession struct {
	sync.RWMutex
	block                 mitumbase.BlockMap
	ops                   []mitumbase.Operation
	opstree               fixedtree.Tree
	sts                   []mitumbase.State
	st                    *currencydigest.Database
	proposal              mitumbase.ProposalSignFact
	opsTreeNodes          map[string]mitumbase.OperationFixedtreeNode
	blockModels           []mongo.WriteModel
	operationModels       []mongo.WriteModel
	accountModels         []mongo.WriteModel
	balanceModels         []mongo.WriteModel
	currencyModels        []mongo.WriteModel
	contractAccountModels []mongo.WriteModel
	didIssuerModels       []mongo.WriteModel
	didCredentialModels   []mongo.WriteModel
	didHolderDIDModels    []mongo.WriteModel
	didTemplateModels     []mongo.WriteModel
	statesValue           *sync.Map
	balanceAddressList    []string
	credentialMap         map[string]struct{}
	buildinfo             string
}

func NewBlockSession(
	st *currencydigest.Database,
	blk mitumbase.BlockMap,
	ops []mitumbase.Operation,
	opstree fixedtree.Tree,
	sts []mitumbase.State,
	proposal mitumbase.ProposalSignFact,
	vs string,
) (*BlockSession, error) {
	if st.Readonly() {
		return nil, errors.Errorf("readonly mode")
	}

	nst, err := st.New()
	if err != nil {
		return nil, err
	}

	return &BlockSession{
		st:            nst,
		block:         blk,
		ops:           ops,
		opstree:       opstree,
		sts:           sts,
		proposal:      proposal,
		statesValue:   &sync.Map{},
		credentialMap: map[string]struct{}{},
		buildinfo:     vs,
	}, nil
}

func (bs *BlockSession) Prepare() error {
	bs.Lock()
	defer bs.Unlock()
	if err := bs.prepareOperationsTree(); err != nil {
		return err
	}
	if err := bs.prepareBlock(); err != nil {
		return err
	}
	if err := bs.prepareOperations(); err != nil {
		return err
	}
	if err := bs.prepareCurrencies(); err != nil {
		return err
	}
	if err := bs.prepareDIDCredential(); err != nil {
		return err
	}

	return bs.prepareAccounts()
}

func (bs *BlockSession) Commit(_ context.Context) error {
	bs.Lock()
	defer bs.Unlock()

	started := time.Now()
	defer func() {
		bs.statesValue.Store("commit", time.Since(started))

		_ = bs.close()
	}()

	_, err := bs.st.MongoClient().WithSession(func(txnCtx mongo.SessionContext, collection func(string) *mongo.Collection) (interface{}, error) {
		if len(bs.didCredentialModels) > 0 {
			for key := range bs.credentialMap {
				parsedKey, err := crcystate.ParseStateKey(key, state.CredentialPrefix, 5)
				if err != nil {
					return nil, err
				}
				err = bs.st.CleanByHeightColName(
					txnCtx,
					bs.block.Manifest().Height(),
					defaultColNameDIDCredential,
					bson.D{{"contract", parsedKey[1]}},
					bson.D{{"template", parsedKey[2]}},
					bson.D{{"credential_id", parsedKey[3]}},
				)
				if err != nil {
					return nil, err
				}
			}

			if err := bs.writeModels(txnCtx, defaultColNameDIDCredential, bs.didCredentialModels); err != nil {
				return nil, err
			}
		}

		if err := bs.writeModels(txnCtx, defaultColNameBlock, bs.blockModels); err != nil {
			return nil, err
		}

		if len(bs.operationModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameOperation, bs.operationModels); err != nil {
				return nil, err
			}
		}

		if len(bs.currencyModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameCurrency, bs.currencyModels); err != nil {
				return nil, err
			}
		}

		if len(bs.accountModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameAccount, bs.accountModels); err != nil {
				return nil, err
			}
		}

		if len(bs.contractAccountModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameContractAccount, bs.contractAccountModels); err != nil {
				return nil, err
			}
		}

		if len(bs.balanceModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameBalance, bs.balanceModels); err != nil {
				return nil, err
			}
		}

		if len(bs.didIssuerModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDIDCredentialService, bs.didIssuerModels); err != nil {
				return nil, err
			}
		}

		if len(bs.didHolderDIDModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameHolder, bs.didHolderDIDModels); err != nil {
				return nil, err
			}
		}

		if len(bs.didTemplateModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameTemplate, bs.didTemplateModels); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

func (bs *BlockSession) Close() error {
	bs.Lock()
	defer bs.Unlock()

	return bs.close()
}

func (bs *BlockSession) prepareOperationsTree() error {
	nodes := map[string]mitumbase.OperationFixedtreeNode{}

	if err := bs.opstree.Traverse(func(_ uint64, no fixedtree.Node) (bool, error) {
		nno := no.(mitumbase.OperationFixedtreeNode)
		if nno.InState() {
			nodes[nno.Key()] = nno
		} else {
			nodes[nno.Key()[:len(nno.Key())-1]] = nno
		}

		return true, nil
	}); err != nil {
		return err
	}

	bs.opsTreeNodes = nodes

	return nil
}

func (bs *BlockSession) prepareBlock() error {
	if bs.block == nil {
		return nil
	}

	bs.blockModels = make([]mongo.WriteModel, 1)

	manifest := isaac.NewManifest(
		bs.block.Manifest().Height(),
		bs.block.Manifest().Previous(),
		bs.block.Manifest().Proposal(),
		bs.block.Manifest().OperationsTree(),
		bs.block.Manifest().StatesTree(),
		bs.block.Manifest().Suffrage(),
		bs.block.Manifest().ProposedAt(),
	)

	doc, err := currencydigest.NewManifestDoc(manifest, bs.st.Encoder(), bs.block.Manifest().Height(), bs.ops, bs.block.SignedAt(), bs.proposal.ProposalFact().Proposer(), bs.proposal.ProposalFact().Point().Round(), bs.buildinfo)
	if err != nil {
		return err
	}
	bs.blockModels[0] = mongo.NewInsertOneModel().SetDocument(doc)

	return nil
}

func (bs *BlockSession) prepareOperations() error {
	if len(bs.ops) < 1 {
		return nil
	}

	node := func(h mitumutil.Hash) (bool, bool, mitumbase.OperationProcessReasonError) {
		no, found := bs.opsTreeNodes[h.String()]
		if !found {
			return false, false, nil
		}

		return true, no.InState(), no.Reason()
	}

	bs.operationModels = make([]mongo.WriteModel, len(bs.ops))

	for i := range bs.ops {
		op := bs.ops[i]

		var doc currencydigest.OperationDoc
		switch found, inState, reason := node(op.Fact().Hash()); {
		case !found:
			return mitumutil.ErrNotFound.Errorf("operation, %v in operations tree", op.Fact().Hash().String())
		default:
			var reasonMsg string
			switch {
			case reason == nil:
				reasonMsg = ""
			default:
				reasonMsg = reason.Msg()
			}
			d, err := currencydigest.NewOperationDoc(
				op,
				bs.st.Encoder(),
				bs.block.Manifest().Height(),
				bs.block.SignedAt(),
				inState,
				reasonMsg,
				uint64(i),
			)
			if err != nil {
				return err
			}
			doc = d
		}

		bs.operationModels[i] = mongo.NewInsertOneModel().SetDocument(doc)
	}

	return nil
}

func (bs *BlockSession) prepareAccounts() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var accountModels []mongo.WriteModel
	var balanceModels []mongo.WriteModel
	var contractAccountModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]

		switch {
		case statecurrency.IsAccountStateKey(st.Key()):
			j, err := bs.handleAccountState(st)
			if err != nil {
				return err
			}
			accountModels = append(accountModels, j...)
		case statecurrency.IsBalanceStateKey(st.Key()):
			j, address, err := bs.handleBalanceState(st)
			if err != nil {
				return err
			}
			balanceModels = append(balanceModels, j...)
			bs.balanceAddressList = append(bs.balanceAddressList, address)
		case stateextension.IsStateContractAccountKey(st.Key()):
			j, err := bs.handleContractAccountState(st)
			if err != nil {
				return err
			}
			contractAccountModels = append(contractAccountModels, j...)
		default:
			continue
		}
	}

	bs.accountModels = accountModels
	bs.contractAccountModels = contractAccountModels
	bs.balanceModels = balanceModels

	return nil
}

func (bs *BlockSession) prepareCurrencies() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var currencyModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case statecurrency.IsDesignStateKey(st.Key()):
			j, err := bs.handleCurrencyState(st)
			if err != nil {
				return err
			}
			currencyModels = append(currencyModels, j...)
		default:
			continue
		}
	}

	bs.currencyModels = currencyModels

	return nil
}

func (bs *BlockSession) writeModels(ctx context.Context, col string, models []mongo.WriteModel) error {
	started := time.Now()
	defer func() {
		bs.statesValue.Store(fmt.Sprintf("write-models-%s", col), time.Since(started))
	}()

	n := len(models)
	if n < 1 {
		return nil
	} else if n <= bulkWriteLimit {
		return bs.writeModelsChunk(ctx, col, models)
	}

	z := n / bulkWriteLimit
	if n%bulkWriteLimit != 0 {
		z++
	}

	for i := 0; i < z; i++ {
		s := i * bulkWriteLimit
		e := s + bulkWriteLimit
		if e > n {
			e = n
		}

		if err := bs.writeModelsChunk(ctx, col, models[s:e]); err != nil {
			return err
		}
	}

	return nil
}

func (bs *BlockSession) writeModelsChunk(ctx context.Context, col string, models []mongo.WriteModel) error {
	opts := options.BulkWrite().SetOrdered(false)
	if res, err := bs.st.MongoClient().Collection(col).BulkWrite(ctx, models, opts); err != nil {
		return err
	} else if res != nil && res.InsertedCount < 1 {
		return errors.Errorf("not inserted to %s", col)
	}

	return nil
}

func (bs *BlockSession) close() error {
	bs.block = nil
	bs.operationModels = nil
	bs.currencyModels = nil
	bs.accountModels = nil
	bs.balanceModels = nil
	bs.contractAccountModels = nil
	bs.didIssuerModels = nil
	bs.didCredentialModels = nil
	bs.didHolderDIDModels = nil
	bs.didTemplateModels = nil
	bs.credentialMap = nil

	return bs.st.Close()
}
