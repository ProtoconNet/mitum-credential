package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
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

func (it Credential) Bytes() []byte {
	if it.holder == nil {
		return util.ConcatBytesSlice(
			[]byte(it.templateID),
			[]byte(it.id),
			[]byte(it.value),
			util.Uint64ToBytes(it.validFrom),
			util.Uint64ToBytes(it.validUntil),
			[]byte(it.did),
		)
	}

	return util.ConcatBytesSlice(
		it.holder.Bytes(),
		[]byte(it.templateID),
		[]byte(it.id),
		[]byte(it.value),
		util.Uint64ToBytes(it.validFrom),
		util.Uint64ToBytes(it.validUntil),
		[]byte(it.did),
	)
}

func (it Credential) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
	); err != nil {
		return err
	}
	if err := util.CheckIsValiders(nil, true,
		it.holder,
	); err != nil {
		return err
	}

	if it.validUntil <= it.validFrom {
		return util.ErrInvalid.Errorf("valid until <= valid from, %q <= %q", it.validUntil, it.validFrom)
	}

	if len(it.id) == 0 {
		return util.ErrInvalid.Errorf("empty id")
	}

	if len(it.did) == 0 {
		return util.ErrInvalid.Errorf("empty did")
	}

	if len(it.value) == 0 {
		return util.ErrInvalid.Errorf("empty value")
	}

	return nil
}

func (it Credential) Holder() base.Address {
	return it.holder
}

func (it Credential) TemplateID() string {
	return it.templateID
}

func (it Credential) ValidFrom() uint64 {
	return it.validFrom
}

func (it Credential) ValidUntil() uint64 {
	return it.validUntil
}

func (it Credential) ID() string {
	return it.id
}

func (it Credential) Value() string {
	return it.value
}

func (it Credential) DID() string {
	return it.did
}
