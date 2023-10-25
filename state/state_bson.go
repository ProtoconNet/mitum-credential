package state

import (
	"github.com/ProtoconNet/mitum-credential/types"
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
	e := util.StringError("failed to decode bson of DesignStateValue")

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	var design types.Design
	if err := design.DecodeBSON(u.CredentialService, enc); err != nil {
		return e.Wrap(err)
	}

	de.Design = design

	if err := de.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

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
	e := util.StringError("failed to decode bson of TemplateStateValue")

	var u TemplateStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	t.BaseHinter = hint.NewBaseHinter(ht)

	var template types.Template
	if err := template.DecodeBSON(u.Template, enc); err != nil {
		return e.Wrap(err)
	}

	t.Template = template

	if err := t.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

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

func (cd *CredentialStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of CredentialStateValue")

	var u CredentialStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	cd.BaseHinter = hint.NewBaseHinter(ht)

	var credential types.Credential
	if err := credential.DecodeBSON(u.Credential, enc); err != nil {
		return e.Wrap(err)
	}

	cd.Credential = credential

	if err := cd.IsValid(nil); err != nil {
		return e.Wrap(err)
	}
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
	e := util.StringError("failed to decode bson of HolderDIDStateValue")

	var u HolderDIDStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	hd.BaseHinter = hint.NewBaseHinter(ht)
	hd.did = u.DID

	if err := hd.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}
