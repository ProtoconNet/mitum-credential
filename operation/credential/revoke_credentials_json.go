package credential

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type RevokeCredentialsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address            `json:"sender"`
	Items  []RevokeCredentialsItem `json:"items"`
}

func (fact RevokeCredentialsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeCredentialsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type RevokeCredentialsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *RevokeCredentialsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RevokeCredentialsFact")

	var uf RevokeCredentialsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Sender, uf.Items)
}

type RevokeCredentialsMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op RevokeCredentials) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeCredentialsMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *RevokeCredentials) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RevokeCredentials")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
