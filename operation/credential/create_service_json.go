package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type CreateServiceFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner    base.Address             `json:"sender"`
	Contract base.Address             `json:"contract"`
	Currency currencytypes.CurrencyID `json:"currency"`
}

func (fact CreateServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	})
}

type CreateServiceFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner    string `json:"sender"`
	Contract string `json:"contract"`
	Currency string `json:"currency"`
}

func (fact *CreateServiceFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf CreateServiceFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Owner, uf.Contract, uf.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type CreateServiceMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateService) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
