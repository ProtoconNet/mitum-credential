package state

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-credential/types"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	CredentialService types.Design `json:"credential_service"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter:        de.BaseHinter,
		CredentialService: de.Design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	Hint              hint.Hint       `json:"_hint"`
	CredentialService json.RawMessage `json:"credential_service"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	de.BaseHinter = hint.NewBaseHinter(u.Hint)

	var design types.Design

	if err := design.DecodeJSON(u.CredentialService, enc); err != nil {
		return e.Wrap(err)
	}

	de.Design = design

	if err := de.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

type TemplateStateValueJSONMarshaler struct {
	hint.BaseHinter
	Template types.Template `json:"template"`
}

func (t TemplateStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TemplateStateValueJSONMarshaler{
		BaseHinter: t.BaseHinter,
		Template:   t.Template,
	})
}

type TemplateStateValueJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Template json.RawMessage `json:"template"`
}

func (t *TemplateStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of TemplateStateValue")

	var u TemplateStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	t.BaseHinter = hint.NewBaseHinter(u.Hint)

	var template types.Template

	if err := template.DecodeJSON(u.Template, enc); err != nil {
		return e.Wrap(err)
	}

	t.Template = template

	if err := t.IsValid(nil); err != nil {
		return e.Wrap(err)
	}
	return nil
}

type CredentialStateValueJSONMarshaler struct {
	hint.BaseHinter
	Credential types.Credential `json:"credential"`
}

func (cd CredentialStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CredentialStateValueJSONMarshaler{
		BaseHinter: cd.BaseHinter,
		Credential: cd.Credential,
	})
}

type CredentialStateValueJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	Credential json.RawMessage `json:"credential"`
}

func (cd *CredentialStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CredentialStateValue")

	var u CredentialStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	cd.BaseHinter = hint.NewBaseHinter(u.Hint)

	var credential types.Credential

	if err := credential.DecodeJSON(u.Credential, enc); err != nil {
		return e.Wrap(err)
	}

	cd.Credential = credential

	if err := cd.IsValid(nil); err != nil {
		return e.Wrap(err)
	}
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
	Hint hint.Hint `json:"_hint"`
	DID  string    `json:"did"`
}

func (hd *HolderDIDStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of HolderDIDStateValue")

	var u HolderDIDStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	hd.BaseHinter = hint.NewBaseHinter(u.Hint)

	hd.did = u.DID

	if err := hd.IsValid(nil); err != nil {
		return e.Wrap(err)
	}
	return nil
}
