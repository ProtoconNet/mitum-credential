package credential

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var revokeItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RevokeItemProcessor)
	},
}

var revokeProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RevokeProcessor)
	},
}

func (Revoke) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RevokeItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            RevokeItem
	credentialCount *uint64
	holders         *[]types.Holder
}

func (ipp *RevokeItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return err
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(it.Holder()), getStateFunc); err != nil {
		return err
	}

	if err := currencystate.CheckNotExistsState(extension.StateKeyContractAccount(it.Holder()), getStateFunc); err != nil {
		return err
	}

	st, err := currencystate.ExistsState(extension.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return err
	}

	ca, err := extension.StateContractAccountValue(st)
	if err != nil {
		return err
	}

	if !(ca.Owner().Equal(ipp.sender) || ca.IsOperator(ipp.sender)) {
		return errors.Errorf(
			"sender is neither the owner nor the operator of the target contract account, %q",
			ipp.sender,
		)
	}

	if st, err := currencystate.ExistsState(state.StateKeyDesign(it.Contract()), "key of design", getStateFunc); err != nil {
		return errors.Wrapf(err, "failed to get design state of credential service")
	} else if de, err := state.StateDesignValue(st); err != nil {
		return errors.Wrapf(err, "failed to get design value of credential service from state")
	} else {
		if err := de.IsValid(nil); err != nil {
			return err
		}
		for i, v := range de.Policy().TemplateIDs() {
			if it.templateID == v {
				break
			}
			if i == len(de.Policy().TemplateIDs())-1 {
				return errors.Errorf("templateID not found")
			}
		}
	}

	st, err = currencystate.ExistsState(state.StateKeyCredential(it.Contract(), it.TemplateID(), it.ID()), "key of credential", getStateFunc)
	if err != nil {
		return err
	}

	credential, isActive, err := state.StateCredentialValue(st)
	if err != nil {
		return err
	}

	if !isActive {
		return errors.Errorf("already revoked credential, %s-%s, %s", it.Contract(), it.ID(), credential.Holder())
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RevokeItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	if *ipp.credentialCount < 1 {
		return nil, errors.Errorf("no credentials to revoke")
	}
	it := ipp.item

	*ipp.credentialCount--

	if len(*ipp.holders) < 1 {
		return nil, errors.Errorf("empty holders, %s", it.Contract())
	}

	st, _ := currencystate.ExistsState(state.StateKeyCredential(it.Contract(), it.TemplateID(), it.ID()), "key of credential", getStateFunc)
	credential, _, _ := state.StateCredentialValue(st)

	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts := []base.StateMergeValue{
		currencystate.NewStateMergeValue(
			state.StateKeyCredential(it.Contract(), it.TemplateID(), it.ID()),
			state.NewCredentialStateValue(credential, false),
		),
	}

	var holders []types.Holder
	for i, h := range *ipp.holders {
		if h.Address().Equal(it.Holder()) {
			if h.CredentialCount()-1 == 0 {
				copy(holders, (*ipp.holders)[:i])
				copy(holders, (*ipp.holders)[i+1:])
				ipp.holders = &holders
			} else {
				(*ipp.holders)[i] = types.NewHolder(h.Address(), h.CredentialCount()-1)
			}
			break
		}

		if i == len(holders)-1 {
			return nil, errors.Errorf("holder not found in credential service holders, %s, %s", it.Contract(), it.Holder())
		}
	}
	return sts, nil
}

func (ipp *RevokeItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RevokeItem{}
	ipp.credentialCount = nil
	ipp.holders = nil

	revokeItemProcessorPool.Put(ipp)
}

type RevokeProcessor struct {
	*base.BaseOperationProcessor
}

func NewRevokeProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RevokeProcessor")

		nopp := revokeProcessorPool.Get()
		opp, ok := nopp.(*RevokeProcessor)
		if !ok {
			return nil, e.Errorf("expected RevokeProcessor, not %T", nopp)
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

func (opp *RevokeProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Revoke")

	fact, ok := op.Fact().(RevokeFact)
	if !ok {
		return ctx, nil, e.Errorf("expected RevokeFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot revoke credential status, %q; %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := revokeItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected RevokeItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = nil
		ipc.holders = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess RevokeItem; %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RevokeProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Revoke")

	fact, ok := op.Fact().(RevokeFact)
	if !ok {
		return nil, nil, e.Errorf("expected RevokeFact, not %T", op.Fact())
	}

	designs := map[string]types.Design{}
	counters := map[string]*uint64{}
	holders := map[string]*[]types.Holder{}

	for _, it := range fact.Items() {
		k := state.StateKeyDesign(it.Contract())

		if _, found := counters[k]; found {
			continue
		}

		st, _ := currencystate.ExistsState(k, "key of design", getStateFunc)

		design, err := state.StateDesignValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("credential service value not found, %s; %w", it.Contract(), err), nil
		}

		designs[k] = design

		count := design.Policy().CredentialCount()
		holder := design.Policy().Holders()

		counters[k] = &count
		holders[k] = &holder
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := revokeItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected RevokeItemProcessor, not %T", ip)
		}

		k := state.StateKeyDesign(it.Contract())

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = counters[k]
		ipc.holders = holders[k]

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RevokeItem; %w", err), nil
		}

		holders[k] = ipc.holders
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

func (opp *RevokeProcessor) Close() error {
	revokeProcessorPool.Put(opp)

	return nil
}
