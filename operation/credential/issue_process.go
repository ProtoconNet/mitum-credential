package credential

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var issueItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueItemProcessor)
	},
}

var issueProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueProcessor)
	},
}

func (Issue) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type IssueItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            IssueItem
	credentialCount *uint64
	holders         *[]types.Holder
}

func (ipp *IssueItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess IssueItemProcessor")
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.DesignStateKey(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", it.Currency())))
	}

	if _, _, _, cErr := currencystate.ExistsCAccount(
		it.Holder(), "holder", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v: holder %v is contract account", cErr, it.Holder())))
	}

	_, cSt, aErr, cErr := currencystate.ExistsCAccount(it.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(cErr)
	}

	_, err := extensioncurrency.CheckCAAuthFromState(cSt, ipp.sender)
	if err != nil {
		return e.Wrap(err)
	}

	if st, err := currencystate.ExistsState(state.StateKeyDesign(it.Contract()), "design", getStateFunc); err != nil {
		return e.Wrap(
			common.ErrServiceNF.Errorf("credential design state for contract account %v", it.Contract()))
	} else if de, err := state.StateDesignValue(st); err != nil {
		return e.Wrap(
			common.ErrServiceNF.Errorf("credential design state value for contract account %v", it.Contract()))
	} else {
		if err := de.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
		for i, v := range de.Policy().TemplateIDs() {
			if it.templateID == v {
				break
			}
			if i == len(de.Policy().TemplateIDs())-1 {
				return e.Wrap(
					common.ErrValueInvalid.Errorf(
						"templateID %v not registered in contract account %v", it.TemplateID(), it.Contract()))
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyCredential(it.Contract(),
		it.TemplateID(),
		it.CredentialID())); {
	case err != nil:
		return e.Wrap(common.ErrStateNF.Errorf(
			"credential %v for template id %v in contract account %v", it.CredentialID(), it.TemplateID(), it.Contract()))
	case !found:
	default:
		if credential, isActive, err := state.StateCredentialValue(st); err != nil {
			return e.Wrap(
				common.ErrStateValInvalid.Errorf(
					"credential %v for template id %v in contract account %v",
					it.CredentialID(), it.TemplateID(), it.Contract()))
		} else if isActive {
			return e.Wrap(
				common.ErrValueInvalid.Errorf(
					"credential %v for template %v is already assigned to holder %v in contract account %v",
					it.CredentialID(), it.TemplateID(), credential.Holder(), it.Contract()))
		}
	}

	return nil
}

func (ipp *IssueItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	*ipp.credentialCount++

	var sts []base.StateMergeValue

	smv, err := currencystate.CreateNotExistAccount(it.Holder(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	credential := types.NewCredential(it.Holder(), it.TemplateID(), it.CredentialID(), it.Value(), it.ValidFrom(), it.ValidUntil(), it.DID())
	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyCredential(it.Contract(), it.TemplateID(), it.CredentialID()),
		state.NewCredentialStateValue(credential, true),
	))

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyHolderDID(it.Contract(), it.Holder()),
		state.NewHolderDIDStateValue(it.DID()),
	))

	if len(*ipp.holders) == 0 {
		*ipp.holders = append(*ipp.holders, types.NewHolder(it.Holder(), 1))
	} else {
		for i, h := range *ipp.holders {
			if h.Address().Equal(it.Holder()) {
				(*ipp.holders)[i] = types.NewHolder(h.Address(), h.CredentialCount()+1)
				break
			}

			if i == len(*ipp.holders)-1 {
				*ipp.holders = append(*ipp.holders, types.NewHolder(it.Holder(), 1))
			}
		}
	}

	return sts, nil
}

func (ipp *IssueItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = IssueItem{}
	ipp.credentialCount = nil
	ipp.holders = nil

	issueItemProcessorPool.Put(ipp)
}

type IssueProcessor struct {
	*base.BaseOperationProcessor
}

func NewIssueProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AssignProcessor")

		nopp := issueProcessorPool.Get()
		opp, ok := nopp.(*IssueProcessor)
		if !ok {
			return nil, e.Errorf("expected AssignProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *IssueProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(IssueFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", IssueFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender %v is contract account", cErr, fact.Sender())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, it := range fact.Items() {
		ip := issueItemProcessorPool.Get()
		ipc, ok := ip.(*IssueItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected AssignItemProcessor, not %T", ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = nil
		ipc.holders = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *IssueProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Assign")

	fact, _ := op.Fact().(IssueFact)
	designs := map[string]types.Design{}
	counters := map[string]*uint64{}
	holders := map[string]*[]types.Holder{}

	for _, it := range fact.Items() {
		k := state.StateKeyDesign(it.Contract())

		if _, found := counters[k]; found {
			continue
		}

		st, _ := currencystate.ExistsState(k, "design", getStateFunc)

		design, _ := state.StateDesignValue(st)
		count := design.Policy().CredentialCount()
		holder := design.Policy().Holders()

		designs[k] = design
		counters[k] = &count
		holders[k] = &holder
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := issueItemProcessorPool.Get()
		ipc, _ := ip.(*IssueItemProcessor)

		k := state.StateKeyDesign(it.Contract())

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = counters[k]
		ipc.holders = holders[k]

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AssignItem; %w", err), nil
		}

		sts = append(sts, st...)

		ipc.Close()
	}

	for k, de := range designs {
		policy := types.NewPolicy(de.Policy().TemplateIDs(), *holders[k], *counters[k])
		design := types.NewDesign(policy)
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid design, %s; %w", k, err), nil
		}

		sts = append(sts,
			currencystate.NewStateMergeValue(
				k,
				state.NewDesignStateValue(design),
			),
		)
	}

	items := make([]CredentialItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	feeReceiverBalSts, required, err := calculateCredentialItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for cid := range sb {
		v, ok := sb[cid].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[cid].Value())
		}

		_, feeReceiverFound := feeReceiverBalSts[cid]

		if feeReceiverFound && (sb[cid].Key() != feeReceiverBalSts[cid].Key()) {
			stmv := common.NewBaseStateMergeValue(
				sb[cid].Key(),
				statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(required[cid][1])),
				func(height base.Height, st base.State) base.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiverBalSts[cid].Value().(statecurrency.BalanceStateValue)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", statecurrency.BalanceStateValue{}, feeReceiverBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiverBalSts[cid].Key(),
					statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height base.Height, st base.State) base.StateValueMerger {
						return statecurrency.NewBalanceStateValueMerger(height, feeReceiverBalSts[cid].Key(), cid, st)
					},
				),
			)

			sts = append(sts, stmv)
		}
	}

	return sts, nil, nil
}

func (opp *IssueProcessor) Close() error {
	issueProcessorPool.Put(opp)

	return nil
}

func calculateCredentialItemsFee(getStateFunc base.GetStateFunc, items []CredentialItem) (
	map[currencytypes.CurrencyID]base.State, map[currencytypes.CurrencyID][2]common.Big, error) {
	feeReceiveSts := map[currencytypes.CurrencyID]base.State{}
	required := map[currencytypes.CurrencyID][2]common.Big{}

	for _, item := range items {
		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := currencystate.ExistsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, nil, err
		}

		switch k, err := policy.Feeer().Fee(common.ZeroBig); {
		case err != nil:
			return nil, nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]common.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]common.Big{rq[0].Add(k), rq[1].Add(k)}
		}

		if policy.Feeer().Receiver() == nil {
			continue
		}

		if err := currencystate.CheckExistsState(statecurrency.AccountStateKey(policy.Feeer().Receiver()), getStateFunc); err != nil {
			return nil, nil, err
		} else if st, found, err := getStateFunc(statecurrency.BalanceStateKey(policy.Feeer().Receiver(), item.Currency())); err != nil {
			return nil, nil, err
		} else if !found {
			return nil, nil, errors.Errorf("feeer receiver account not found, %s", policy.Feeer().Receiver())
		} else {
			feeReceiveSts[item.Currency()] = st
		}

	}

	return feeReceiveSts, required, nil

}
