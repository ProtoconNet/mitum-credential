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

type TestIssueProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Issue, IssueItem]
	templateID string
	id         string
	value      string
	validFrom  uint64
	validUntil uint64
	did        string
}

func NewTestIssueProcessor(tp *test.TestProcessor) TestIssueProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Issue, IssueItem](tp)
	return TestIssueProcessor{BaseTestOperationProcessorWithItem: &t}
}

func (t *TestIssueProcessor) Create() *TestIssueProcessor {
	t.Opr, _ = NewIssueProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestIssueProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestIssueProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestIssueProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestIssueProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestIssueProcessor) SetService(
	contract base.Address,
) *TestIssueProcessor {
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

func (t *TestIssueProcessor) LoadOperation(fileName string,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestIssueProcessor) Print(fileName string,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestIssueProcessor) SetTemplate(
	templateID,
	id,
	value string,
	validFrom,
	validUntil uint64,
	did string,
) *TestIssueProcessor {
	t.templateID = templateID
	t.id = id
	t.value = value
	t.validFrom = validFrom
	t.validUntil = validUntil
	t.did = did

	return t
}

func (t *TestIssueProcessor) MakeItem(
	contract, holder test.Account, currency types.CurrencyID, targetItems []IssueItem,
) *TestIssueProcessor {
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

func (t *TestIssueProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []IssueItem,
) *TestIssueProcessor {
	op := NewIssue(NewIssueFact(
		[]byte("token"),
		sender,
		items,
	))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestIssueProcessor) RunPreProcess() *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestIssueProcessor) RunProcess() *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestIssueProcessor) IsValid() *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestIssueProcessor) Decode(fileName string) *TestIssueProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
