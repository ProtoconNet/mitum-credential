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

var assignCredentialsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AssignCredentialsItemProcessor)
	},
}

var assignCredentialsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AssignCredentialsProcessor)
	},
}

func (AssignCredentials) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AssignCredentialsItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            AssignCredentialsItem
	credentialCount *uint64
	holders         *[]Holder
}

func (ipp *AssignCredentialsItemProcessor) PreProcess(
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

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AssignCredentialsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	if st, _ := existsState(StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()), "key of credential", getStateFunc); st != nil {
		credential, err := StateCredentialValue(st)
		if err != nil {
			return nil, err
		}

		if credential.Holder() == nil {
			*ipp.credentialCount++
		}
	}

	sts := make([]base.StateMergeValue, 2)

	credential := NewCredential(it.Holder(), it.TemplateID(), it.ID(), it.Value(), it.ValidFrom(), it.ValidUntil(), it.DID())
	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts[0] = NewStateMergeValue(
		StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()),
		NewCredentialStateValue(credential),
	)

	sts[1] = NewStateMergeValue(
		StateKeyHolderDID(it.Contract(), it.CredentialServiceID(), it.Holder()),
		NewHolderDIDStateValue(it.DID()),
	)

	if len(*ipp.holders) == 0 {
		*ipp.holders = append(*ipp.holders, NewHolder(it.Holder(), 1))
	} else {
		for i, h := range *ipp.holders {
			if h.Address().Equal(it.Holder()) {
				(*ipp.holders)[i] = NewHolder(h.Address(), h.CredentialCount()+1)
				break
			}

			if i == len(*ipp.holders)-1 {
				*ipp.holders = append(*ipp.holders, NewHolder(it.Holder(), 1))
			}
		}
	}

	return sts, nil
}

func (ipp *AssignCredentialsItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AssignCredentialsItem{}
	ipp.credentialCount = nil
	ipp.holders = nil

	assignCredentialsItemProcessorPool.Put(ipp)

	return nil
}

type AssignCredentialsProcessor struct {
	*base.BaseOperationProcessor
}

func NewAssignCredentialsProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new AssignCredentialsProcessor")

		nopp := assignCredentialsProcessorPool.Get()
		opp, ok := nopp.(*AssignCredentialsProcessor)
		if !ok {
			return nil, e(nil, "expected AssignCredentialsProcessor, not %T", nopp)
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

func (opp *AssignCredentialsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess AssignCredentials")

	fact, ok := op.Fact().(AssignCredentialsFact)
	if !ok {
		return ctx, nil, e(nil, "expected AssignCredentialsFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot assign credential status, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := assignCredentialsItemProcessorPool.Get()
		ipc, ok := ip.(*AssignCredentialsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AssignCredentialsItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = nil
		ipc.holders = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess AssignCredentialsItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AssignCredentialsProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process AssignCredentials")

	fact, ok := op.Fact().(AssignCredentialsFact)
	if !ok {
		return nil, nil, e(nil, "expected AssignCredentialsFact, not %T", op.Fact())
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
		ip := assignCredentialsItemProcessorPool.Get()
		ipc, ok := ip.(*AssignCredentialsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AssignCredentialsItemProcessor, not %T", ip)
		}

		k := StateKeyDesign(it.Contract(), it.CredentialServiceID())

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = counters[k]
		ipc.holders = holders[k]

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AssignCredentialsItem: %w", err), nil
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

func (opp *AssignCredentialsProcessor) Close() error {
	assignCredentialsProcessorPool.Put(opp)

	return nil
}

func calculateCredentialItemsFee(getStateFunc base.GetStateFunc, items []CredentialItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for _, item := range items {
		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := existsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]currency.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
