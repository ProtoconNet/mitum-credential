package credential

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	statec "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stetee "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
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
	e := util.StringError("process RevokeItemProcessor")
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := cstate.CheckExistsState(statec.DesignStateKey(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", it.Currency())))
	}

	if _, _, aErr, cErr := cstate.ExistsCAccount(it.Holder(), "holder", true, false, getStateFunc); aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(cErr))
	}

	_, cSt, aErr, cErr := cstate.ExistsCAccount(it.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(cErr)
	}

	_, err := stetee.CheckCAAuthFromState(cSt, ipp.sender)
	if err != nil {
		return e.Wrap(err)
	}

	if st, err := cstate.ExistsState(state.StateKeyDesign(it.Contract()), "design", getStateFunc); err != nil {
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
				return e.Wrap(common.ErrValueInvalid.Errorf("not registered template %v", it.TemplateID()))
			}
		}
	}

	st, err := cstate.ExistsState(
		state.StateKeyCredential(it.Contract(), it.TemplateID(), it.CredentialID()), "credential", getStateFunc)
	if err != nil {
		return e.Wrap(
			common.ErrStateNF.Errorf(
				"credential %v for template %v in contract account %v", it.CredentialID(), it.TemplateID(), it.Contract()))
	}

	crd, isActive, err := state.StateCredentialValue(st)
	if err != nil {
		return e.Wrap(
			common.ErrStateValInvalid.Errorf(
				"credential %v for template %v in contract account %v", it.CredentialID(), it.TemplateID(), it.Contract()))
	}

	if !crd.Holder().Equal(it.Holder()) {
		return e.Wrap(
			common.ErrValueInvalid.Errorf(
				"holder %v has not owned credential %v for template %v in contract account %v",
				it.Holder(), it.CredentialID(), it.TemplateID(), it.Contract()))
	}

	if !isActive {
		return e.Wrap(common.ErrValueInvalid.Errorf(
			"already revoked credential %v for template %v in contract account %v",
			it.CredentialID(), it.TemplateID(), it.Contract()))
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

	st, _ := cstate.ExistsState(state.StateKeyCredential(it.Contract(), it.TemplateID(), it.CredentialID()), "credential", getStateFunc)
	credential, _, _ := state.StateCredentialValue(st)

	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts := []base.StateMergeValue{
		cstate.NewStateMergeValue(
			state.StateKeyCredential(it.Contract(), it.TemplateID(), it.CredentialID()),
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

func NewRevokeProcessor() ctypes.GetNewProcessor {
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
	fact, ok := op.Fact().(RevokeFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", RevokeFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if _, _, aErr, cErr := cstate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender is contract account, %q", fact.Sender(), cErr)), nil
	}

	if err := cstate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, it := range fact.Items() {
		ip := revokeItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected RevokeItemProcessor, not %T", ip)), nil
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

		st, _ := cstate.ExistsState(k, "design", getStateFunc)

		design, err := state.StateDesignValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("credential design value not found, %s; %w", it.Contract(), err), nil
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
			cstate.NewStateMergeValue(
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
		v, ok := sb[cid].Value().(statec.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[cid].Value())
		}

		_, feeReceiverFound := feeReceiverBalSts[cid]

		if feeReceiverFound && (sb[cid].Key() != feeReceiverBalSts[cid].Key()) {
			stmv := common.NewBaseStateMergeValue(
				sb[cid].Key(),
				statec.NewDeductBalanceStateValue(v.Amount.WithBig(required[cid][1])),
				func(height base.Height, st base.State) base.StateValueMerger {
					return statec.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiverBalSts[cid].Value().(statec.BalanceStateValue)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", statec.BalanceStateValue{}, feeReceiverBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiverBalSts[cid].Key(),
					statec.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height base.Height, st base.State) base.StateValueMerger {
						return statec.NewBalanceStateValueMerger(height, feeReceiverBalSts[cid].Key(), cid, st)
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
