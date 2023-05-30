package credential

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AssignCredentialsItemJSONMarshaler struct {
	hint.BaseHinter
	Contract          base.Address                 `json:"contract"`
	CredentialService extensioncurrency.ContractID `json:"credential_service_id"`
	Holder            base.Address                 `json:"holder"`
	TemplateID        Uint256                      `json:"template_id"`
	ID                string                       `json:"id"`
	Value             string                       `json:"value"`
	ValidFrom         Uint256                      `json:"valid_from"`
	ValidUntil        Uint256                      `json:"valid_until"`
	DID               string                       `json:"did"`
	Currency          currency.CurrencyID          `json:"currency"`
}

func (it AssignCredentialsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AssignCredentialsItemJSONMarshaler{
		BaseHinter:        it.BaseHinter,
		Contract:          it.contract,
		CredentialService: it.credentialServiceID,
		Holder:            it.holder,
		TemplateID:        it.templateID,
		ID:                it.id,
		Value:             it.value,
		ValidFrom:         it.validfrom,
		ValidUntil:        it.validuntil,
		DID:               it.did,
		Currency:          it.currency,
	})
}

type AssignCredentialsItemJSONUnMarshaler struct {
	Hint              hint.Hint `json:"_hint"`
	Contract          string    `json:"contract"`
	CredentialService string    `json:"credential_service_id"`
	Holder            string    `json:"holder"`
	TemplateID        string    `json:"template_id"`
	ID                string    `json:"id"`
	Value             string    `json:"value"`
	ValidFrom         string    `json:"valid_from"`
	ValidUntil        string    `json:"valid_until"`
	DID               string    `json:"did"`
	Currency          string    `json:"currency"`
}

func (it *AssignCredentialsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AssignCredentialsItem")

	var uit AssignCredentialsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.CredentialService,
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
