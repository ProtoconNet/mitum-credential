package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RevokeCredentialsItemHint = hint.MustNewHint("mitum-credential-revoke-credentials-item-v0.0.1")

type RevokeCredentialsItem struct {
	hint.BaseHinter
	contract            base.Address
	credentialServiceID currencybase.ContractID
	holder              base.Address
	templateID          Uint256
	id                  string
	currency            currencybase.CurrencyID
}

func NewRevokeCredentialsItem(
	contract base.Address,
	credentialServiceID currencybase.ContractID,
	holder base.Address,
	templateID Uint256,
	id string,
	currency currencybase.CurrencyID,
) RevokeCredentialsItem {
	return RevokeCredentialsItem{
		BaseHinter:          hint.NewBaseHinter(RevokeCredentialsItemHint),
		contract:            contract,
		credentialServiceID: credentialServiceID,
		holder:              holder,
		templateID:          templateID,
		id:                  id,
		currency:            currency,
	}
}

func (it RevokeCredentialsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.credentialServiceID.Bytes(),
		it.holder.Bytes(),
		it.templateID.Bytes(),
		[]byte(it.id),
		it.currency.Bytes(),
	)
}

func (it RevokeCredentialsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.credentialServiceID,
		it.contract,
		it.holder,
		it.templateID,
		it.currency,
	); err != nil {
		return err
	}

	if it.contract.Equal(it.holder) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", it.holder)
	}

	if len(it.id) == 0 {
		return util.ErrInvalid.Errorf("empty id")
	}

	return nil
}

func (it RevokeCredentialsItem) CredentialServiceID() currencybase.ContractID {
	return it.credentialServiceID
}

func (it RevokeCredentialsItem) Contract() base.Address {
	return it.contract
}

func (it RevokeCredentialsItem) Holder() base.Address {
	return it.holder
}

func (it RevokeCredentialsItem) TemplateID() Uint256 {
	return it.templateID
}

func (it RevokeCredentialsItem) ID() string {
	return it.id
}

func (it RevokeCredentialsItem) Currency() currencybase.CurrencyID {
	return it.currency
}

func (it RevokeCredentialsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
