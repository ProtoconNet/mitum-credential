package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type AddTemplateFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner          base.Address      `json:"sender"`
	Contract       base.Address      `json:"contract"`
	TemplateID     string            `json:"template_id"`
	TemplateName   string            `json:"template_name"`
	ServiceDate    types.Date        `json:"service_date"`
	ExpirationDate types.Date        `json:"expiration_date"`
	TemplateShare  types.Bool        `json:"template_share"`
	MultiAudit     types.Bool        `json:"multi_audit"`
	DisplayName    string            `json:"display_name"`
	SubjectKey     string            `json:"subject_key"`
	Description    string            `json:"description"`
	Creator        base.Address      `json:"creator"`
	Currency       ctypes.CurrencyID `json:"currency"`
}

func (fact AddTemplateFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddTemplateFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		TemplateID:            fact.templateID,
		TemplateName:          fact.templateName,
		ServiceDate:           fact.serviceDate,
		ExpirationDate:        fact.expirationDate,
		TemplateShare:         fact.templateShare,
		MultiAudit:            fact.multiAudit,
		DisplayName:           fact.displayName,
		SubjectKey:            fact.subjectKey,
		Description:           fact.description,
		Creator:               fact.creator,
		Currency:              fact.currency,
	})
}

type AddTemplateFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner          string `json:"sender"`
	Contract       string `json:"contract"`
	TemplateID     string `json:"template_id"`
	TemplateName   string `json:"template_name"`
	ServiceDate    string `json:"service_date"`
	ExpirationDate string `json:"expiration_date"`
	TemplateShare  bool   `json:"template_share"`
	MultiAudit     bool   `json:"multi_audit"`
	DisplayName    string `json:"display_name"`
	SubjectKey     string `json:"subject_key"`
	Description    string `json:"description"`
	Creator        string `json:"creator"`
	Currency       string `json:"currency"`
}

func (fact *AddTemplateFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf AddTemplateFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.TemplateID,
		uf.TemplateName,
		uf.ServiceDate,
		uf.ExpirationDate,
		uf.TemplateShare,
		uf.MultiAudit,
		uf.DisplayName,
		uf.SubjectKey,
		uf.Description,
		uf.Creator,
		uf.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type OperationMarshaler struct {
	common.BaseOperationJSONMarshaler
	extras.BaseOperationExtensionsJSONMarshaler
}

func (op AddTemplate) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *AddTemplate) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
