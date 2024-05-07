package credential

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type AssignFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address `json:"sender"`
	Items  []AssignItem `json:"items"`
}

func (fact AssignFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AssignFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type AssignFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *AssignFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf AssignFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type AssignMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Assign) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AssignMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Assign) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
