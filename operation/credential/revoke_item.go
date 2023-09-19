package credential

import (
	"unicode/utf8"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RevokeItemHint = hint.MustNewHint("mitum-credential-item-v0.0.1")

type RevokeItem struct {
	hint.BaseHinter
	contract   base.Address
	holder     base.Address
	templateID string
	id         string
	currency   currencytypes.CurrencyID
}

func NewRevokeItem(
	contract base.Address,
	holder base.Address,
	templateID, id string,
	currency currencytypes.CurrencyID,
) RevokeItem {
	return RevokeItem{
		BaseHinter: hint.NewBaseHinter(RevokeItemHint),
		contract:   contract,
		holder:     holder,
		templateID: templateID,
		id:         id,
		currency:   currency,
	}
}

func (it RevokeItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.holder.Bytes(),
		[]byte(it.templateID),
		[]byte(it.id),
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
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", it.holder)
	}

	if l := utf8.RuneCountInString(it.templateID); l < 1 || l > MaxLengthTemplateID {
		return util.ErrInvalid.Errorf("invalid length of template ID, 0 <= length <= %d", MaxLengthTemplateID)
	}

	if l := utf8.RuneCountInString(it.id); l < 1 || l > MaxLengthCredentialID {
		return util.ErrInvalid.Errorf("invalid length of ID, 0 <= length <= %d", MaxLengthCredentialID)
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

func (it RevokeItem) ID() string {
	return it.id
}

func (it RevokeItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it RevokeItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
