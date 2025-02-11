package credential

import (
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type TestRevokeProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Revoke, RevokeItem]
	templateID string
	id         string
}

func NewTestRevokeProcessor(tp *test.TestProcessor) TestRevokeProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Revoke, RevokeItem](tp)
	return TestRevokeProcessor{BaseTestOperationProcessorWithItem: &t}
}

func (t *TestRevokeProcessor) Create() *TestRevokeProcessor {
	t.Opr, _ = NewRevokeProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestRevokeProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestRevokeProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestRevokeProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestRevokeProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestRevokeProcessor) LoadOperation(fileName string,
) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestRevokeProcessor) Print(fileName string,
) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestRevokeProcessor) SetTemplate(
	templateID,
	id string,
) *TestRevokeProcessor {
	t.templateID = templateID
	t.id = id

	return t
}

func (t *TestRevokeProcessor) SetService(
	contract base.Address,
) *TestRevokeProcessor {
	var templates []string
	var holders []types.Holder

	policy := types.NewPolicy(templates, holders, 0)
	design := types.NewDesign(policy)

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

	nstatus := status.SetActive(true)
	cState := common.NewBaseState(base.Height(1), extension.StateKeyContractAccount(contract), extension.NewContractAccountStateValue(nstatus), nil, []util.Hash{})
	t.SetState(cState, true)

	return t
}

func (t *TestRevokeProcessor) MakeItem(
	contract, holder test.Account, currency ctypes.CurrencyID, targetItems []RevokeItem,
) *TestRevokeProcessor {
	item := NewRevokeItem(
		contract.Address(),
		holder.Address(),
		t.templateID,
		t.id,
		currency,
	)
	test.UpdateSlice[RevokeItem](item, targetItems)

	return t
}

func (t *TestRevokeProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []RevokeItem,
) *TestRevokeProcessor {
	op := NewRevoke(
		NewRevokeFact(
			[]byte("token"),
			sender,
			items,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestRevokeProcessor) RunPreProcess() *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestRevokeProcessor) RunProcess() *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestRevokeProcessor) IsValid() *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestRevokeProcessor) Decode(fileName string) *TestRevokeProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
