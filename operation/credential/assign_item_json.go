package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AssignItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address             `json:"contract"`
	ServiceID  types.ServiceID          `json:"service_id"`
	Holder     base.Address             `json:"holder"`
	TemplateID string                   `json:"template_id"`
	ID         string                   `json:"id"`
	Value      string                   `json:"value"`
	ValidFrom  uint64                   `json:"valid_from"`
	ValidUntil uint64                   `json:"valid_until"`
	DID        string                   `json:"did"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (it AssignItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AssignItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		ServiceID:  it.serviceID,
		Holder:     it.holder,
		TemplateID: it.templateID,
		ID:         it.id,
		Value:      it.value,
		ValidFrom:  it.validfrom,
		ValidUntil: it.validuntil,
		DID:        it.did,
		Currency:   it.currency,
	})
}

type AssignItemJSONUnMarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	ServiceID  string    `json:"service_id"`
	Holder     string    `json:"holder"`
	TemplateID string    `json:"template_id"`
	ID         string    `json:"id"`
	Value      string    `json:"value"`
	ValidFrom  uint64    `json:"valid_from"`
	ValidUntil uint64    `json:"valid_until"`
	DID        string    `json:"did"`
	Currency   string    `json:"currency"`
}

func (it *AssignItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of AssignItem")

	var uit AssignItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.ServiceID,
		uit.Holder,
		uit.TemplateID,
		uit.ID,
		uit.Value,
		uit.ValidFrom,
		uit.ValidUntil,
		uit.DID,
		uit.Currency,
	)
}
