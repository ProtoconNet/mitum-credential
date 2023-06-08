package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateCredentialServiceFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner               base.Address            `json:"sender"`
	Contract            base.Address            `json:"contract"`
	CredentialServiceID currencybase.ContractID `json:"credential_service_id"`
	Currency            currencybase.CurrencyID `json:"currency"`
}

func (fact CreateCredentialServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateCredentialServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		CredentialServiceID:   fact.credentialServiceID,
		Currency:              fact.currency,
	})
}

type CreateCredentialServiceFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner               string `json:"sender"`
	Contract            string `json:"contract"`
	CredentialServiceID string `json:"credential_service_id"`
	Currency            string `json:"currency"`
}

func (fact *CreateCredentialServiceFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CreateCredentialServiceFact")

	var uf CreateCredentialServiceFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Contract, uf.CredentialServiceID, uf.Currency)
}

type CreateCredentialServiceMarshaler struct {
	currencybase.BaseOperationJSONMarshaler
}

func (op CreateCredentialService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateCredentialServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateCredentialService) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CreateCredentialService")

	var ubo currencybase.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
