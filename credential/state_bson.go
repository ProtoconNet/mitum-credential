package credential

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (de DesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":              de.Hint().String(),
			"credential_service": de.Design,
		},
	)
}

type DesignStateValueBSONUnmarshaler struct {
	Hint              string   `bson:"_hint"`
	CredentialService bson.Raw `bson:"credential_service"`
}

func (de *DesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of DesignStateValue")

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	var design Design
	if err := design.DecodeBSON(u.CredentialService, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

func (t TemplateStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    t.Hint().String(),
			"template": t.Template,
		},
	)
}

type TemplateStateValueBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Template bson.Raw `bson:"template"`
}

func (t *TemplateStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of TemplateStateValue")

	var u TemplateStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	t.BaseHinter = hint.NewBaseHinter(ht)

	var template Template
	if err := template.DecodeBSON(u.Template, enc); err != nil {
		return e(err, "")
	}

	t.Template = template

	return nil
}

func (cd CredentialStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      cd.Hint().String(),
			"credential": cd.Credential,
		},
	)
}

type CredentialStateValueBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Credential bson.Raw `bson:"credential"`
}

func (t *CredentialStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CredentialStateValue")

	var u CredentialStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	t.BaseHinter = hint.NewBaseHinter(ht)

	var credential Credential
	if err := credential.DecodeBSON(u.Credential, enc); err != nil {
		return e(err, "")
	}

	t.Credential = credential

	return nil
}

func (hd HolderDIDStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": hd.Hint().String(),
			"did":   hd.did,
		},
	)
}

type HolderDIDStateValueBSONUnmarshaler struct {
	Hint string `bson:"_hint"`
	DID  string `bson:"did"`
}

func (hd *HolderDIDStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of HolderDIDStateValue")

	var u HolderDIDStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	hd.BaseHinter = hint.NewBaseHinter(ht)
	hd.did = u.DID

	return nil
}
