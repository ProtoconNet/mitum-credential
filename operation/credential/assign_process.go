package credential

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var assignItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AssignItemProcessor)
	},
}

var assignProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AssignProcessor)
	},
}

func (Assign) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AssignItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            AssignItem
	credentialCount *uint64
	holders         *[]types.Holder
}

func (ipp *AssignItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return errors.Wrap(err, " invalid AssignItem")
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return errors.Errorf("failed to get fee Currency %s state", it.Currency())
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(it.Holder()), getStateFunc); err != nil {
		return errors.Wrapf(err, "failed to get Holder %s state", it.Holder())
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(it.Holder()), getStateFunc); err != nil {
		return errors.Wrapf(err, "Holder %s is contract account, contract account cannot be holder", it.Holder())
	}

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return errors.Wrapf(err, "failed to get target contract account %s state", it.Contract())
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return errors.Wrap(err, "failed to get contract account value from state")
	}

	if !(ca.Owner().Equal(ipp.sender) || ca.IsOperator(ipp.sender)) {
		return errors.Errorf(
			"sender is neither the owner nor the operator of the target contract account, %q",
			ipp.sender,
		)
	}

	if st, err := currencystate.ExistsState(state.StateKeyDesign(it.Contract()), "key of design", getStateFunc); err != nil {
		return errors.Wrapf(err, "failed to get design state of credential service")
	} else if _, err = state.StateDesignValue(st); err != nil {
		return errors.Wrapf(err, "failed to get design value of credential service from state")
	}

	return nil
}

func (ipp *AssignItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	if st, _ := currencystate.ExistsState(
		state.StateKeyCredential(it.Contract(),
			it.TemplateID(),
			it.ID()), "key of credential",
		getStateFunc,
	); st != nil {
		credential, err := state.StateCredentialValue(st)
		if err != nil {
			return nil, err
		}

		if credential.Holder() == nil {
			*ipp.credentialCount++
		}
	}

	sts := make([]base.StateMergeValue, 2)

	credential := types.NewCredential(it.Holder(), it.TemplateID(), it.ID(), it.Value(), it.ValidFrom(), it.ValidUntil(), it.DID())
	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts[0] = currencystate.NewStateMergeValue(
		state.StateKeyCredential(it.Contract(), it.TemplateID(), it.ID()),
		state.NewCredentialStateValue(credential),
	)

	sts[1] = currencystate.NewStateMergeValue(
		state.StateKeyHolderDID(it.Contract(), it.Holder()),
		state.NewHolderDIDStateValue(it.DID()),
	)

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

func (ipp *AssignItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AssignItem{}
	ipp.credentialCount = nil
	ipp.holders = nil

	assignItemProcessorPool.Put(ipp)
}

type AssignProcessor struct {
	*base.BaseOperationProcessor
}

func NewAssignProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AssignProcessor")

		nopp := assignProcessorPool.Get()
		opp, ok := nopp.(*AssignProcessor)
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

func (opp *AssignProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Assign")

	fact, ok := op.Fact().(AssignFact)
	if !ok {
		return ctx, nil, e.Errorf("expected AssignFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(
		stateextension.StateKeyContractAccount(fact.Sender()),
		getStateFunc,
	); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			"sender is contract account and contract account cannot assign credential status, %q; %w",
			fact.Sender(),
			err,
		), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := assignItemProcessorPool.Get()
		ipc, ok := ip.(*AssignItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected AssignItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = nil
		ipc.holders = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to preprocess AssignItem; %w",
				err,
			), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AssignProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Assign")

	fact, _ := op.Fact().(AssignFact)
	designs := map[string]types.Design{}
	counters := map[string]*uint64{}
	holders := map[string]*[]types.Holder{}

	for _, it := range fact.Items() {
		k := state.StateKeyDesign(it.Contract())

		if _, found := counters[k]; found {
			continue
		}

		st, _ := currencystate.ExistsState(k, "key of design", getStateFunc)

		design, _ := state.StateDesignValue(st)
		count := design.Policy().CredentialCount()
		holder := design.Policy().Holders()

		designs[k] = design
		counters[k] = &count
		holders[k] = &holder
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := assignItemProcessorPool.Get()
		ipc, _ := ip.(*AssignItemProcessor)

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

	required, err := calculateCredentialItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currencystate.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *AssignProcessor) Close() error {
	assignProcessorPool.Put(opp)

	return nil
}

func calculateCredentialItemsFee(getStateFunc base.GetStateFunc, items []CredentialItem) (map[currencytypes.CurrencyID][2]common.Big, error) {
	required := map[currencytypes.CurrencyID][2]common.Big{}

	for _, item := range items {
		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := currencystate.ExistsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(common.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]common.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]common.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
