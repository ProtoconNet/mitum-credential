package credential

import (
	"context"
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

var issueItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueItemProcessor)
	},
}

var issueProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueProcessor)
	},
}

func (Issue) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type IssueItemProcessor struct {
	h               util.Hash
	sender          base.Address
	item            IssueItem
	credentialCount *uint64
	holders         *[]types.Holder
}

func (ipp *IssueItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess IssueItemProcessor")
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		it.Holder(), "holder", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v: holder %v is contract account", cErr, it.Holder())))
	}

	if st, err := cstate.ExistsState(state.StateKeyDesign(it.Contract()), "design", getStateFunc); err != nil {
		return e.Wrap(
			common.ErrServiceNF.Errorf("credential service state for contract account %v", it.Contract()))
	} else if de, err := state.StateDesignValue(st); err != nil {
		return e.Wrap(
			common.ErrServiceNF.Errorf("credential service state value for contract account %v", it.Contract()))
	} else {
		if err := de.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
		for i, v := range de.Policy().TemplateIDs() {
			if it.templateID == v {
				break
			}
			if i == len(de.Policy().TemplateIDs())-1 {
				return e.Wrap(
					common.ErrValueInvalid.Errorf(
						"templateID %v not registered in contract account %v", it.TemplateID(), it.Contract()))
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyCredential(it.Contract(),
		it.TemplateID(),
		it.CredentialID())); {
	case err != nil:
		return e.Wrap(common.ErrStateNF.Errorf(
			"credential %v for template id %v in contract account %v", it.CredentialID(), it.TemplateID(), it.Contract()))
	case !found:
	default:
		if credential, isActive, err := state.StateCredentialValue(st); err != nil {
			return e.Wrap(
				common.ErrStateValInvalid.Errorf(
					"credential %v for template id %v in contract account %v",
					it.CredentialID(), it.TemplateID(), it.Contract()))
		} else if isActive {
			return e.Wrap(
				common.ErrValueInvalid.Errorf(
					"credential %v for template %v is already issued to holder %v in contract account %v",
					it.CredentialID(), it.TemplateID(), credential.Holder(), it.Contract()))
		}
	}

	return nil
}

func (ipp *IssueItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	*ipp.credentialCount++

	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(it.Holder(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	credential := types.NewCredential(it.Holder(), it.TemplateID(), it.CredentialID(), it.Value(), it.ValidFrom(), it.ValidUntil(), it.DID())
	if err := credential.IsValid(nil); err != nil {
		return nil, err
	}

	sts = append(sts, cstate.NewStateMergeValue(
		state.StateKeyCredential(it.Contract(), it.TemplateID(), it.CredentialID()),
		state.NewCredentialStateValue(credential, true),
	))

	sts = append(sts, cstate.NewStateMergeValue(
		state.StateKeyHolderDID(it.Contract(), it.Holder()),
		state.NewHolderDIDStateValue(it.DID()),
	))

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

func (ipp *IssueItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = IssueItem{}
	ipp.credentialCount = nil
	ipp.holders = nil

	issueItemProcessorPool.Put(ipp)
}

type IssueProcessor struct {
	*base.BaseOperationProcessor
}

func NewIssueProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new IssueProcessor")

		nopp := issueProcessorPool.Get()
		opp, ok := nopp.(*IssueProcessor)
		if !ok {
			return nil, e.Errorf("expected IssueProcessor, not %T", nopp)
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

func (opp *IssueProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(IssueFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", IssueFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	for _, it := range fact.Items() {
		ip := issueItemProcessorPool.Get()
		ipc, ok := ip.(*IssueItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected IssueItemProcessor, not %T", ip)), nil
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

func (opp *IssueProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(IssueFact)
	designs := map[string]types.Design{}
	counters := map[string]*uint64{}
	holders := map[string]*[]types.Holder{}

	for _, it := range fact.Items() {
		k := state.StateKeyDesign(it.Contract())

		if _, found := counters[k]; found {
			continue
		}

		st, _ := cstate.ExistsState(k, "design", getStateFunc)

		design, _ := state.StateDesignValue(st)
		count := design.Policy().CredentialCount()
		holder := design.Policy().Holders()

		designs[k] = design
		counters[k] = &count
		holders[k] = &holder
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := issueItemProcessorPool.Get()
		ipc, _ := ip.(*IssueItemProcessor)

		k := state.StateKeyDesign(it.Contract())

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.credentialCount = counters[k]
		ipc.holders = holders[k]

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process IssueItem; %w", err), nil
		}

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

	return sts, nil, nil
}

func (opp *IssueProcessor) Close() error {
	issueProcessorPool.Put(opp)

	return nil
}
