package credential

import (
	"context"
	"sort"
	"sync"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
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

func NewAddTemplateProcessor() ctypes.GetNewProcessor {
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
	fact, ok := op.Fact().(AddTemplateFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", AddTemplateFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if _, _, _, cErr := cstate.ExistsCAccount(fact.Creator(), "creator", true, false, getStateFunc); cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: creator %v is contract account", cErr, fact.Contract())), nil
	}

	st, err := cstate.ExistsState(state.StateKeyDesign(fact.Contract()), "design", getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).
				Errorf("credential design state for contract account %v", fact.Contract())), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Wrap(common.ErrMServiceNF).
				Errorf("credential design state value for contract account %v", fact.Contract())), nil
	}

	for _, templateID := range design.Policy().TemplateIDs() {
		if templateID == fact.TemplateID() {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Wrap(common.ErrMStateE).
					Errorf("already registered template %v in contract account %v", fact.TemplateID(), fact.Contract())), nil
		}
	}

	return ctx, nil, nil
}

func (opp *AddTemplateProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(AddTemplateFact)
	st, _ := cstate.ExistsState(state.StateKeyDesign(fact.Contract()), "design", getStateFunc)
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

	smv, err := cstate.CreateNotExistAccount(fact.creator, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("%w", err), nil
	} else if smv != nil {
		sts = append(sts, smv)
	}

	sts = append(sts, cstate.NewStateMergeValue(
		state.StateKeyDesign(fact.Contract()),
		state.NewDesignStateValue(design),
	))

	sts = append(sts, cstate.NewStateMergeValue(
		state.StateKeyTemplate(fact.Contract(), fact.TemplateID()),
		state.NewTemplateStateValue(template),
	))

	return sts, nil, nil
}

func (opp *AddTemplateProcessor) Close() error {
	addTemplateProcessorPool.Put(opp)

	return nil
}
