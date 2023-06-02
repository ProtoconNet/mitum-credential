package credential

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
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
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RevokeCredentialsItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            RevokeCredentialsItem
	credentialCount *uint64
	holders         *[]Holder
}

func (ipp *RevokeCredentialsItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return err
	}

	if err := checkExistsState(currency.StateKeyAccount(it.Holder()), getStateFunc); err != nil {
		return err
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(it.Holder()), getStateFunc); err != nil {
		return err
	}

	st, err := existsState(extensioncurrency.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return err
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return err
	}

	if !ca.Owner().Equal(ipp.sender) {
		return errors.Errorf("sender is not contract owner, %s", ipp.sender)
	}

	if err := checkExistsState(StateKeyDesign(it.Contract(), it.CredentialServiceID()), getStateFunc); err != nil {
		return err
	}

	st, err = existsState(StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()), "key of credential", getStateFunc)
	if err != nil {
		return err
	}

	c, err := StateCredentialValue(st)
	if err != nil {
		return err
	}

	if c.Holder() == nil {
		return errors.Errorf("already revoked credential, %s-%s-%s, %s", it.Contract(), it.CredentialServiceID(), it.ID(), c.Holder())
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RevokeCredentialsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	st, err := existsState(StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()), "key of credential", getStateFunc)
	if err != nil {
		return nil, err
	}

	credential, err := StateCredentialValue(st)
	if err != nil {
		return nil, err
	}

	credential = NewCredential(nil, credential.TemplateID(), credential.ID(), credential.Value(), credential.ValidFrom(), credential.ValidUntil(), credential.DID())
	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts := []base.StateMergeValue{
		NewStateMergeValue(
			StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()),
			NewCredentialStateValue(credential),
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
				holders[i] = NewHolder(h.Address(), h.CredentialCount()-1)
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

func NewRevokeCredentialsProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new RevokeCredentialsProcessor")

		nopp := revokeCredentialsProcessorPool.Get()
		opp, ok := nopp.(*RevokeCredentialsProcessor)
		if !ok {
			return nil, e(nil, "expected RevokeCredentialsProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *RevokeCredentialsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess RevokeCredentials")

	fact, ok := op.Fact().(RevokeCredentialsFact)
	if !ok {
		return ctx, nil, e(nil, "expected RevokeCredentialsFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot revoke credential status, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := revokeCredentialsItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeCredentialsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected RevokeCredentialsItemProcessor, not %T", ip)
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
	e := util.StringErrorFunc("failed to process RevokeCredentials")

	fact, ok := op.Fact().(RevokeCredentialsFact)
	if !ok {
		return nil, nil, e(nil, "expected RevokeCredentialsFact, not %T", op.Fact())
	}

	designs := map[string]Design{}
	counters := map[string]*uint64{}
	holders := map[string]*[]Holder{}

	for _, it := range fact.Items() {
		k := StateKeyDesign(it.Contract(), it.CredentialServiceID())

		if _, found := counters[k]; found {
			continue
		}

		st, err := existsState(k, "key of design", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("credential service not found, %s-%s:%w", it.Contract(), it.CredentialServiceID(), err), nil
		}

		design, err := StateDesignValue(st)
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
			return nil, nil, e(nil, "expected RevokeCredentialsItemProcessor, not %T", ip)
		}

		k := StateKeyDesign(it.Contract(), it.CredentialServiceID())

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
		policy := NewPolicy(de.Policy().Templates(), *holders[k], *counters[k])
		design := NewDesign(de.CredentialServiceID(), policy)
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid design, %s: %w", k, err), nil
		}

		sts = append(sts,
			NewStateMergeValue(
				k,
				NewDesignStateValue(design),
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
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *RevokeCredentialsProcessor) Close() error {
	revokeCredentialsProcessorPool.Put(opp)

	return nil
}
