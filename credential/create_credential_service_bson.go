package credential // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/base"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact CreateCredentialServiceFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":                 fact.Hint().String(),
			"sender":                fact.sender,
			"contract":              fact.contract,
			"credential_service_id": fact.credentialServiceID,
			"currency":              fact.currency,
			"hash":                  fact.BaseFact.Hash().String(),
			"token":                 fact.BaseFact.Token(),
		},
	)
}

type CreateCredentialServiceFactBSONUnmarshaler struct {
	Hint                string `bson:"_hint"`
	Sender              string `bson:"sender"`
	Contract            string `bson:"contract"`
	CredentialServiceID string `bson:"credential_service_id"`
	Currency            string `bson:"currency"`
}

func (fact *CreateCredentialServiceFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateCredentialServiceFact")

	var ubf base.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf CreateCredentialServiceFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.Sender, uf.Contract, uf.CredentialServiceID, uf.Currency)
}

func (op CreateCredentialService) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CreateCredentialService) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateCredentialService")

	var ubo base.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
