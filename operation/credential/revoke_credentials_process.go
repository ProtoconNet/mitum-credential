package credential

import (
	"context"
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

var revokeCredentialsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RevokeCredentialsItemProcessor)
	},
}

var revokeCredentialsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RevokeCredentialsProcessor)
	},
}

func (RevokeCredentials) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RevokeCredentialsItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            RevokeCredentialsItem
	credentialCount *uint64
	holders         *[]types.Holder
}

func (ipp *RevokeCredentialsItemProcessor) PreProcess(
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

	if !ca.Owner().Equal(ipp.sender) {
		return errors.Errorf("sender is not contract owner, %s", ipp.sender)
	}

	if err := currencystate.CheckExistsState(state.StateKeyDesign(it.Contract(), it.CredentialServiceID()), getStateFunc); err != nil {
		return err
	}

	st, err = currencystate.ExistsState(state.StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()), "key of credential", getStateFunc)
	if err != nil {
		return err
	}

	c, err := state.StateCredentialValue(st)
	if err != nil {
		return err
	}

	if c.Holder() == nil {
		return errors.Errorf("already revoked credential, %s-%s-%s, %s", it.Contract(), it.CredentialServiceID(), it.ID(), c.Holder())
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RevokeCredentialsItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	st, err := currencystate.ExistsState(state.StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()), "key of credential", getStateFunc)
	if err != nil {
		return nil, err
	}

	credential, err := state.StateCredentialValue(st)
	if err != nil {
		return nil, err
	}

	credential = types.NewCredential(nil, credential.TemplateID(), credential.ID(), credential.Value(), credential.ValidFrom(), credential.ValidUntil(), credential.DID())
	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts := []base.StateMergeValue{
		state.NewStateMergeValue(
			state.StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()),
			state.NewCredentialStateValue(credential),
		),
	}

	*ipp.credentialCount--

	holders := *ipp.holders

	if len(holders) == 0 {
		return nil, errors.Errorf("empty holders, %s-%s", it.Contract(), it.CredentialServiceID())
	}

	for i, h := range holders {
		if h.Address().Equal(it.Holder()) {
			if h.CredentialCount()-1 == 0 {
				if i < len(holders)-1 {
					copy(holders[i:], holders[i+1:])
				}
				holders = holders[:len(holders)-1]
			} else {
				holders[i] = types.NewHolder(h.Address(), h.CredentialCount()-1)
			}
			break
		}

		if i == len(holders)-1 {
			return nil, errors.Errorf("holder not found in credential service holders, %s-%s, %s", it.Contract(), it.CredentialServiceID(), it.Holder())
		}
	}

	ipp.holders = &holders

	return sts, nil
}

func (ipp *RevokeCredentialsItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RevokeCredentialsItem{}
	ipp.credentialCount = nil
	ipp.holders = nil

	revokeCredentialsItemProcessorPool.Put(ipp)

	return nil
}

type RevokeCredentialsProcessor struct {
	*base.BaseOperationProcessor
}

func NewRevokeCredentialsProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RevokeCredentialsProcessor")

		nopp := revokeCredentialsProcessorPool.Get()
		opp, ok := nopp.(*RevokeCredentialsProcessor)
		if !ok {
			return nil, e.Errorf("expected RevokeCredentialsProcessor, not %T", nopp)
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

func (opp *RevokeCredentialsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess RevokeCredentials")

	fact, ok := op.Fact().(RevokeCredentialsFact)
	if !ok {
		return ctx, nil, e.Errorf("expected RevokeCredentialsFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot revoke credential status, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := revokeCredentialsItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeCredentialsItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected RevokeCredentialsItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = nil
		ipc.holders = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess RevokeCredentialsItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RevokeCredentialsProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process RevokeCredentials")

	fact, ok := op.Fact().(RevokeCredentialsFact)
	if !ok {
		return nil, nil, e.Errorf("expected RevokeCredentialsFact, not %T", op.Fact())
	}

	designs := map[string]types.Design{}
	counters := map[string]*uint64{}
	holders := map[string]*[]types.Holder{}

	for _, it := range fact.Items() {
		k := state.StateKeyDesign(it.Contract(), it.CredentialServiceID())

		if _, found := counters[k]; found {
			continue
		}

		st, err := currencystate.ExistsState(k, "key of design", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("credential service not found, %s-%s:%w", it.Contract(), it.CredentialServiceID(), err), nil
		}

		design, err := state.StateDesignValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("credential service value not found, %s-%s:%w", it.Contract(), it.CredentialServiceID(), err), nil
		}

		designs[k] = design

		count := design.Policy().CredentialCount()
		holder := design.Policy().Holders()

		counters[k] = &count
		holders[k] = &holder
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := revokeCredentialsItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeCredentialsItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected RevokeCredentialsItemProcessor, not %T", ip)
		}

		k := state.StateKeyDesign(it.Contract(), it.CredentialServiceID())

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = counters[k]
		ipc.holders = holders[k]

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RevokeCredentialsItem: %w", err), nil
		}

		sts = append(sts, st...)

		ipc.Close()
	}

	for k, de := range designs {
		policy := types.NewPolicy(de.Policy().TemplateIDs(), *holders[k], *counters[k])
		design := types.NewDesign(de.CredentialServiceID(), policy)
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid design, %s: %w", k, err), nil
		}

		sts = append(sts,
			state.NewStateMergeValue(
				k,
				state.NewDesignStateValue(design),
			),
		)
	}

	fitems := fact.Items()
	items := make([]CredentialItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := calculateCredentialItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, state.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *RevokeCredentialsProcessor) Close() error {
	revokeCredentialsProcessorPool.Put(opp)

	return nil
}
