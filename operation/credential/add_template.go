package credential

import (
	"unicode/utf8"

	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	AddTemplateFactHint = hint.MustNewHint("mitum-credential-add-template-operation-fact-v0.0.1")
	AddTemplateHint     = hint.MustNewHint("mitum-credential-add-template-operation-v0.0.1")
)

type AddTemplateFact struct {
	base.BaseFact
	sender         base.Address
	contract       base.Address
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
	currency       ctypes.CurrencyID
}

func NewAddTemplateFact(
	token []byte,
	sender base.Address,
	contract base.Address,
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
	currency ctypes.CurrencyID,
) AddTemplateFact {
	bf := base.NewBaseFact(AddTemplateFactHint, token)
	fact := AddTemplateFact{
		BaseFact:       bf,
		sender:         sender,
		contract:       contract,
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
	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.serviceDate,
		fact.expirationDate,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := utf8.RuneCountInString(fact.templateID); l < 1 || l > types.MaxLengthTemplateID {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of template ID <= %d, but %d", types.MaxLengthTemplateID, l)))
	}

	if !ctypes.ReValidSpcecialCh.Match([]byte(fact.templateID)) {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("template ID %s, must match regex `^[^\\s:/?#\\[\\]$@]*$`", fact.TemplateID())))
	}

	if l := utf8.RuneCountInString(fact.templateName); l < 1 || l > types.MaxLengthTemplateName {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of template name <= %d, but %d", types.MaxLengthTemplateName, l)))
	}

	if l := utf8.RuneCountInString(fact.displayName); l < 1 || l > types.MaxLengthDisplayName {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of display name <= %d, but %d", types.MaxLengthDisplayName, l)))
	}

	if l := utf8.RuneCountInString(fact.subjectKey); l < 1 || l > types.MaxLengthSubjectKey {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of subjectKey <= %d, but %d", types.MaxLengthSubjectKey, l)))
	}

	if l := utf8.RuneCountInString(fact.description); l < 1 || l > types.MaxLengthDescription {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of description <= %d, but %d", types.MaxLengthDescription, l)))
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if fact.creator.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("creator %v is same with contract account", fact.creator)))
	}

	serviceDate, err := fact.serviceDate.Parse()
	if err != nil {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(err))
	}

	expire, err := fact.serviceDate.Parse()
	if err != nil {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(err))
	}

	if expire.UnixNano() < serviceDate.UnixNano() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("expire date <= service date, %s <= %s", fact.expirationDate, fact.serviceDate)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
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

func (fact AddTemplateFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact AddTemplateFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 3)
	as[0] = fact.sender
	as[1] = fact.contract
	as[2] = fact.creator
	return as, nil
}

func (fact AddTemplateFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact AddTemplateFact) FeePayer() base.Address {
	return fact.sender
}

func (fact AddTemplateFact) FactUser() base.Address {
	return fact.sender
}

func (fact AddTemplateFact) Signer() base.Address {
	return fact.sender
}

func (fact AddTemplateFact) ActiveContractOwnerHandlerOnly() [][2]base.Address {
	return [][2]base.Address{{fact.contract, fact.sender}}
}

type AddTemplate struct {
	extras.ExtendedOperation
}

func NewAddTemplate(fact AddTemplateFact) AddTemplate {
	return AddTemplate{
		ExtendedOperation: extras.NewExtendedOperation(AddTemplateHint, fact),
	}
}
