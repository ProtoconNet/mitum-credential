package credential

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type AddTemplateFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner               base.Address                 `json:"sender"`
	Contract            base.Address                 `json:"contract"`
	CredentialServiceID extensioncurrency.ContractID `json:"credential_service_id"`
	TemplateID          Uint256                      `json:"template_id"`
	TemplateName        string                       `json:"template_name"`
	ServiceDate         Date                         `json:"service_date"`
	ExpirationDate      Date                         `json:"expiration_date"`
	TemplateShare       Bool                         `json:"temcdplate_share"`
	MultiAudit          Bool                         `json:"multi_audit"`
	DisplayName         string                       `json:"display_name"`
	SubjectKey          string                       `json:"subject_key"`
	Description         string                       `json:"description"`
	Creator             base.Address                 `json:"creator"`
	Currency            currency.CurrencyID          `json:"currency"`
}

func (fact AddTemplateFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddTemplateFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		CredentialServiceID:   fact.credentialServiceID,
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
	Owner               string `json:"sender"`
	Contract            string `json:"contract"`
	CredentialServiceID string `json:"credential_service_id"`
	TemplateID          string `json:"template_id"`
	TemplateName        string `json:"template_name"`
	ServiceDate         string `json:"service_date"`
	ExpirationDate      string `json:"expiration_date"`
	TemplateShare       bool   `json:"template_share"`
	MultiAudit          bool   `json:"multi_audit"`
	DisplayName         string `json:"display_name"`
	SubjectKey          string `json:"subject_key"`
	Description         string `json:"description"`
	Creator             string `json:"creator"`
	Currency            string `json:"currency"`
}

func (fact *AddTemplateFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AddTemplateFact")

	var uf AddTemplateFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.CredentialServiceID,
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
	)
}

type AddTemplateMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op AddTemplate) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddTemplateMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AddTemplate) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AddTemplate")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
