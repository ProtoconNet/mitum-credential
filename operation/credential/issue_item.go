package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"unicode/utf8"

	"github.com/ProtoconNet/mitum-credential/types"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var IssueItemHint = hint.MustNewHint("mitum-credential-issue-item-v0.0.1")

type IssueItem struct {
	hint.BaseHinter
	contract     base.Address
	holder       base.Address
	templateID   string
	credentialID string
	value        string
	validFrom    uint64
	validUntil   uint64
	did          string
	currency     crcytypes.CurrencyID
}

func NewIssueItem(
	contract base.Address,
	holder base.Address,
	templateID string,
	credentialID string,
	value string,
	validFrom uint64,
	validUntil uint64,
	did string,
	currency crcytypes.CurrencyID,
) IssueItem {
	return IssueItem{
		BaseHinter:   hint.NewBaseHinter(IssueItemHint),
		contract:     contract,
		holder:       holder,
		templateID:   templateID,
		credentialID: credentialID,
		value:        value,
		validFrom:    validFrom,
		validUntil:   validUntil,
		did:          did,
		currency:     currency,
	}
}

func (it IssueItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.holder.Bytes(),
		[]byte(it.templateID),
		[]byte(it.credentialID),
		[]byte(it.value),
		util.Uint64ToBytes(it.validFrom),
		util.Uint64ToBytes(it.validUntil),
		[]byte(it.did),
		it.currency.Bytes(),
	)
}

func (it IssueItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.holder,
		it.currency,
	); err != nil {
		return common.ErrItemInvalid.Wrap(err)
	}

	if it.contract.Equal(it.holder) {
		return common.ErrItemInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("contract address is same with holder, %q", it.holder)))
	}

	if it.validUntil <= it.validFrom {
		return common.ErrItemInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("valid until <= valid from, %q <= %q", it.validUntil, it.validFrom)))
	}

	if l := utf8.RuneCountInString(it.templateID); l < 1 || l > types.MaxLengthTemplateID {
		return common.ErrItemInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of template ID <= %d", types.MaxLengthTemplateID)))
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(it.templateID)) {
		return common.ErrItemInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("template ID %s, must match regex `^[^\\s:/?#\\[\\]@]*$`", it.templateID)))
	}

	if l := utf8.RuneCountInString(it.credentialID); l < 1 || l > types.MaxLengthCredentialID {
		return common.ErrItemInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of credential ID <= %d", types.MaxLengthCredentialID)))
	}

	if !crcytypes.ReSpcecialChar.Match([]byte(it.credentialID)) {
		return common.ErrItemInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("credential ID %s, must match regex `^[^\\s:/?#\\[\\]@]*$`", it.credentialID)))
	}

	if len(it.did) == 0 {
		return common.ErrItemInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("empty did")))
	}

	if l := utf8.RuneCountInString(it.value); l < 1 || l > types.MaxLengthCredentialValue {
		return common.ErrItemInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("0 <= length of credential value <= %d", types.MaxLengthCredentialValue)))
	}

	return nil
}

func (it IssueItem) Contract() base.Address {
	return it.contract
}

func (it IssueItem) Holder() base.Address {
	return it.holder
}

func (it IssueItem) TemplateID() string {
	return it.templateID
}

func (it IssueItem) ValidFrom() uint64 {
	return it.validFrom
}

func (it IssueItem) ValidUntil() uint64 {
	return it.validUntil
}

func (it IssueItem) CredentialID() string {
	return it.credentialID
}

func (it IssueItem) Value() string {
	return it.value
}

func (it IssueItem) DID() string {
	return it.did
}

func (it IssueItem) Currency() crcytypes.CurrencyID {
	return it.currency
}

func (it IssueItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.holder

	return ad
}
