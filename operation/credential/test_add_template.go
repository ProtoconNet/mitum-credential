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
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TestAddTemplateProcessor struct {
	*test.BaseTestOperationProcessorNoItem[AddTemplate]
	templateID     string
	templateName   string
	serviceDate    credentialtypes.Date
	expirationDate credentialtypes.Date
	templateShare  credentialtypes.Bool
	multiAudit     credentialtypes.Bool
	displayName    string
	subjectKey     string
	description    string
}

func NewTestAddTemplateProcessor(encs *encoder.Encoders) TestAddTemplateProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[AddTemplate](encs)
	return TestAddTemplateProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestAddTemplateProcessor) Create() *TestAddTemplateProcessor {
	t.Opr, _ = NewAddTemplateProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestAddTemplateProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestAddTemplateProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestAddTemplateProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestAddTemplateProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestAddTemplateProcessor) SetService(
	contract base.Address,
) *TestAddTemplateProcessor {
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

func (t *TestAddTemplateProcessor) LoadOperation(fileName string,
) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestAddTemplateProcessor) Print(fileName string,
) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestAddTemplateProcessor) SetTemplate(
	templateID, templateName string, serviceDate, expirationDate credentialtypes.Date,
	templateShare, multiAudit credentialtypes.Bool, displayName, subjectKey, description string,
) *TestAddTemplateProcessor {
	t.templateID = templateID
	t.templateName = templateName
	t.serviceDate = serviceDate
	t.expirationDate = expirationDate
	t.templateShare = templateShare
	t.multiAudit = multiAudit
	t.displayName = displayName
	t.subjectKey = subjectKey
	t.description = description

	return t
}

func (t *TestAddTemplateProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract, creator base.Address, currency types.CurrencyID,
) *TestAddTemplateProcessor {
	op := NewAddTemplate(
		NewAddTemplateFact(
			[]byte("token"),
			sender,
			contract,
			t.templateID,
			t.templateName,
			t.serviceDate,
			t.expirationDate,
			t.templateShare,
			t.multiAudit,
			t.displayName,
			t.subjectKey,
			t.description,
			creator,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestAddTemplateProcessor) RunPreProcess() *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestAddTemplateProcessor) RunProcess() *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestAddTemplateProcessor) IsValid() *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestAddTemplateProcessor) Decode(fileName string) *TestAddTemplateProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
