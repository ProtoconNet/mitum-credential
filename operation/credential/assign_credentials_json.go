package credential

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type AssignCredentialsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address            `json:"sender"`
	Items  []AssignCredentialsItem `json:"items"`
}

func (fact AssignCredentialsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AssignCredentialsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type AssignCredentialsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *AssignCredentialsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AssignCredentialsFact")

	var uf AssignCredentialsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Sender, uf.Items)
}

type AssignCredentialsMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op AssignCredentials) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AssignCredentialsMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AssignCredentials) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AssignCredentials")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
