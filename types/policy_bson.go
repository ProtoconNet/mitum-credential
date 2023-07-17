package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (h Holder) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            h.Hint().String(),
			"address":          h.address,
			"credential_count": h.credentialCount,
		},
	)
}

type HolderBSONUnmarshaler struct {
	Hint            string `bson:"_hint"`
	Address         string `bson:"address"`
	CredentialCount uint64 `bson:"credential_count"`
}

func (h *Holder) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Holder")

	var upo HolderBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return h.unpack(enc, ht, upo.Address, upo.CredentialCount)
}

func (po Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            po.Hint().String(),
			"templates":        po.templateIDs,
			"holders":          po.holders,
			"credential_count": po.credentialCount,
		},
	)
}

type PolicyBSONUnmarshaler struct {
	Hint            string   `bson:"_hint"`
	Templates       []uint64 `bson:"templates"`
	Holders         bson.Raw `bson:"holders"`
	CredentialCount uint64   `bson:"credential_count"`
}

func (po *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Policy")

	var upo PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return po.unpack(enc, ht, upo.Templates, upo.Holders, upo.CredentialCount)
}
