package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (t Template) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":           t.Hint().String(),
			"template_id":     t.templateID,
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
	TemplateID     string `bson:"template_id"`
	TemplateName   string `bson:"template_name"`
	ServiceDate    string `bson:"service_date"`
	ExpirationDate string `bson:"expiration_date"`
	TemplateShare  bool   `bson:"template_share"`
	MultiAudit     bool   `bson:"multi_audit"`
	DisplayName    string `bson:"display_name"`
	SubjectKey     string `bson:"subject_key"`
	Description    string `bson:"description"`
	Creator        string `bson:"creator"`
}

func (t *Template) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of Template")

	var u TemplateBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
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
