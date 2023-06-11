package credential

import (
	"context"
	"sync"

	state "github.com/ProtoconNet/mitum-credential/state/credential"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencyoperation "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
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
	holders         *[]types.Holder
}

func (ipp *AssignCredentialsItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return err
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(it.Holder()), getStateFunc); err != nil {
		return err
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(it.Holder()), getStateFunc); err != nil {
		return err
	}

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
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

	if err := currencystate.CheckExistsState(state.StateKeyDesign(it.Contract(), it.CredentialServiceID()), getStateFunc); err != nil {
		return err
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AssignCredentialsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	if st, _ := currencystate.ExistsState(state.StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()), "key of credential", getStateFunc); st != nil {
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
		state.StateKeyCredential(it.Contract(), it.CredentialServiceID(), it.TemplateID(), it.ID()),
		state.NewCredentialStateValue(credential),
	)

	sts[1] = currencystate.NewStateMergeValue(
		state.StateKeyHolderDID(it.Contract(), it.CredentialServiceID(), it.Holder()),
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

func NewAssignCredentialsProcessor() currencytypes.GetNewProcessor {
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

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot assign credential status, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
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
		ip := assignCredentialsItemProcessorPool.Get()
		ipc, ok := ip.(*AssignCredentialsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AssignCredentialsItemProcessor, not %T", ip)
		}

		k := state.StateKeyDesign(it.Contract(), it.CredentialServiceID())

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
		policy := types.NewPolicy(de.Policy().Templates(), *holders[k], *counters[k])
		design := types.NewDesign(de.CredentialServiceID(), policy)
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid design, %s: %w", k, err), nil
		}

		sts = append(sts,
			currencystate.NewStateMergeValue(
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
	sb, err := currencyoperation.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currencystate.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *AssignCredentialsProcessor) Close() error {
	assignCredentialsProcessorPool.Put(opp)

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
