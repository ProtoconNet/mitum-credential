package credential // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it RevokeCredentialsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":                 it.Hint().String(),
			"contract":              it.contract,
			"credential_service_id": it.credentialServiceID,
			"holder":                it.holder,
			"template_id":           it.templateID,
			"id":                    it.id,
			"currency":              it.currency,
		},
	)
}

type RevokeCredentialsItemBSONUnmarshaler struct {
	Hint              string `bson:"_hint"`
	Contract          string `bson:"contract"`
	CredentialService string `json:"credential_service_id"`
	Holder            string `json:"holder"`
	TemplateID        string `json:"template_id"`
	ID                string `json:"id"`
	Currency          string `bson:"currency"`
}

func (it *RevokeCredentialsItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of RevokeCredentialsItem")

	var uit RevokeCredentialsItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, ht,
		uit.Contract,
		uit.CredentialService,
		uit.Holder,
		uit.TemplateID,
		uit.ID,
		uit.Currency,
	)
}
