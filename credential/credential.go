package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var CredentialHint = hint.MustNewHint("mitum-credential-credential-v0.0.1")

type Credential struct {
	hint.BaseHinter
	holder     base.Address
	templateID Uint256
	id         string
	value      string
	validfrom  Uint256
	validuntil Uint256
	did        string
}

func NewCredential(
	holder base.Address,
	templateID Uint256,
	id string,
	value string,
	validfrom Uint256,
	validuntil Uint256,
	did string,
) Credential {
	return Credential{
		BaseHinter: hint.NewBaseHinter(CredentialHint),
		holder:     holder,
		templateID: templateID,
		id:         id,
		value:      value,
		validfrom:  validfrom,
		validuntil: validuntil,
		did:        did,
	}
}

func (it Credential) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.holder.Bytes(),
		it.templateID.Bytes(),
		[]byte(it.id),
		[]byte(it.value),
		it.validfrom.Bytes(),
		it.validuntil.Bytes(),
		[]byte(it.did),
	)
}

func (it Credential) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.holder,
		it.templateID,
		it.validfrom,
		it.validuntil,
	); err != nil {
		return err
	}

	if it.validuntil.n.Cmp(&it.validfrom.n) <= 0 {
		return util.ErrInvalid.Errorf("valid until <= valid from, %q <= %q", it.validuntil, it.validfrom)
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

func (it Credential) TemplateID() Uint256 {
	return it.templateID
}

func (it Credential) ValidFrom() Uint256 {
	return it.validfrom
}

func (it Credential) ValidUntil() Uint256 {
	return it.validuntil
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
