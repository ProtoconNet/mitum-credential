package credential

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-credential/state"
	crendentialtypes "github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"

	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var createServiceProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateServiceProcessor)
	},
}

func (CreateService) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateServiceProcessor struct {
	*base.BaseOperationProcessor
}

func NewCreateServiceProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateServiceProcessor")

		nopp := createServiceProcessorPool.Get()
		opp, ok := nopp.(*CreateServiceProcessor)
		if !ok {
			return nil, errors.Errorf("expected CreateServiceProcessor, not %T", nopp)
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

func (opp *CreateServiceProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess CreateService")

	fact, ok := op.Fact().(CreateServiceFact)
	if !ok {
		return ctx, nil, e.Errorf("not CreateServiceFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}
	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot create service, %q; %w", fact.Sender(), err), nil
	}

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %q; %w", fact.Contract(), err), nil
	}

	if !ca.Owner().Equal(fact.sender) {
		return nil, base.NewBaseOperationProcessReasonError("not contract account owner, %q", fact.sender), nil
	}

	if ca.IsActive() {
		return nil, base.NewBaseOperationProcessReasonError("a design is already registered, %q", fact.Contract().String()), nil
	}

	if err := currencystate.CheckNotExistsState(state.StateKeyDesign(fact.Contract(), fact.ServiceID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service already exists, %s-%s; %w", fact.Contract(), fact.ServiceID(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *CreateServiceProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process CreateService")

	fact, ok := op.Fact().(CreateServiceFact)
	if !ok {
		return nil, nil, e.Errorf("expected CreateServiceFact, not %T", op.Fact())
	}

	var templates []string
	var holders []crendentialtypes.Holder

	policy := crendentialtypes.NewPolicy(templates, holders, 0)
	if err := policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid credential policy, %s-%s; %w", fact.Contract(), fact.ServiceID(), err), nil
	}

	design := crendentialtypes.NewDesign(fact.ServiceID(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid credential design, %s-%s; %w", fact.Contract(), fact.ServiceID(), err), nil
	}

	sts := make([]base.StateMergeValue, 3)

	sts[0] = state.NewStateMergeValue(
		state.StateKeyDesign(fact.Contract(), fact.ServiceID()),
		state.NewDesignStateValue(design),
	)

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Contract(), err), nil
	}
	ca.SetIsActive(true)

	sts[1] = currencystate.NewStateMergeValue(
		extensioncurrency.StateKeyContractAccount(fact.Contract()),
		extensioncurrency.NewContractAccountStateValue(ca),
	)

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q; %w", fact.Currency(), err), nil
	}

	st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q; %w", fact.Sender(), err), nil
	}
	sb := state.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q; %w", currency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[2] = state.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))))

	return sts, nil, nil
}

func (opp *CreateServiceProcessor) Close() error {
	createServiceProcessorPool.Put(opp)

	return nil
}
