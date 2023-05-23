package credential

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	CredentialService Design `json:"credential_service"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter:        de.BaseHinter,
		CredentialService: de.Design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	CredentialService json.RawMessage `json:"credential_service"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var design Design

	if err := design.DecodeJSON(u.CredentialService, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

type TemplateStateValueJSONMarshaler struct {
	hint.BaseHinter
	Template Template `json:"template"`
}

func (t TemplateStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TemplateStateValueJSONMarshaler{
		BaseHinter: t.BaseHinter,
		Template:   t.Template,
	})
}

type TemplateStateValueJSONUnmarshaler struct {
	Template json.RawMessage `json:"template"`
}

func (t *TemplateStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of TemplateStateValue")

	var u TemplateStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var template Template

	if err := template.DecodeJSON(u.Template, enc); err != nil {
		return e(err, "")
	}

	t.Template = template

	return nil
}

type CredentialStateValueJSONMarshaler struct {
	hint.BaseHinter
	Credential Credential `json:"credential"`
}

func (cd CredentialStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CredentialStateValueJSONMarshaler{
		BaseHinter: cd.BaseHinter,
		Credential: cd.Credential,
	})
}

type CredentialStateValueJSONUnmarshaler struct {
	Credential json.RawMessage `json:"credential"`
}

func (cd *CredentialStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CredentialStateValue")

	var u CredentialStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var credential Credential

	if err := credential.DecodeJSON(u.Credential, enc); err != nil {
		return e(err, "")
	}

	cd.Credential = credential

	return nil
}

type HolderDIDStateValueJSONMarshaler struct {
	hint.BaseHinter
	DID string `json:"did"`
}

func (hd HolderDIDStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(HolderDIDStateValueJSONMarshaler{
		BaseHinter: hd.BaseHinter,
		DID:        hd.did,
	})
}

type HolderDIDStateValueJSONUnmarshaler struct {
	DID string `json:"did"`
}

func (hd *HolderDIDStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of HolderDIDStateValue")

	var u HolderDIDStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	hd.did = u.DID

	return nil
}
