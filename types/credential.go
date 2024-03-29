package types

import (
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"unicode/utf8"
)

var CredentialHint = hint.MustNewHint("mitum-credential-credential-v0.0.1")

type Credential struct {
	hint.BaseHinter
	holder     base.Address
	templateID string
	id         string
	value      string
	validFrom  uint64
	validUntil uint64
	did        string
}

func NewCredential(
	holder base.Address,
	templateID string,
	id string,
	value string,
	validFrom uint64,
	validUntil uint64,
	did string,
) Credential {
	return Credential{
		BaseHinter: hint.NewBaseHinter(CredentialHint),
		holder:     holder,
		templateID: templateID,
		id:         id,
		value:      value,
		validFrom:  validFrom,
		validUntil: validUntil,
		did:        did,
	}
}

func (c Credential) Bytes() []byte {
	if c.holder == nil {
		return util.ConcatBytesSlice(
			[]byte(c.templateID),
			[]byte(c.id),
			[]byte(c.value),
			util.Uint64ToBytes(c.validFrom),
			util.Uint64ToBytes(c.validUntil),
			[]byte(c.did),
		)
	}

	return util.ConcatBytesSlice(
		c.holder.Bytes(),
		[]byte(c.templateID),
		[]byte(c.id),
		[]byte(c.value),
		util.Uint64ToBytes(c.validFrom),
		util.Uint64ToBytes(c.validUntil),
		[]byte(c.did),
	)
}

func (c Credential) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		c.BaseHinter,
	); err != nil {
		return err
	}
	if err := util.CheckIsValiders(nil, true,
		c.holder,
	); err != nil {
		return err
	}

	if c.validUntil <= c.validFrom {
		return util.ErrInvalid.Errorf("valid until <= valid from, %q <= %q", c.validUntil, c.validFrom)
	}

	if l := utf8.RuneCountInString(c.templateID); l < 1 || l > MaxLengthTemplateID {
		return util.ErrInvalid.Errorf("invalid length of credential ID, 0 <= length <= %d", MaxLengthTemplateID)
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(c.templateID)) {
		return util.ErrInvalid.Errorf("invalid template ID due to the inclusion of special characters")
	}

	if l := utf8.RuneCountInString(c.id); l < 1 || l > MaxLengthCredentialID {
		return util.ErrInvalid.Errorf("invalid length of credential ID, 0 <= length <= %d", MaxLengthCredentialID)
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(c.id)) {
		return util.ErrInvalid.Errorf("invalid credential ID due to the inclusion of special characters")
	}

	if len(c.did) == 0 {
		return util.ErrInvalid.Errorf("empty did")
	}

	if len(c.value) == 0 {
		return util.ErrInvalid.Errorf("empty value")
	}

	return nil
}

func (c Credential) Holder() base.Address {
	return c.holder
}

func (c Credential) TemplateID() string {
	return c.templateID
}

func (c Credential) ValidFrom() uint64 {
	return c.validFrom
}

func (c Credential) ValidUntil() uint64 {
	return c.validUntil
}

func (c Credential) ID() string {
	return c.id
}

func (c Credential) Value() string {
	return c.value
}

func (c Credential) DID() string {
	return c.did
}
