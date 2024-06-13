package credential

import (
	"github.com/ProtoconNet/mitum-credential/state"
	credentialtypes "github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type TestAssignProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Issue, IssueItem]
	templateID string
	id         string
	value      string
	validFrom  uint64
	validUntil uint64
	did        string
}

func NewTestAssignProcessor(tp *test.TestProcessor) TestAssignProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Issue, IssueItem](tp)
	return TestAssignProcessor{BaseTestOperationProcessorWithItem: &t}
}

func (t *TestAssignProcessor) Create() *TestAssignProcessor {
	t.Opr, _ = NewIssueProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestAssignProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestAssignProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestAssignProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestAssignProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestAssignProcessor) SetService(
	contract base.Address,
) *TestAssignProcessor {
	var templates []string
	var holders []credentialtypes.Holder

	policy := credentialtypes.NewPolicy(templates, holders, 0)
	design := credentialtypes.NewDesign(policy)

	st := common.NewBaseState(base.Height(1), state.StateKeyDesign(contract), state.NewDesignStateValue(design), nil, []util.Hash{})
	t.SetState(st, true)

	cst, found, _ := t.MockGetter.Get(extension.StateKeyContractAccount(contract))
	if !found {
		panic("contract account not set")
	}
	status, err := extension.StateContractAccountValue(cst)
	if err != nil {
		panic(err)
	}

	nstatus := status.SetIsActive(true)
	cState := common.NewBaseState(base.Height(1), extension.StateKeyContractAccount(contract), extension.NewContractAccountStateValue(nstatus), nil, []util.Hash{})
	t.SetState(cState, true)

	return t
}

func (t *TestAssignProcessor) LoadOperation(fileName string,
) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestAssignProcessor) Print(fileName string,
) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestAssignProcessor) SetTemplate(
	templateID,
	id,
	value string,
	validFrom,
	validUntil uint64,
	did string,
) *TestAssignProcessor {
	t.templateID = templateID
	t.id = id
	t.value = value
	t.validFrom = validFrom
	t.validUntil = validUntil
	t.did = did

	return t
}

func (t *TestAssignProcessor) MakeItem(
	contract, holder test.Account, currency types.CurrencyID, targetItems []IssueItem,
) *TestAssignProcessor {
	item := NewIssueItem(
		contract.Address(),
		holder.Address(),
		t.templateID,
		t.id,
		t.value,
		t.validFrom,
		t.validUntil,
		t.did,
		currency,
	)
	test.UpdateSlice[IssueItem](item, targetItems)

	return t
}

func (t *TestAssignProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []IssueItem,
) *TestAssignProcessor {
	op := NewAssign(
		NewIssueFact(
			[]byte("token"),
			sender,
			items,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestAssignProcessor) RunPreProcess() *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestAssignProcessor) RunProcess() *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestAssignProcessor) IsValid() *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestAssignProcessor) Decode(fileName string) *TestAssignProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
