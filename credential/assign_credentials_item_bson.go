package credential // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it AssignCredentialsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":                 it.Hint().String(),
			"contract":              it.contract,
			"credential_service_id": it.credentialServiceID,
			"holder":                it.holder,
			"template_id":           it.templateID,
			"id":                    it.id,
			"value":                 it.value,
			"valid_from":            it.validfrom,
			"valid_until":           it.validuntil,
			"did":                   it.did,
			"currency":              it.currency,
		},
	)
}

type AssignCredentialsItemBSONUnmarshaler struct {
	Hint              string `bson:"_hint"`
	Contract          string `bson:"contract"`
	CredentialService string `json:"credential_service_id"`
	Holder            string `json:"holder"`
	TemplateID        string `json:"template_id"`
	ID                string `json:"id"`
	Value             string `json:"value"`
	ValidFrom         string `json:"valid_from"`
	ValidUntil        string `json:"valid_until"`
	DID               string `json:"did"`
	Currency          string `bson:"currency"`
}

func (it *AssignCredentialsItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AssignCredentialsItem")

	var uit AssignCredentialsItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht,
		uit.Contract,
		uit.CredentialService,
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
