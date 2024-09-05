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

type TestRegisterModelProcessor struct {
	*test.BaseTestOperationProcessorNoItem[RegisterModel]
}

func NewTestRegisterModelProcessor(tp *test.TestProcessor) TestRegisterModelProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[RegisterModel](tp)
	return TestRegisterModelProcessor{&t}
}

func (t *TestRegisterModelProcessor) Create() *TestRegisterModelProcessor {
	t.Opr, _ = NewRegisterModelProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestRegisterModelProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestRegisterModelProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestRegisterModelProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterModelProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterModelProcessor) SetService(
	contract base.Address,
) *TestRegisterModelProcessor {
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

func (t *TestRegisterModelProcessor) LoadOperation(fileName string,
) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestRegisterModelProcessor) Print(fileName string,
) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestRegisterModelProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, currency types.CurrencyID,
) *TestRegisterModelProcessor {
	op := NewRegisterModel(
		NewRegisterModelFact(
			[]byte("token"),
			sender,
			contract,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestRegisterModelProcessor) RunPreProcess() *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestRegisterModelProcessor) RunProcess() *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestRegisterModelProcessor) IsValid() *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestRegisterModelProcessor) Decode(fileName string) *TestRegisterModelProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
