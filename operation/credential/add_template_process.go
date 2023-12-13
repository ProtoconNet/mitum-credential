package credential

import (
	"context"
	"sort"
	"sync"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
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
	_ context.Context, _ base.GetStateFunc,
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
		e := util.StringError("failed to create new AddTemplateProcessor")

		nopp := addTemplateProcessorPool.Get()
		opp, ok := nopp.(*AddTemplateProcessor)
		if !ok {
			return nil, errors.Errorf("expected AddTemplateProcessor, not %T", nopp)
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

func (opp *AddTemplateProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess AddTemplate")

	fact, ok := op.Fact().(AddTemplateFact)
	if !ok {
		return ctx, nil, e.Errorf("not %T, %T", AddTemplateFact{}, op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender account state not found, %q; %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return ctx, nil, e.WithMessage(err, "fee Currency state not found")
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender account is contract account, %q; %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Creator()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("creator account state not found, %q; %w", fact.Creator(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Creator()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("creator account is contract account, %q; %w", fact.Creator(), err), nil
	}

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target contract account state not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found from state, %q; %w", fact.Contract(), err), nil
	}

	if !(ca.Owner().Equal(fact.sender) || ca.IsOperator(fact.Sender())) {
		return nil, base.NewBaseOperationProcessReasonError("sender account is neither the owner nor the operator of the target contract account, %q", fact.sender), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyDesign(fact.Contract()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service state not found, %s; %w", fact.Contract(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("credential service value not found from state, %s; %w", fact.Contract(), err), nil
	}

	for _, templateID := range design.Policy().TemplateIDs() {
		if templateID == fact.TemplateID() {
			return nil, base.NewBaseOperationProcessReasonError("already registered template, %q, %s", fact.TemplateID(), fact.Contract()), nil
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *AddTemplateProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(AddTemplateFact)
	st, _ := currencystate.ExistsState(state.StateKeyDesign(fact.Contract()), "key of design", getStateFunc)
	design, _ := state.StateDesignValue(st)
	templateIDs := design.Policy().TemplateIDs()
	templateIDs = append(templateIDs, fact.templateID)
	sort.Slice(templateIDs, func(i int, j int) bool {
		return templateIDs[i] < templateIDs[j]
	})
	policy := types.NewPolicy(templateIDs, design.Policy().Holders(), design.Policy().CredentialCount())
	if err := policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid credential policy, %s; %w", fact.Contract(), err), nil
	}

	design = types.NewDesign(policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid credential design, %s; %w", fact.Contract(), err), nil
	}

	template := types.NewTemplate(
		fact.TemplateID(), fact.TemplateName(), fact.ServiceDate(), fact.ExpirationDate(),
		fact.TemplateShare(), fact.MultiAudit(), fact.DisplayName(), fact.SubjectKey(),
		fact.Description(), fact.Creator(),
	)
	if err := template.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid template, %q; %w", fact.TemplateID(), err), nil
	}

	var sts []base.StateMergeValue

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyDesign(fact.Contract()),
		state.NewDesignStateValue(design),
	))

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyTemplate(fact.Contract(), fact.TemplateID()),
		state.NewTemplateStateValue(template),
	))

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	if currencyPolicy.Feeer().Receiver() == nil {
		return sts, nil, nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %q; %w",
			fact.Currency(),
			err,
		), nil
	}

	senderBalSt, err := currencystate.ExistsState(
		currency.StateKeyBalance(fact.Sender(), fact.Currency()),
		"key of sender balance",
		getStateFunc,
	)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"sender balance not found, %q; %w",
			fact.Sender(),
			err,
		), nil
	}

	switch senderBal, err := currency.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			currency.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := senderBalSt.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
		return nil, nil, err
	} else if feeRcvrSt, found, err := getStateFunc(currency.StateKeyBalance(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
		return nil, nil, err
	} else if !found {
		return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
	} else if feeRcvrSt.Key() != senderBalSt.Key() {
		r, ok := feeRcvrSt.Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, errors.Errorf("expected %T, not %T", currency.BalanceStateValue{}, feeRcvrSt.Value())
		}
		sts = append(sts, common.NewBaseStateMergeValue(
			feeRcvrSt.Key(),
			currency.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
			},
		))

		sts = append(sts, common.NewBaseStateMergeValue(
			senderBalSt.Key(),
			currency.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
			},
		))
	}
	
	return sts, nil, nil
}

func (opp *AddTemplateProcessor) Close() error {
	addTemplateProcessorPool.Put(opp)

	return nil
}
