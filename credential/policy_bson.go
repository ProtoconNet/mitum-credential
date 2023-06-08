package credential

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            po.Hint().String(),
			"templates":        po.templates,
			"holders":          po.holders,
			"credential_count": po.credentialCount,
		},
	)
}

type PolicyBSONUnmarshaler struct {
	Hint            string   `bson:"_hint"`
	Templates       []string `bson:"templates"`
	Holders         bson.Raw `bson:"holders"`
	CredentialCount uint64   `bson:"credential_count"`
}

func (po *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Policy")

	var upo PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e(err, "")
	}

	return po.unpack(enc, ht, upo.Templates, upo.Holders, upo.CredentialCount)
}
