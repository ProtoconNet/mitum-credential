package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	AddTemplateFactHint = hint.MustNewHint("mitum-credential-add-template-operation-fact-v0.0.1")
	AddTemplateHint     = hint.MustNewHint("mitum-credential-add-template-operation-v0.0.1")
)

type AddTemplateFact struct {
	base.BaseFact
	sender              base.Address
	contract            base.Address
	credentialServiceID currencybase.ContractID
	templateID          Uint256
	templateName        string
	serviceDate         Date
	expirationDate      Date
	templateShare       Bool
	multiAudit          Bool
	displayName         string
	subjectKey          string
	description         string
	creator             base.Address
	currency            currencybase.CurrencyID
}

func NewAddTemplateFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	credentialServiceID currencybase.ContractID,
	templateID Uint256,
	templateName string,
	serviceDate Date,
	expirationDate Date,
	templateShare Bool,
	multiAudit Bool,
	displayName string,
	subjectKey string,
	description string,
	creator base.Address,
	currency currencybase.CurrencyID,
) AddTemplateFact {
	bf := base.NewBaseFact(AddTemplateFactHint, token)
	fact := AddTemplateFact{
		BaseFact:            bf,
		sender:              sender,
		contract:            contract,
		credentialServiceID: credentialServiceID,
		templateID:          templateID,
		templateName:        templateName,
		serviceDate:         serviceDate,
		expirationDate:      expirationDate,
		templateShare:       templateShare,
		multiAudit:          multiAudit,
		displayName:         displayName,
		subjectKey:          subjectKey,
		description:         description,
		creator:             creator,
		currency:            currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AddTemplateFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AddTemplateFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddTemplateFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.credentialServiceID.Bytes(),
		fact.templateID.Bytes(),
		[]byte(fact.templateName),
		fact.serviceDate.Bytes(),
		fact.expirationDate.Bytes(),
		fact.templateShare.Bytes(),
		fact.multiAudit.Bytes(),
		[]byte(fact.displayName),
		[]byte(fact.subjectKey),
		[]byte(fact.description),
		fact.creator.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact AddTemplateFact) IsValid(b []byte) error {
	if err := currencybase.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.credentialServiceID,
		fact.templateID,
		fact.serviceDate,
		fact.expirationDate,
		fact.currency,
	); err != nil {
		return err
	}

	if len(fact.templateName) == 0 {
		return util.ErrInvalid.Errorf("empty template name")
	}

	if len(fact.displayName) == 0 {
		return util.ErrInvalid.Errorf("empty display name")
	}

	if len(fact.subjectKey) == 0 {
		return util.ErrInvalid.Errorf("empty subject key")
	}

	if len(fact.description) == 0 {
		return util.ErrInvalid.Errorf("empty description")
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if fact.creator.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with creator, %q", fact.creator)
	}

	service, err := fact.serviceDate.Parse()
	if err != nil {
		return err
	}

	expire, err := fact.serviceDate.Parse()
	if err != nil {
		return err
	}

	if expire.UnixNano() < service.UnixNano() {
		return util.ErrInvalid.Errorf("expire date <= service date, %s <= %s", fact.expirationDate, fact.serviceDate)
	}

	return nil
}

func (fact AddTemplateFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AddTemplateFact) Sender() base.Address {
	return fact.sender
}

func (fact AddTemplateFact) Contract() base.Address {
	return fact.contract
}

func (fact AddTemplateFact) CredentialServiceID() currencybase.ContractID {
	return fact.credentialServiceID
}

func (fact AddTemplateFact) TemplateID() Uint256 {
	return fact.templateID
}

func (fact AddTemplateFact) TemplateName() string {
	return fact.templateName
}

func (fact AddTemplateFact) ServiceDate() Date {
	return fact.serviceDate
}

func (fact AddTemplateFact) ExpirationDate() Date {
	return fact.expirationDate
}

func (fact AddTemplateFact) TemplateShare() Bool {
	return fact.templateShare
}

func (fact AddTemplateFact) MultiAudit() Bool {
	return fact.multiAudit
}

func (fact AddTemplateFact) DisplayName() string {
	return fact.displayName
}

func (fact AddTemplateFact) SubjectKey() string {
	return fact.subjectKey
}

func (fact AddTemplateFact) Description() string {
	return fact.description
}

func (fact AddTemplateFact) Creator() base.Address {
	return fact.creator
}

func (fact AddTemplateFact) Currency() currencybase.CurrencyID {
	return fact.currency
}

func (fact AddTemplateFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 3)
	as[0] = fact.sender
	as[1] = fact.contract
	as[2] = fact.creator
	return as, nil
}

type AddTemplate struct {
	currencybase.BaseOperation
}

func NewAddTemplate(fact AddTemplateFact) (AddTemplate, error) {
	return AddTemplate{BaseOperation: currencybase.NewBaseOperation(AddTemplateHint, fact)}, nil
}

func (op *AddTemplate) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
