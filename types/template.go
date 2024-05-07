package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"unicode/utf8"
)

var TemplateHint = hint.MustNewHint("mitum-credential-template-v0.0.1")

var (
	MaxLengthTemplateID      = 20
	MaxLengthCredentialID    = 20
	MaxLengthTemplateName    = 20
	MaxLengthDisplayName     = 20
	MaxLengthSubjectKey      = 256
	MaxLengthCredentialValue = 1024
	MaxLengthDescription     = 1024
)

type Template struct {
	hint.BaseHinter
	templateID     string
	templateName   string
	serviceDate    Date
	expirationDate Date
	templateShare  Bool
	multiAudit     Bool
	displayName    string
	subjectKey     string
	description    string
	creator        base.Address
}

func NewTemplate(
	templateID,
	templateName string,
	serviceDate,
	expirationDate Date,
	templateShare,
	multiAudit Bool,
	displayName,
	subjectKey,
	description string,
	creator base.Address,
) Template {
	return Template{
		BaseHinter:     hint.NewBaseHinter(TemplateHint),
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
	}
}

func (t Template) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		t.BaseHinter,
		t.serviceDate,
		t.expirationDate,
		t.creator,
	); err != nil {
		return err
	}

	if l := utf8.RuneCountInString(t.templateID); l < 1 || l > MaxLengthTemplateID {
		return common.ErrValOOR.Errorf("0 <= length of credential ID <= %d", MaxLengthTemplateID)
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(t.templateID)) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("template ID %s, must match regex `^[^\\s:/?#\\[\\]@]*$`", t.templateID))
	}

	if len(t.templateName) == 0 {
		return common.ErrValOOR.Errorf("empty template name")
	}

	if len(t.displayName) == 0 {
		return common.ErrValOOR.Errorf("empty display name")
	}

	if len(t.subjectKey) == 0 {
		return common.ErrValOOR.Errorf("empty subject key")
	}

	serviceDate, err := t.serviceDate.Parse()
	if err != nil {
		return common.ErrValueInvalid.Wrap(err)
	}

	expireDate, err := t.expirationDate.Parse()
	if err != nil {
		return common.ErrValueInvalid.Wrap(err)
	}

	if expireDate.UnixNano() < serviceDate.UnixNano() {
		return common.ErrValOOR.Errorf("expire date <= service date, but %s <= %s", t.expirationDate, t.serviceDate)
	}

	return nil
}

func (t Template) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(t.templateID),
		[]byte(t.templateName),
		t.serviceDate.Bytes(),
		t.expirationDate.Bytes(),
		t.templateShare.Bytes(),
		t.multiAudit.Bytes(),
		[]byte(t.displayName),
		[]byte(t.subjectKey),
		[]byte(t.description),
		t.creator.Bytes(),
	)
}

func (t Template) TemplateID() string {
	return t.templateID
}

func (t Template) TemplateName() string {
	return t.templateName
}

func (t Template) ServiceDate() Date {
	return t.serviceDate
}

func (t Template) ExpirationDate() Date {
	return t.expirationDate
}

func (t Template) TemplateShare() Bool {
	return t.templateShare
}

func (t Template) MultiAudit() Bool {
	return t.multiAudit
}

func (t Template) DisplayName() string {
	return t.displayName
}

func (t Template) SubjectKey() string {
	return t.subjectKey
}

func (t Template) Description() string {
	return t.description
}

func (t Template) Creator() base.Address {
	return t.creator
}
