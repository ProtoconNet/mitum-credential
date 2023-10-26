package credential // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it AssignItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       it.Hint().String(),
			"contract":    it.contract,
			"holder":      it.holder,
			"template_id": it.templateID,
			"id":          it.id,
			"value":       it.value,
			"valid_from":  it.validFrom,
			"valid_until": it.validUntil,
			"did":         it.did,
			"currency":    it.currency,
		},
	)
}

type AssignItemBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Contract   string `bson:"contract"`
	Holder     string `bson:"holder"`
	TemplateID string `bson:"template_id"`
	ID         string `bson:"id"`
	Value      string `bson:"value"`
	ValidFrom  uint64 `bson:"valid_from"`
	ValidUntil uint64 `bson:"valid_until"`
	DID        string `bson:"did"`
	Currency   string `bson:"currency"`
}

func (it *AssignItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of AssignItem")

	var uit AssignItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, ht,
		uit.Contract,
		uit.Holder,
		uit.TemplateID,
		uit.ID,
		uit.Value,
		uit.ValidFrom,
		uit.ValidUntil,
		uit.DID,
		uit.Currency,
	)
}
