package credential

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RevokeFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address `json:"sender"`
	Items  []RevokeItem `json:"items"`
}

func (fact RevokeFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type RevokeFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *RevokeFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf RevokeFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type RevokeMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Revoke) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Revoke) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
