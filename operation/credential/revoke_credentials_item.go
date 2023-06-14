package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RevokeCredentialsItemHint = hint.MustNewHint("mitum-credential-revoke-credentials-item-v0.0.1")

type RevokeCredentialsItem struct {
	hint.BaseHinter
	contract            base.Address
	credentialServiceID currencytypes.ContractID
	holder              base.Address
	templateID          types.Uint256
	id                  string
	currency            currencytypes.CurrencyID
}

func NewRevokeCredentialsItem(
	contract base.Address,
	credentialServiceID currencytypes.ContractID,
	holder base.Address,
	templateID types.Uint256,
	id string,
	currency currencytypes.CurrencyID,
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

func (it RevokeCredentialsItem) CredentialServiceID() currencytypes.ContractID {
	return it.credentialServiceID
}

func (it RevokeCredentialsItem) Contract() base.Address {
	return it.contract
}

func (it RevokeCredentialsItem) Holder() base.Address {
	return it.holder
}

func (it RevokeCredentialsItem) TemplateID() types.Uint256 {
	return it.templateID
}

func (it RevokeCredentialsItem) ID() string {
	return it.id
}

func (it RevokeCredentialsItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it RevokeCredentialsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
