package digest

import (
	"context"
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultColNameAccount              = "digest_ac"
	defaultColNameContractAccount      = "digest_ca"
	defaultColNameBalance              = "digest_bl"
	defaultColNameCurrency             = "digest_cr"
	defaultColNameOperation            = "digest_op"
	defaultColNameBlock                = "digest_bm"
	defaultColNameDIDCredentialService = "digest_did_issuer"
	defaultColNameDIDCredential        = "digest_did_credential"
	defaultColNameHolder               = "digest_did_holder_did"
	defaultColNameTemplate             = "digest_did_template"
)

var maxLimit int64 = 50

func CredentialService(st *currencydigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta mitumbase.State
	var err error
	if err := st.MongoClient().GetByFilter(
		defaultColNameDIDCredentialService,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
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

func Credential(st *currencydigest.Database, contract, templateID, credentialID string) (*types.Credential, bool, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("template", templateID)
	filter = filter.Add("credential_id", credentialID)

	var credential *types.Credential
	var isActive bool
	var sta mitumbase.State
	var err error
	if err = st.MongoClient().GetByFilter(
		defaultColNameDIDCredential,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
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

func Template(st *currencydigest.Database, contract, templateID string) (*types.Template, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("template", templateID)

	var template *types.Template
	var sta mitumbase.State
	var err error
	if err = st.MongoClient().GetByFilter(
		defaultColNameTemplate,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
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

func HolderDID(st *currencydigest.Database, contract, holder string) (string, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)

	var did string
	var sta mitumbase.State
	var err error
	if err = st.MongoClient().GetByFilter(
		defaultColNameHolder,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
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
	st *currencydigest.Database,
	contract,
	templateID string,
	reverse bool,
	offset string,
	limit int64,
	callback func(types.Credential, bool, mitumbase.State) (bool, error),
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
		defaultColNameDIDCredential,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.Encoders())
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
	st *currencydigest.Database,
	contract, holder string,
	callback func(types.Credential, bool, mitumbase.State) (bool, error),
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
		defaultColNameDIDCredential,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.Encoders())
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
