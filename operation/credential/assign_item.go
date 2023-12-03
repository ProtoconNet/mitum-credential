package credential

import (
	"unicode/utf8"

	"github.com/ProtoconNet/mitum-credential/types"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AssignItemHint = hint.MustNewHint("mitum-credential-assign-item-v0.0.1")

type AssignItem struct {
	hint.BaseHinter
	contract   base.Address
	holder     base.Address
	templateID string
	id         string
	value      string
	validFrom  uint64
	validUntil uint64
	did        string
	currency   crcytypes.CurrencyID
}

func NewAssignItem(
	contract base.Address,
	holder base.Address,
	templateID string,
	id string,
	value string,
	validFrom uint64,
	validUntil uint64,
	did string,
	currency crcytypes.CurrencyID,
) AssignItem {
	return AssignItem{
		BaseHinter: hint.NewBaseHinter(AssignItemHint),
		contract:   contract,
		holder:     holder,
		templateID: templateID,
		id:         id,
		value:      value,
		validFrom:  validFrom,
		validUntil: validUntil,
		did:        did,
		currency:   currency,
	}
}

func (it AssignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.holder.Bytes(),
		[]byte(it.templateID),
		[]byte(it.id),
		[]byte(it.value),
		util.Uint64ToBytes(it.validFrom),
		util.Uint64ToBytes(it.validUntil),
		[]byte(it.did),
		it.currency.Bytes(),
	)
}

func (it AssignItem) IsValid([]byte) error {
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

	if it.validUntil <= it.validFrom {
		return util.ErrInvalid.Errorf("valid until <= valid from, %q <= %q", it.validUntil, it.validFrom)
	}

	if l := utf8.RuneCountInString(it.templateID); l < 1 || l > types.MaxLengthTemplateID {
		return util.ErrInvalid.Errorf("invalid length of template ID, 0 <= length <= %d", types.MaxLengthTemplateID)
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(it.templateID)) {
		return util.ErrInvalid.Errorf("invalid templateID due to the inclusion of special characters")
	}

	if l := utf8.RuneCountInString(it.id); l < 1 || l > types.MaxLengthCredentialID {
		return util.ErrInvalid.Errorf("invalid length of credential ID, 0 <= length <= %d", types.MaxLengthCredentialID)
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(it.id)) {
		return util.ErrInvalid.Errorf("invalid credential ID due to the inclusion of special characters")
	}

	if len(it.did) == 0 {
		return util.ErrInvalid.Errorf("empty did")
	}

	if l := utf8.RuneCountInString(it.value); l < 1 || l > types.MaxLengthCredentialValue {
		return util.ErrInvalid.Errorf("invalid length of value, 0 <= length <= %d", types.MaxLengthCredentialValue)
	}

	return nil
}

func (it AssignItem) Contract() base.Address {
	return it.contract
}

func (it AssignItem) Holder() base.Address {
	return it.holder
}

func (it AssignItem) TemplateID() string {
	return it.templateID
}

func (it AssignItem) ValidFrom() uint64 {
	return it.validFrom
}

func (it AssignItem) ValidUntil() uint64 {
	return it.validUntil
}

func (it AssignItem) ID() string {
	return it.id
}

func (it AssignItem) Value() string {
	return it.value
}

func (it AssignItem) DID() string {
	return it.did
}

func (it AssignItem) Currency() crcytypes.CurrencyID {
	return it.currency
}

func (it AssignItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
