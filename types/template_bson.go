package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t Template) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":           t.Hint().String(),
			"templateid":      t.templateID,
			"template_name":   t.templateName,
			"service_date":    t.serviceDate,
			"expiration_date": t.expirationDate,
			"template_share":  t.templateShare,
			"multi_audit":     t.multiAudit,
			"display_name":    t.displayName,
			"subject_key":     t.subjectKey,
			"description":     t.description,
			"creator":         t.creator,
		},
	)
}

type TemplateBSONUnmarshaler struct {
	Hint           string `bson:"_hint"`
	TemplateID     string `json:"templateid"`
	TemplateName   string `json:"template_name"`
	ServiceDate    string `json:"service_date"`
	ExpirationDate string `json:"expiration_date"`
	TemplateShare  bool   `json:"template_share"`
	MultiAudit     bool   `json:"multi_audit"`
	DisplayName    string `json:"display_name"`
	SubjectKey     string `json:"subject_key"`
	Description    string `json:"description"`
	Creator        string `json:"creator"`
}

func (t *Template) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Template")

	var u TemplateBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return t.unpack(enc, ht,
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
