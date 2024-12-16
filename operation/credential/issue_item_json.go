package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type IssueItemJSONMarshaler struct {
	hint.BaseHinter
	Contract     base.Address      `json:"contract"`
	Holder       base.Address      `json:"holder"`
	TemplateID   string            `json:"template_id"`
	CredentialID string            `json:"credential_id"`
	Value        string            `json:"value"`
	ValidFrom    uint64            `json:"valid_from"`
	ValidUntil   uint64            `json:"valid_until"`
	DID          string            `json:"did"`
	Currency     ctypes.CurrencyID `json:"currency"`
}

func (it IssueItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueItemJSONMarshaler{
		BaseHinter:   it.BaseHinter,
		Contract:     it.contract,
		Holder:       it.holder,
		TemplateID:   it.templateID,
		CredentialID: it.credentialID,
		Value:        it.value,
		ValidFrom:    it.validFrom,
		ValidUntil:   it.validUntil,
		DID:          it.did,
		Currency:     it.currency,
	})
}

type IssueItemJSONUnMarshaler struct {
	Hint         hint.Hint `json:"_hint"`
	Contract     string    `json:"contract"`
	Holder       string    `json:"holder"`
	TemplateID   string    `json:"template_id"`
	CredentialID string    `json:"credential_id"`
	Value        string    `json:"value"`
	ValidFrom    uint64    `json:"valid_from"`
	ValidUntil   uint64    `json:"valid_until"`
	DID          string    `json:"did"`
	Currency     string    `json:"currency"`
}

func (it *IssueItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uit IssueItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.Holder,
		uit.TemplateID,
		uit.CredentialID,
		uit.Value,
		uit.ValidFrom,
		uit.ValidUntil,
		uit.DID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
