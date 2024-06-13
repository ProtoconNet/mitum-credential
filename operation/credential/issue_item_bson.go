package credential // nolint:dupl

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it IssueItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":         it.Hint().String(),
			"contract":      it.contract,
			"holder":        it.holder,
			"template_id":   it.templateID,
			"credential_id": it.credentialID,
			"value":         it.value,
			"valid_from":    it.validFrom,
			"valid_until":   it.validUntil,
			"did":           it.did,
			"currency":      it.currency,
		},
	)
}

type IssueItemBSONUnmarshaler struct {
	Hint         string `bson:"_hint"`
	Contract     string `bson:"contract"`
	Holder       string `bson:"holder"`
	TemplateID   string `bson:"template_id"`
	CredentialID string `bson:"credential_id"`
	Value        string `bson:"value"`
	ValidFrom    uint64 `bson:"valid_from"`
	ValidUntil   uint64 `bson:"valid_until"`
	DID          string `bson:"did"`
	Currency     string `bson:"currency"`
}

func (it *IssueItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit IssueItemBSONUnmarshaler
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
		uit.CredentialID,
		uit.Value,
		uit.ValidFrom,
		uit.ValidUntil,
		uit.DID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	return nil
}
