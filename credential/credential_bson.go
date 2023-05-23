package credential

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (cd Credential) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       cd.Hint().String(),
			"holder":      cd.holder,
			"templateid":  cd.templateID,
			"id":          cd.id,
			"value":       cd.value,
			"valid_from":  cd.validfrom,
			"valid_until": cd.validuntil,
			"did":         cd.did,
		},
	)
}

type CredentialBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Holder     string `bson:"holder"`
	TemplateID string `bson:"templateid"`
	ID         string `bson:"id"`
	Value      string `bson:"value"`
	ValidFrom  string `bson:"valid_from"`
	ValidUntil string `bson:"valid_until"`
	DID        string `bson:"did"`
}

func (cd *Credential) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Credential")

	var u CredentialBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
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
