package credential

import (
	"fmt"
	"strings"

	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	CredentialServicePrefix = "credentialservice:"
	DesignStateValueHint    = hint.MustNewHint("mitum-credential-design-state-value-v0.0.1")
	DesignSuffix            = ":design"
)

func StateKeyCredentialServicePrefix(ca base.Address, credentialServiceID currencytypes.ContractID) string {
	return fmt.Sprintf("%s%s:%s", CredentialServicePrefix, ca.String(), credentialServiceID)
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

func (hd DesignStateValue) Hint() hint.Hint {
	return hd.BaseHinter.Hint()
}

func (hd DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := hd.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := hd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (hd DesignStateValue) HashBytes() []byte {
	return hd.Design.Bytes()
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
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, DesignSuffix)
}

func StateKeyDesign(ca base.Address, crid currencytypes.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyCredentialServicePrefix(ca, crid), DesignSuffix)
}

var (
	TemplateStateValueHint = hint.MustNewHint("mitum-credential-template-state-value-v0.0.1")
	TemplateSuffix         = ":credential-template"
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

func (sv TemplateStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv TemplateStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TemplateStateValue")

	if err := sv.BaseHinter.IsValid(TemplateStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Template.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv TemplateStateValue) HashBytes() []byte {
	return sv.Template.Bytes()
}

func StateKeyTemplate(ca base.Address, credentialServiceID currencytypes.ContractID, templateID types.Uint256) string {
	return fmt.Sprintf("%s-%s%s", StateKeyCredentialServicePrefix(ca, credentialServiceID), templateID.String(), TemplateSuffix)
}

func IsStateTemplateKey(key string) bool {
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, TemplateSuffix)
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
	CredentialSuffix         = ":credential"
)

type CredentialStateValue struct {
	hint.BaseHinter
	Credential types.Credential
}

func NewCredentialStateValue(credential types.Credential) CredentialStateValue {
	return CredentialStateValue{
		BaseHinter: hint.NewBaseHinter(CredentialStateValueHint),
		Credential: credential,
	}
}

func (sv CredentialStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv CredentialStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CredentialStateValue")

	if err := sv.BaseHinter.IsValid(CredentialStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Credential.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv CredentialStateValue) HashBytes() []byte {
	return sv.Credential.Bytes()
}

func StateKeyCredential(ca base.Address, credentialServiceID currencytypes.ContractID, templateID types.Uint256, id string) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeyCredentialServicePrefix(ca, credentialServiceID), templateID.String(), id, CredentialSuffix)
}

func IsStateCredentialKey(key string) bool {
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, CredentialSuffix)
}

func StateCredentialValue(st base.State) (types.Credential, error) {
	v := st.Value()
	if v == nil {
		return types.Credential{}, util.ErrNotFound.Errorf("crednetial not found in State")
	}

	c, ok := v.(CredentialStateValue)
	if !ok {
		return types.Credential{}, errors.Errorf("invalid credential value found, %T", v)
	}

	return c.Credential, nil
}

var (
	HolderDIDStateValueHint = hint.MustNewHint("mitum-credential-holder-did-state-value-v0.0.1")
	HolderDIDSuffix         = ":holder-did"
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

	return nil
}

func (hd HolderDIDStateValue) HashBytes() []byte {
	return []byte(hd.did)
}

func StateHolderDIDValue(st base.State) (*string, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("holder did not found in State")
	}

	d, ok := v.(HolderDIDStateValue)
	if !ok {
		return nil, errors.Errorf("invalid holder did value found, %T", v)
	}

	return &d.did, nil
}

func IsStateHolderDIDKey(key string) bool {
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, HolderDIDSuffix)
}

func StateKeyHolderDID(ca base.Address, credentialServiceID currencytypes.ContractID, ha base.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyCredentialServicePrefix(ca, credentialServiceID), ha.String(), HolderDIDSuffix)
}
