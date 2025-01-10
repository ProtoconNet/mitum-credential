package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var didCredentialServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "credential_service_contract_height"),
	},
}

var didCredentialIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "template", Value: 1},
			bson.E{Key: "credential_id", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "credential_id_contract_template_height"),
	},
}

var didCredentialHolderIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "holder", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "credential_holder_contract_height"),
	},
}

var DidCredentialTemplateIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "template", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "credential_template_contract_height"),
	},
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	DefaultIndexes[DefaultColNameDIDCredentialService] = didCredentialServiceIndexModels
	DefaultIndexes[DefaultColNameDIDCredential] = didCredentialIndexModels
	DefaultIndexes[DefaultColNameHolder] = didCredentialHolderIndexModels
	DefaultIndexes[DefaultColNameTemplate] = DidCredentialTemplateIndexModels
}
