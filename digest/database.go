package digest

import (
	"context"
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultColNameDIDCredentialService = "digest_did_issuer"
	DefaultColNameDIDCredential        = "digest_did_credential"
	DefaultColNameHolder               = "digest_did_holder_did"
	DefaultColNameTemplate             = "digest_did_template"
)

var maxLimit int64 = 50

func CredentialService(st *cdigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNameDIDCredentialService,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			de, err := state.StateDesignValue(sta)
			if err != nil {
				return err
			}
			design = &de

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return design, nil
}

func Credential(st *cdigest.Database, contract, templateID, credentialID string) (*types.Credential, bool, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("template", templateID)
	filter = filter.Add("credential_id", credentialID)

	var credential *types.Credential
	var isActive bool
	var sta base.State
	var err error
	if err = st.MongoClient().GetByFilter(
		DefaultColNameDIDCredential,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			cre, active, err := state.StateCredentialValue(sta)
			if err != nil {
				return err
			}
			credential = &cre
			isActive = active
			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, false, err
	}

	return credential, isActive, nil
}

func Template(st *cdigest.Database, contract, templateID string) (*types.Template, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("template", templateID)

	var template *types.Template
	var sta base.State
	var err error
	if err = st.MongoClient().GetByFilter(
		DefaultColNameTemplate,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			te, err := state.StateTemplateValue(sta)
			if err != nil {
				return err
			}
			template = &te
			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return template, nil
}

func HolderDID(st *cdigest.Database, contract, holder string) (string, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)

	var did string
	var sta base.State
	var err error
	if err = st.MongoClient().GetByFilter(
		DefaultColNameHolder,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			did, err = state.StateHolderDIDValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return "", err
	}

	return did, nil
}

func CredentialsByServiceTemplate(
	st *cdigest.Database,
	contract,
	templateID string,
	reverse bool,
	offset string,
	limit int64,
	callback func(types.Credential, bool, base.State) (bool, error),
) error {
	filter, err := buildCredentialFilterByServiceTemplate(contract, templateID, offset, reverse)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.MongoClient().Find(
		context.Background(),
		DefaultColNameDIDCredential,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := cdigest.LoadState(cursor.Decode, st.Encoders())
			if err != nil {
				return false, err
			}
			credential, isActive, err := state.StateCredentialValue(st)
			if err != nil {
				return false, err
			}
			return callback(credential, isActive, st)
		},
		opt,
	)
}

func buildCredentialFilterByServiceTemplate(contract, templateID string, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter for matching template
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterTemplate := bson.D{{"template", bson.D{{"$in", []string{templateID}}}}}
	filterA = append(filterA, filterContract)
	filterA = append(filterA, filterTemplate)

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"credential_id", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
		} else {
			filterHeight := bson.D{
				{"credential_id", bson.D{{"$lt", offset}}},
			}
			filterA = append(filterA, filterHeight)
		}
	}

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}

func CredentialsByServiceHolder(
	st *cdigest.Database,
	contract, holder string,
	callback func(types.Credential, bool, base.State) (bool, error),
) error {
	filter, err := buildCredentialFilterByServiceHolder(contract, holder)
	if err != nil {
		return err
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", 1).D(),
	)

	opt = opt.SetLimit(1000)

	return st.MongoClient().Find(
		context.Background(),
		DefaultColNameDIDCredential,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := cdigest.LoadState(cursor.Decode, st.Encoders())
			if err != nil {
				return false, err
			}
			credential, isActive, err := state.StateCredentialValue(st)
			if err != nil {
				return false, err
			}
			return callback(credential, isActive, st)
		},
		opt,
	)
}

func buildCredentialFilterByServiceHolder(contract, holder string) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterHolder := bson.D{{"d.value.credential.holder", holder}}
	filterA = append(filterA, filterContract)
	filterA = append(filterA, filterHolder)

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}
