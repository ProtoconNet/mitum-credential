package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (c Credential) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       c.Hint().String(),
			"holder":      c.holder,
			"template_id": c.templateID,
			"id":          c.id,
			"value":       c.value,
			"valid_from":  c.validFrom,
			"valid_until": c.validUntil,
			"did":         c.did,
		},
	)
}

type CredentialBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Holder     string `bson:"holder"`
	TemplateID string `bson:"template_id"`
	ID         string `bson:"id"`
	Value      string `bson:"value"`
	ValidFrom  uint64 `bson:"valid_from"`
	ValidUntil uint64 `bson:"valid_until"`
	DID        string `bson:"did"`
}

func (c *Credential) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Credential")

	var u CredentialBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return c.unpack(enc, ht,
		u.Holder,
		u.TemplateID,
		u.ID,
		u.Value,
		u.ValidFrom,
		u.ValidUntil,
		u.DID,
	)
}
