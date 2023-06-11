package credential

import (
	"context"
	"sync"

	state "github.com/ProtoconNet/mitum-credential/state/credential"
	common "github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"

	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var addTemplateProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddTemplateProcessor)
	},
}

func (AddTemplate) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AddTemplateProcessor struct {
	*base.BaseOperationProcessor
}

func NewAddTemplateProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new AddTemplateProcessor")

		nopp := addTemplateProcessorPool.Get()
		opp, ok := nopp.(*AddTemplateProcessor)
		if !ok {
			return nil, errors.Errorf("expected AddTemplateProcessor, not %T", nopp)
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

func (opp *AddTemplateProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess AddTemplate")

	fact, ok := op.Fact().(AddTemplateFact)
	if !ok {
		return ctx, nil, e(nil, "not AddTemplateFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot add template, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Creator()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("creator not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Creator()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot be creator, %q: %w", fact.Creator(), err), nil
	}

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account not found, %q: %w", fact.Contract(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %q: %w", fact.Contract(), err), nil
	}

	if !ca.Owner().Equal(fact.sender) {
		return nil, base.NewBaseOperationProcessReasonError("not contract account owner, %q", fact.sender), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.CredentialServiceID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service not found, %s-%s: %w", fact.Contract(), fact.CredentialServiceID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service value not found, %s-%s: %w", fact.Contract(), fact.CredentialServiceID(), err), nil
	}

	for _, t := range design.Policy().Templates() {
		ft := fact.TemplateID().N()
		if t.N().Cmp(ft) == 0 {
			return nil, base.NewBaseOperationProcessReasonError("already registered template, %q, %s-%s", fact.TemplateID(), fact.Contract(), fact.CredentialServiceID()), nil
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *AddTemplateProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process AddTemplate")

	fact, ok := op.Fact().(AddTemplateFact)
	if !ok {
		return nil, nil, e(nil, "expected AddTemplateFact, not %T", op.Fact())
	}

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.CredentialServiceID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service not found, %s-%s: %w", fact.Contract(), fact.CredentialServiceID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service value not found, %s-%s: %w", fact.Contract(), fact.CredentialServiceID(), err), nil
	}

	templates := design.Policy().Templates()

	for _, t := range templates {
		ft := fact.TemplateID().N()
		if t.N().Cmp(ft) == 0 {
			return nil, base.NewBaseOperationProcessReasonError("already registered template, %q, %s-%s", fact.TemplateID(), fact.Contract(), fact.CredentialServiceID()), nil
		}
	}

	templates = append(templates, fact.templateID)

	policy := types.NewPolicy(templates, design.Policy().Holders(), design.Policy().CredentialCount())
	if err := policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid credential policy, %s-%s: %w", fact.Contract(), fact.CredentialServiceID(), err), nil
	}

	design = types.NewDesign(fact.CredentialServiceID(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid credential design, %s-%s: %w", fact.Contract(), fact.CredentialServiceID(), err), nil
	}

	template := types.NewTemplate(
		fact.TemplateID(), fact.TemplateName(), fact.ServiceDate(), fact.ExpirationDate(),
		fact.TemplateShare(), fact.MultiAudit(), fact.DisplayName(), fact.SubjectKey(),
		fact.Description(), fact.Creator(),
	)
	if err := template.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid template, %q: %w", fact.TemplateID(), err), nil
	}

	sts := make([]base.StateMergeValue, 3)

	sts[0] = currencystate.NewStateMergeValue(
		state.StateKeyDesign(fact.Contract(), fact.CredentialServiceID()),
		state.NewDesignStateValue(design),
	)

	sts[1] = currencystate.NewStateMergeValue(
		state.StateKeyTemplate(fact.Contract(), fact.CredentialServiceID(), fact.TemplateID()),
		state.NewTemplateStateValue(template),
	)

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}

	sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[2] = currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))))

	return sts, nil, nil
}

func (opp *AddTemplateProcessor) Close() error {
	addTemplateProcessorPool.Put(opp)

	return nil
}