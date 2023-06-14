package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RevokeCredentialsItemJSONMarshaler struct {
	hint.BaseHinter
	Contract          base.Address             `json:"contract"`
	CredentialService currencytypes.ContractID `json:"credential_service_id"`
	Holder            base.Address             `json:"holder"`
	TemplateID        types.Uint256            `json:"template_id"`
	ID                string                   `json:"id"`
	Currency          currencytypes.CurrencyID `json:"currency"`
}

func (it RevokeCredentialsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeCredentialsItemJSONMarshaler{
		BaseHinter:        it.BaseHinter,
		Contract:          it.contract,
		CredentialService: it.credentialServiceID,
		Holder:            it.holder,
		TemplateID:        it.templateID,
		ID:                it.id,
		Currency:          it.currency,
	})
}

type RevokeCredentialsItemJSONUnMarshaler struct {
	Hint              hint.Hint `json:"_hint"`
	Contract          string    `json:"contract"`
	CredentialService string    `json:"credential_service_id"`
	Holder            string    `json:"holder"`
	TemplateID        string    `json:"template_id"`
	ID                string    `json:"id"`
	Currency          string    `json:"currency"`
}

func (it *RevokeCredentialsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RevokeCredentialsItem")

	var uit RevokeCredentialsItemJSONUnMarshaler
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
		uit.Currency,
	)
}
