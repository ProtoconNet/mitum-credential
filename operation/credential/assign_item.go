package credential

import (
	"unicode/utf8"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
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
	validfrom  uint64
	validuntil uint64
	did        string
	currency   currencytypes.CurrencyID
}

func NewAssignItem(
	contract base.Address,
	holder base.Address,
	templateID string,
	id string,
	value string,
	validfrom uint64,
	validuntil uint64,
	did string,
	currency currencytypes.CurrencyID,
) AssignItem {
	return AssignItem{
		BaseHinter: hint.NewBaseHinter(AssignItemHint),
		contract:   contract,
		holder:     holder,
		templateID: templateID,
		id:         id,
		value:      value,
		validfrom:  validfrom,
		validuntil: validuntil,
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
		util.Uint64ToBytes(it.validfrom),
		util.Uint64ToBytes(it.validuntil),
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

	if it.validuntil <= it.validfrom {
		return util.ErrInvalid.Errorf("valid until <= valid from, %q <= %q", it.validuntil, it.validfrom)
	}

	if l := utf8.RuneCountInString(it.templateID); l < 1 || l > MaxLengthTemplateID {
		return util.ErrInvalid.Errorf("invalid length of template ID, 0 <= length <= %d", MaxLengthTemplateID)
	}

	if l := utf8.RuneCountInString(it.id); l < 1 || l > MaxLengthCredentialID {
		return util.ErrInvalid.Errorf("invalid length of ID, 0 <= length <= %d", MaxLengthCredentialID)
	}

	if len(it.did) == 0 {
		return util.ErrInvalid.Errorf("empty did")
	}

	if l := utf8.RuneCountInString(it.value); l < 1 || l > MaxLengthCredentialValue {
		return util.ErrInvalid.Errorf("invalid length of value, 0 <= length <= %d", MaxLengthCredentialValue)
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
	return it.validfrom
}

func (it AssignItem) ValidUntil() uint64 {
	return it.validuntil
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

func (it AssignItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it AssignItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
