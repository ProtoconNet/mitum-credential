package credential

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AssignCredentialsItemHint = hint.MustNewHint("mitum-credential-assign-credentials-item-v0.0.1")

type AssignCredentialsItem struct {
	hint.BaseHinter
	contract            base.Address
	credentialServiceID extensioncurrency.ContractID
	holder              base.Address
	templateID          Uint256
	id                  string
	value               string
	validfrom           Uint256
	validuntil          Uint256
	did                 string
	currency            currency.CurrencyID
}

func NewAssignCredentialsItem(
	contract base.Address,
	credentialServiceID extensioncurrency.ContractID,
	holder base.Address,
	templateID Uint256,
	id string,
	value string,
	validfrom Uint256,
	validuntil Uint256,
	did string,
	currency currency.CurrencyID,
) AssignCredentialsItem {
	return AssignCredentialsItem{
		BaseHinter:          hint.NewBaseHinter(AssignCredentialsItemHint),
		contract:            contract,
		credentialServiceID: credentialServiceID,
		holder:              holder,
		templateID:          templateID,
		id:                  id,
		value:               value,
		validfrom:           validfrom,
		validuntil:          validuntil,
		did:                 did,
		currency:            currency,
	}
}

func (it AssignCredentialsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.credentialServiceID.Bytes(),
		it.holder.Bytes(),
		it.templateID.Bytes(),
		[]byte(it.id),
		[]byte(it.value),
		it.validfrom.Bytes(),
		it.validuntil.Bytes(),
		[]byte(it.did),
		it.currency.Bytes(),
	)
}

func (it AssignCredentialsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.credentialServiceID,
		it.contract,
		it.holder,
		it.templateID,
		it.validfrom,
		it.validuntil,
		it.currency,
	); err != nil {
		return err
	}

	if it.contract.Equal(it.holder) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", it.holder)
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

func (it AssignCredentialsItem) CredentialServiceID() extensioncurrency.ContractID {
	return it.credentialServiceID
}

func (it AssignCredentialsItem) Contract() base.Address {
	return it.contract
}

func (it AssignCredentialsItem) Holder() base.Address {
	return it.holder
}

func (it AssignCredentialsItem) TemplateID() Uint256 {
	return it.templateID
}

func (it AssignCredentialsItem) ValidFrom() Uint256 {
	return it.validfrom
}

func (it AssignCredentialsItem) ValidUntil() Uint256 {
	return it.validuntil
}

func (it AssignCredentialsItem) ID() string {
	return it.id
}

func (it AssignCredentialsItem) Value() string {
	return it.value
}

func (it AssignCredentialsItem) DID() string {
	return it.did
}

func (it AssignCredentialsItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it AssignCredentialsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
