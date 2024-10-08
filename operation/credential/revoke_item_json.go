package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RevokeItemJSONMarshaler struct {
	hint.BaseHinter
	Contract     base.Address     `json:"contract"`
	Holder       base.Address     `json:"holder"`
	TemplateID   string           `json:"template_id"`
	CredentialID string           `json:"credential_id"`
	Currency     types.CurrencyID `json:"currency"`
}

func (it RevokeItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeItemJSONMarshaler{
		BaseHinter:   it.BaseHinter,
		Contract:     it.contract,
		Holder:       it.holder,
		TemplateID:   it.templateID,
		CredentialID: it.credentialID,
		Currency:     it.currency,
	})
}

type RevokeItemJSONUnmarshaler struct {
	Hint         hint.Hint `json:"_hint"`
	Contract     string    `json:"contract"`
	Holder       string    `json:"holder"`
	TemplateID   string    `json:"template_id"`
	CredentialID string    `json:"credential_id"`
	Currency     string    `json:"currency"`
}

func (it *RevokeItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uit RevokeItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.Holder,
		uit.TemplateID,
		uit.CredentialID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
