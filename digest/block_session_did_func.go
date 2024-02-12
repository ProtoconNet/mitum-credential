package digest

import (
	"context"

	"github.com/ProtoconNet/mitum-credential/state"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	crcystate "github.com/ProtoconNet/mitum-currency/v3/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	PrepareFuncs = map[string]currencydigest.BlockSessionPrepareFunc{
		"PrepareDID": PrepareDID,
	}
	HandlerFuncs = map[string]currencydigest.BlockSessionHandlerFunc{
		"HandleDIDServiceState": HandleDIDServiceState,
		"HandleCredentialState": HandleCredentialState,
		"HandleHolderDIDState":  HandleHolderDIDState,
		"HandleTemplateState":   HandleTemplateState,
	}
	CommitFuncs = map[string]currencydigest.BlockSessionCommitFunc{
		defaultColNameDIDCredentialService: currencydigest.CommitFunc,
		defaultColNameDIDCredential:        CommitCredentialFunc,
		defaultColNameHolder:               currencydigest.CommitFunc,
		defaultColNameTemplate:             currencydigest.CommitFunc,
	}
)

var credentialMap = map[string]struct{}{}

func PrepareDID(bs *currencydigest.BlockSession) error {
	if len(bs.States()) < 1 {
		return nil
	}

	var didModels []mongo.WriteModel
	var credentialModels []mongo.WriteModel
	var holderDIDModels []mongo.WriteModel
	var templateModels []mongo.WriteModel

	for i := range bs.States() {
		st := bs.States()[i]
		switch {
		case state.IsStateDesignKey(st.Key()):
			j, err := bs.HandlerFuncs["HandleDIDServiceState"](bs, st)
			if err != nil {
				return err
			}
			didModels = append(didModels, j...)
		case state.IsStateCredentialKey(st.Key()):
			j, err := bs.HandlerFuncs["HandleCredentialState"](bs, st)
			if err != nil {
				return err
			}
			credentialMap[st.Key()] = struct{}{}
			credentialModels = append(credentialModels, j...)
		case state.IsStateHolderDIDKey(st.Key()):
			j, err := bs.HandlerFuncs["HandleHolderDIDState"](bs, st)
			if err != nil {
				return err
			}
			holderDIDModels = append(holderDIDModels, j...)
		case state.IsStateTemplateKey(st.Key()):
			j, err := bs.HandlerFuncs["HandleTemplateState"](bs, st)
			if err != nil {
				return err
			}
			templateModels = append(templateModels, j...)
		default:
			continue
		}
	}

	err := bs.SetWriteModel(defaultColNameDIDCredentialService, didModels)
	if err != nil {
		return err
	}
	err = bs.SetWriteModel(defaultColNameDIDCredential, credentialModels)
	if err != nil {
		return err
	}
	err = bs.SetWriteModel(defaultColNameHolder, holderDIDModels)
	if err != nil {
		return err
	}
	return bs.SetWriteModel(defaultColNameTemplate, templateModels)
}

func HandleDIDServiceState(bs *currencydigest.BlockSession, st mitumbase.State) ([]mongo.WriteModel, error) {
	if issuerDoc, err := NewServiceDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(issuerDoc),
		}, nil
	}
}

func HandleCredentialState(bs *currencydigest.BlockSession, st mitumbase.State) ([]mongo.WriteModel, error) {
	if credentialDoc, err := NewCredentialDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(credentialDoc),
		}, nil
	}
}

func HandleHolderDIDState(bs *currencydigest.BlockSession, st mitumbase.State) ([]mongo.WriteModel, error) {
	if holderDidDoc, err := NewHolderDIDDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(holderDidDoc),
		}, nil
	}
}

func HandleTemplateState(bs *currencydigest.BlockSession, st mitumbase.State) ([]mongo.WriteModel, error) {
	if templateDoc, err := NewTemplateDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(templateDoc),
		}, nil
	}
}

func CommitCredentialFunc(bs *currencydigest.BlockSession, ctx context.Context, colName string, writeModels []mongo.WriteModel) error {
	if len(writeModels) > 0 {
		for key := range credentialMap {
			parsedKey, err := crcystate.ParseStateKey(key, state.CredentialPrefix, 5)
			if err != nil {
				return err
			}
			err = bs.Database().CleanByHeightColName(
				ctx,
				bs.BLock().Manifest().Height(),
				colName,
				"contract", parsedKey[1],
				"template", parsedKey[2],
				"credential_id", parsedKey[3],
			)
			if err != nil {
				return err
			}
		}

		if err := bs.WriteWriteModels(ctx, colName, writeModels); err != nil {
			return err
		}
	}

	return nil
}
