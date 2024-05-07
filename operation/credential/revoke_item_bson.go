package credential // nolint:dupl

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it RevokeItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       it.Hint().String(),
			"contract":    it.contract,
			"holder":      it.holder,
			"template_id": it.templateID,
			"id":          it.id,
			"currency":    it.currency,
		},
	)
}

type RevokeItemBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Contract   string `bson:"contract"`
	Holder     string `bson:"holder"`
	TemplateID string `bson:"template_id"`
	ID         string `bson:"id"`
	Currency   string `bson:"currency"`
}

func (it *RevokeItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit RevokeItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht,
		uit.Contract,
		uit.Holder,
		uit.TemplateID,
		uit.ID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	return nil
}
