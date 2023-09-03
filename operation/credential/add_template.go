package credential

import (
	"unicode/utf8"

	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
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
	sender         base.Address
	contract       base.Address
	serviceID      types.ServiceID
	templateID     string
	templateName   string
	serviceDate    types.Date
	expirationDate types.Date
	templateShare  types.Bool
	multiAudit     types.Bool
	displayName    string
	subjectKey     string
	description    string
	creator        base.Address
	currency       currencytypes.CurrencyID
}

func NewAddTemplateFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	serviceID types.ServiceID,
	templateID string,
	templateName string,
	serviceDate types.Date,
	expirationDate types.Date,
	templateShare types.Bool,
	multiAudit types.Bool,
	displayName string,
	subjectKey string,
	description string,
	creator base.Address,
	currency currencytypes.CurrencyID,
) AddTemplateFact {
	bf := base.NewBaseFact(AddTemplateFactHint, token)
	fact := AddTemplateFact{
		BaseFact:       bf,
		sender:         sender,
		contract:       contract,
		serviceID:      serviceID,
		templateID:     templateID,
		templateName:   templateName,
		serviceDate:    serviceDate,
		expirationDate: expirationDate,
		templateShare:  templateShare,
		multiAudit:     multiAudit,
		displayName:    displayName,
		subjectKey:     subjectKey,
		description:    description,
		creator:        creator,
		currency:       currency,
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
		fact.serviceID.Bytes(),
		[]byte(fact.templateID),
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
	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.serviceID,
		fact.serviceDate,
		fact.expirationDate,
		fact.currency,
	); err != nil {
		return err
	}

	if l := utf8.RuneCountInString(fact.templateID); l < 1 || l > MaxLengthTemplateID {
		return util.ErrInvalid.Errorf("invalid length of template ID, 0 <= length <= %d", MaxLengthTemplateID)
	}

	if l := utf8.RuneCountInString(fact.templateName); l < 1 || l > MaxLengthTemplateName {
		return util.ErrInvalid.Errorf("invalid length of template name, 0 <= length <= %d", MaxLengthTemplateName)
	}

	if l := utf8.RuneCountInString(fact.displayName); l < 1 || l > MaxLengthDisplayName {
		return util.ErrInvalid.Errorf("invalid length of display name, 0 <= length <= %d", MaxLengthDisplayName)
	}

	if l := utf8.RuneCountInString(fact.subjectKey); l < 1 || l > MaxLengthSubjectKey {
		return util.ErrInvalid.Errorf("invalid length of subjectKey, 0 <= length <= %d", MaxLengthSubjectKey)
	}

	if l := utf8.RuneCountInString(fact.description); l < 1 || l > MaxLengthDescription {
		return util.ErrInvalid.Errorf("invalid length of description, 0 <= length <= %d", MaxLengthDescription)
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

func (fact AddTemplateFact) ServiceID() types.ServiceID {
	return fact.serviceID
}

func (fact AddTemplateFact) TemplateID() string {
	return fact.templateID
}

func (fact AddTemplateFact) TemplateName() string {
	return fact.templateName
}

func (fact AddTemplateFact) ServiceDate() types.Date {
	return fact.serviceDate
}

func (fact AddTemplateFact) ExpirationDate() types.Date {
	return fact.expirationDate
}

func (fact AddTemplateFact) TemplateShare() types.Bool {
	return fact.templateShare
}

func (fact AddTemplateFact) MultiAudit() types.Bool {
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

func (fact AddTemplateFact) Currency() currencytypes.CurrencyID {
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
	common.BaseOperation
}

func NewAddTemplate(fact AddTemplateFact) (AddTemplate, error) {
	return AddTemplate{BaseOperation: common.NewBaseOperation(AddTemplateHint, fact)}, nil
}

func (op *AddTemplate) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}

var (
	MaxLengthTemplateID      = 20
	MaxLengthCredentialID    = 20
	MaxLengthTemplateName    = 20
	MaxLengthDisplayName     = 20
	MaxLengthSubjectKey      = 256
	MaxLengthCredentialValue = 1024
	MaxLengthDescription     = 1024
)
