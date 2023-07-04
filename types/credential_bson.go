package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (cd Credential) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       cd.Hint().String(),
			"holder":      cd.holder,
			"template_id": cd.templateID,
			"id":          cd.id,
			"value":       cd.value,
			"valid_from":  cd.validFrom,
			"valid_until": cd.validUntil,
			"did":         cd.did,
		},
	)
}

type CredentialBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Holder     string `bson:"holder"`
	TemplateID uint64 `bson:"template_id"`
	ID         string `bson:"id"`
	Value      string `bson:"value"`
	ValidFrom  uint64 `bson:"valid_from"`
	ValidUntil uint64 `bson:"valid_until"`
	DID        string `bson:"did"`
}

func (cd *Credential) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Credential")

	var u CredentialBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return cd.unpack(enc, ht,
		u.Holder,
		u.TemplateID,
		u.ID,
		u.Value,
		u.ValidFrom,
		u.ValidUntil,
		u.DID,
	)
}
