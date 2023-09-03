package credential // nolint: dupl

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact AddTemplateFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":           fact.Hint().String(),
			"sender":          fact.sender,
			"contract":        fact.contract,
			"service_id":      fact.serviceID,
			"template_id":     fact.templateID,
			"template_name":   fact.templateName,
			"service_date":    fact.serviceDate,
			"expiration_date": fact.expirationDate,
			"template_share":  fact.templateShare,
			"multi_audit":     fact.multiAudit,
			"display_name":    fact.displayName,
			"subject_key":     fact.subjectKey,
			"description":     fact.description,
			"creator":         fact.creator,
			"currency":        fact.currency,
			"hash":            fact.BaseFact.Hash().String(),
			"token":           fact.BaseFact.Token(),
		},
	)
}

type AddTemplateFactBSONUnmarshaler struct {
	Hint           string `bson:"_hint"`
	Sender         string `bson:"sender"`
	Contract       string `bson:"contract"`
	ServiceID      string `bson:"service_id"`
	TemplateID     string `bson:"template_id"`
	TemplateName   string `bson:"template_name"`
	ServiceDate    string `bson:"service_date"`
	ExpirationDate string `bson:"expiration_date"`
	TemplateShare  bool   `bson:"template_share"`
	MultiAudit     bool   `bson:"multi_audit"`
	DisplayName    string `bson:"display_name"`
	SubjectKey     string `bson:"subject_key"`
	Description    string `bson:"description"`
	Creator        string `bson:"creator"`
	Currency       string `bson:"currency"`
}

func (fact *AddTemplateFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of AddTemplateFact")

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf AddTemplateFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc,
		uf.Sender,
		uf.Contract,
		uf.ServiceID,
		uf.TemplateID,
		uf.TemplateName,
		uf.ServiceDate,
		uf.ExpirationDate,
		uf.TemplateShare,
		uf.MultiAudit,
		uf.DisplayName,
		uf.SubjectKey,
		uf.Description,
		uf.Creator,
		uf.Currency)
}

func (op AddTemplate) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *AddTemplate) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of AddTemplate")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
