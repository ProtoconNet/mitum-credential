package credential

import (
	"unicode/utf8"

	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var RevokeItemHint = hint.MustNewHint("mitum-credential-revoke-item-v0.0.1")

type RevokeItem struct {
	hint.BaseHinter
	contract     base.Address
	holder       base.Address
	templateID   string
	credentialID string
	currency     crcytypes.CurrencyID
}

func NewRevokeItem(
	contract base.Address,
	holder base.Address,
	templateID, credentialID string,
	currency crcytypes.CurrencyID,
) RevokeItem {
	return RevokeItem{
		BaseHinter:   hint.NewBaseHinter(RevokeItemHint),
		contract:     contract,
		holder:       holder,
		templateID:   templateID,
		credentialID: credentialID,
		currency:     currency,
	}
}

func (it RevokeItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.holder.Bytes(),
		[]byte(it.templateID),
		[]byte(it.credentialID),
		it.currency.Bytes(),
	)
}

func (it RevokeItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.holder,
		it.currency,
	); err != nil {
		return err
	}

	if it.contract.Equal(it.holder) {
		return common.ErrItemInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("contract address is same with holder, %q", it.holder)))
	}

	if l := utf8.RuneCountInString(it.templateID); l < 1 || l > types.MaxLengthTemplateID {
		return common.ErrItemInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of template ID <= %d", types.MaxLengthTemplateID)))
	}

	if !crcytypes.ReValidSpcecialCh.Match([]byte(it.templateID)) {
		return common.ErrItemInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("template ID %s, must match regex `^[^\\s:/?#\\[\\]$@]*$`", it.templateID)))
	}

	if l := utf8.RuneCountInString(it.credentialID); l < 1 || l > types.MaxLengthCredentialID {
		return common.ErrItemInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of credential ID <= %d", types.MaxLengthCredentialID)))
	}

	if !crcytypes.ReValidSpcecialCh.Match([]byte(it.credentialID)) {
		return common.ErrItemInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("credential ID %s, must match regex `^[^\\s:/?#\\[\\]$@]*$`", it.credentialID)))
	}

	return nil
}

func (it RevokeItem) Contract() base.Address {
	return it.contract
}

func (it RevokeItem) Holder() base.Address {
	return it.holder
}

func (it RevokeItem) TemplateID() string {
	return it.templateID
}

func (it RevokeItem) CredentialID() string {
	return it.credentialID
}

func (it RevokeItem) Currency() crcytypes.CurrencyID {
	return it.currency
}

func (it RevokeItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
