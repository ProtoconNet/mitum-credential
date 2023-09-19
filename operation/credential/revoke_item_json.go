package credential

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RevokeItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address             `json:"contract"`
	Holder     base.Address             `json:"holder"`
	TemplateID string                   `json:"template_id"`
	ID         string                   `json:"id"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (it RevokeItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Holder:     it.holder,
		TemplateID: it.templateID,
		ID:         it.id,
		Currency:   it.currency,
	})
}

type RevokeItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Holder     string    `json:"holder"`
	TemplateID string    `json:"template_id"`
	ID         string    `json:"id"`
	Currency   string    `json:"currency"`
}

func (it *RevokeItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RevokeItem")

	var uit RevokeItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.Holder,
		uit.TemplateID,
		uit.ID,
		uit.Currency,
	)
}
