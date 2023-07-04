package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TemplateJSONMarshaler struct {
	hint.BaseHinter
	TemplateID     uint64       `json:"template_id"`
	TemplateName   string       `json:"template_name"`
	ServiceDate    Date         `json:"service_date"`
	ExpirationDate Date         `json:"expiration_date"`
	TemplateShare  Bool         `json:"template_share"`
	MultiAudit     Bool         `json:"multi_audit"`
	DisplayName    string       `json:"display_name"`
	SubjectKey     string       `json:"subject_key"`
	Description    string       `json:"description"`
	Creator        base.Address `json:"creator"`
}

func (t Template) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TemplateJSONMarshaler{
		BaseHinter:     t.BaseHinter,
		TemplateID:     t.templateID,
		TemplateName:   t.templateName,
		ServiceDate:    t.serviceDate,
		ExpirationDate: t.expirationDate,
		TemplateShare:  t.templateShare,
		MultiAudit:     t.multiAudit,
		DisplayName:    t.displayName,
		SubjectKey:     t.subjectKey,
		Description:    t.description,
		Creator:        t.creator,
	})
}

type TemplateJSONUnmarshaler struct {
	Hint           hint.Hint `json:"_hint"`
	TemplateID     uint64    `json:"template_id"`
	TemplateName   string    `json:"template_name"`
	ServiceDate    string    `json:"service_date"`
	ExpirationDate string    `json:"expiration_date"`
	TemplateShare  bool      `json:"template_share"`
	MultiAudit     bool      `json:"multi_audit"`
	DisplayName    string    `json:"display_name"`
	SubjectKey     string    `json:"subject_key"`
	Description    string    `json:"description"`
	Creator        string    `json:"creator"`
}

func (t *Template) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Template")

	var u TemplateJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return t.unpack(enc, u.Hint,
		u.TemplateID,
		u.TemplateName,
		u.ServiceDate,
		u.ExpirationDate,
		u.TemplateShare,
		u.MultiAudit,
		u.DisplayName,
		u.SubjectKey,
		u.Description,
		u.Creator,
	)
}
