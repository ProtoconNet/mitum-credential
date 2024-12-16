package state

import (
	"fmt"

	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"strings"
)

var (
	CredentialPrefix     = "credential"
	DesignStateValueHint = hint.MustNewHint("mitum-credential-design-state-value-v0.0.1")
	DesignSuffix         = "design"
)

func StateKeyCredentialPrefix(contract base.Address) string {
	return fmt.Sprintf("%s:%s", CredentialPrefix, contract.String())
}

type DesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (de DesignStateValue) Hint() hint.Hint {
	return de.BaseHinter.Hint()
}

func (de DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid credential DesignStateValue")

	if err := de.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := de.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (de DesignStateValue) HashBytes() []byte {
	return de.Design.Bytes()
}

func StateDesignValue(st base.State) (types.Design, error) {
	v := st.Value()
	if v == nil {
		return types.Design{}, util.ErrNotFound.Errorf("credential design not found in State")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return types.Design{}, errors.Errorf("invalid credential design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, CredentialPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func StateKeyDesign(contract base.Address) string {
	return fmt.Sprintf("%s:%s", StateKeyCredentialPrefix(contract), DesignSuffix)
}

var (
	TemplateStateValueHint = hint.MustNewHint("mitum-credential-template-state-value-v0.0.1")
	TemplateSuffix         = "template"
)

type TemplateStateValue struct {
	hint.BaseHinter
	Template types.Template
}

func NewTemplateStateValue(template types.Template) TemplateStateValue {
	return TemplateStateValue{
		BaseHinter: hint.NewBaseHinter(TemplateStateValueHint),
		Template:   template,
	}
}

func (t TemplateStateValue) Hint() hint.Hint {
	return t.BaseHinter.Hint()
}

func (t TemplateStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid credential TemplateStateValue")

	if err := t.BaseHinter.IsValid(TemplateStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := t.Template.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (t TemplateStateValue) HashBytes() []byte {
	return t.Template.Bytes()
}

func StateKeyTemplate(contract base.Address, templateID string) string {
	return fmt.Sprintf("%s:%s:%s",
		StateKeyCredentialPrefix(contract),
		templateID,
		TemplateSuffix,
	)
}

func IsStateTemplateKey(key string) bool {
	return strings.HasPrefix(key, CredentialPrefix) && strings.HasSuffix(key, TemplateSuffix)
}

func StateTemplateValue(st base.State) (types.Template, error) {
	v := st.Value()
	if v == nil {
		return types.Template{}, util.ErrNotFound.Errorf("template not found in State")
	}

	t, ok := v.(TemplateStateValue)
	if !ok {
		return types.Template{}, errors.Errorf("invalid template value found, %T", v)
	}

	return t.Template, nil
}

var (
	CredentialStateValueHint = hint.MustNewHint("mitum-credential-credential-state-value-v0.0.1")
	CredentialSuffix         = "credential"
)

type CredentialStateValue struct {
	hint.BaseHinter
	Credential types.Credential
	IsActive   bool
}

func NewCredentialStateValue(credential types.Credential, isActive bool) CredentialStateValue {
	return CredentialStateValue{
		BaseHinter: hint.NewBaseHinter(CredentialStateValueHint),
		Credential: credential,
		IsActive:   isActive,
	}
}

func (cd CredentialStateValue) Hint() hint.Hint {
	return cd.BaseHinter.Hint()
}

func (cd CredentialStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CredentialStateValue")

	if err := cd.BaseHinter.IsValid(CredentialStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := cd.Credential.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (cd CredentialStateValue) HashBytes() []byte {
	var v int8
	if cd.IsActive {
		v = 1
	}
	return util.ConcatBytesSlice([]byte{byte(v)}, cd.Credential.Bytes())
}

func StateKeyCredential(contract base.Address, templateID string, id string) string {
	return fmt.Sprintf(
		"%s:%s:%s:%s",
		StateKeyCredentialPrefix(contract), templateID,
		id,
		CredentialSuffix,
	)
}

func IsStateCredentialKey(key string) bool {
	return strings.HasPrefix(key, CredentialPrefix) && strings.HasSuffix(key, CredentialSuffix)
}

func StateCredentialValue(st base.State) (types.Credential, bool, error) {
	v := st.Value()
	if v == nil {
		return types.Credential{}, false, util.ErrNotFound.Errorf("credential not found in State")
	}

	c, ok := v.(CredentialStateValue)
	if !ok {
		return types.Credential{}, false, errors.Errorf("invalid credential value found, %T", v)
	}

	return c.Credential, c.IsActive, nil
}

var (
	HolderDIDStateValueHint = hint.MustNewHint("mitum-credential-holder-did-state-value-v0.0.1")
	HolderDIDSuffix         = "holder-did"
)

type HolderDIDStateValue struct {
	hint.BaseHinter
	did string
}

func NewHolderDIDStateValue(did string) HolderDIDStateValue {
	return HolderDIDStateValue{
		BaseHinter: hint.NewBaseHinter(HolderDIDStateValueHint),
		did:        did,
	}
}

func (hd HolderDIDStateValue) Hint() hint.Hint {
	return hd.BaseHinter.Hint()
}

func (hd HolderDIDStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid credential HolderDIDStateValue")

	if err := hd.BaseHinter.IsValid(HolderDIDStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if len(hd.did) < 1 {
		return e.Errorf("empty did")
	}

	return nil
}

func (hd HolderDIDStateValue) HashBytes() []byte {
	return []byte(hd.did)
}

func StateHolderDIDValue(st base.State) (string, error) {
	v := st.Value()
	if v == nil {
		return "", util.ErrNotFound.Errorf("holder did not found in State")
	}

	d, ok := v.(HolderDIDStateValue)
	if !ok {
		return "", errors.Errorf("invalid holder did value found, %T", v)
	}

	return d.did, nil
}

func IsStateHolderDIDKey(key string) bool {
	return strings.HasPrefix(key, CredentialPrefix) && strings.HasSuffix(key, HolderDIDSuffix)
}

func StateKeyHolderDID(contract base.Address, holder base.Address) string {
	return fmt.Sprintf("%s:%s:%s", StateKeyCredentialPrefix(contract), holder.String(), HolderDIDSuffix)
}
