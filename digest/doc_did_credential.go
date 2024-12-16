package digest

import (
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	mongodbst "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DIDCredentialDesignDoc struct {
	mongodbst.BaseDoc
	st base.State
	de types.Design
}

func NewDIDCredentialDesignDoc(st base.State, enc encoder.Encoder) (DIDCredentialDesignDoc, error) {
	de, err := state.StateDesignValue(st)
	if err != nil {
		return DIDCredentialDesignDoc{}, err
	}
	b, err := mongodbst.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DIDCredentialDesignDoc{}, err
	}

	return DIDCredentialDesignDoc{
		BaseDoc: b,
		st:      st,
		de:      de,
	}, nil
}

func (doc DIDCredentialDesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.CredentialPrefix, 3)
	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type TemplateDoc struct {
	mongodbst.BaseDoc
	st       base.State
	template types.Template
}

func NewTemplateDoc(st base.State, enc encoder.Encoder) (*TemplateDoc, error) {
	template, err := state.StateTemplateValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbst.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &TemplateDoc{
		BaseDoc:  b,
		st:       st,
		template: template,
	}, nil
}

func (doc TemplateDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.CredentialPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["template"] = parsedKey[2]
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type CredentialDoc struct {
	mongodbst.BaseDoc
	st         base.State
	credential types.Credential
	isActive   bool
}

func NewCredentialDoc(st base.State, enc encoder.Encoder) (*CredentialDoc, error) {
	credential, isActive, err := state.StateCredentialValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbst.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &CredentialDoc{
		BaseDoc:    b,
		st:         st,
		credential: credential,
		isActive:   isActive,
	}, nil
}

func (doc CredentialDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.CredentialPrefix, 5)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["template"] = parsedKey[2]
	m["credential_id"] = parsedKey[3]
	m["is_active"] = doc.isActive
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type HolderDIDDoc struct {
	mongodbst.BaseDoc
	st  base.State
	did string
}

func NewHolderDIDDoc(st base.State, enc encoder.Encoder) (*HolderDIDDoc, error) {
	did, err := state.StateHolderDIDValue(st)
	if err != nil {
		return nil, err
	}

	b, err := mongodbst.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &HolderDIDDoc{
		BaseDoc: b,
		st:      st,
		did:     did,
	}, nil
}

func (doc HolderDIDDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.CredentialPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["did"] = doc.did
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
