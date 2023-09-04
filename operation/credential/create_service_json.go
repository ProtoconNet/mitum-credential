package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateServiceFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner     base.Address             `json:"sender"`
	Contract  base.Address             `json:"contract"`
	ServiceID types.ServiceID          `json:"service_id"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (fact CreateServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		ServiceID:             fact.serviceID,
		Currency:              fact.currency,
	})
}

type CreateServiceFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner     string `json:"sender"`
	Contract  string `json:"contract"`
	ServiceID string `json:"service_id"`
	Currency  string `json:"currency"`
}

func (fact *CreateServiceFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateServiceFact")

	var uf CreateServiceFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Contract, uf.ServiceID, uf.Currency)
}

type CreateServiceMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateService) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateService")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
