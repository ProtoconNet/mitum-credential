package credential

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
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

type StateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewStateValueMerger(height base.Height, key string, st base.State) *StateValueMerger {
	s := &StateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	StateValueMergerFunc := func(height base.Height, st base.State) base.StateValueMerger {
		return NewStateValueMerger(height, key, st)
	}

	return base.NewBaseStateMergeValue(
		key,
		stv,
		StateValueMergerFunc,
	)
}

func StateKeyCredentialServicePrefix(ca base.Address, credentialServiceID extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s:%s", CredentialServicePrefix, ca.String(), credentialServiceID)
}

type DesignStateValue struct {
	hint.BaseHinter
	Design Design
}

func NewDesignStateValue(design Design) DesignStateValue {
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

func StateDesignValue(st base.State) (Design, error) {
	v := st.Value()
	if v == nil {
		return Design{}, util.ErrNotFound.Errorf("credential design not found in State")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return Design{}, errors.Errorf("invalid credential design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, DesignSuffix)
}

func StateKeyDesign(ca base.Address, crid extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyCredentialServicePrefix(ca, crid), DesignSuffix)
}

var (
	TemplateStateValueHint = hint.MustNewHint("mitum-credential-template-state-value-v0.0.1")
	TemplateSuffix         = ":credential-template"
)

type TemplateStateValue struct {
	hint.BaseHinter
	Template Template
}

func NewTemplateStateValue(template Template) TemplateStateValue {
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

func StateKeyTemplate(ca base.Address, credentialServiceID extensioncurrency.ContractID, templateID Uint256) string {
	return fmt.Sprintf("%s-%s%s", StateKeyCredentialServicePrefix(ca, credentialServiceID), templateID.String(), TemplateSuffix)
}

func IsStateTemplateKey(key string) bool {
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, TemplateSuffix)
}

func StateTemplateValue(st base.State) (Template, error) {
	v := st.Value()
	if v == nil {
		return Template{}, util.ErrNotFound.Errorf("template not found in State")
	}

	t, ok := v.(TemplateStateValue)
	if !ok {
		return Template{}, errors.Errorf("invalid template value found, %T", v)
	}

	return t.Template, nil
}

var (
	CredentialStateValueHint = hint.MustNewHint("mitum-credential-credential-state-value-v0.0.1")
	CredentialSuffix         = ":credential"
)

type CredentialStateValue struct {
	hint.BaseHinter
	Credential Credential
}

func NewCredentialStateValue(credential Credential) CredentialStateValue {
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

func StateKeyCredential(ca base.Address, credentialServiceID extensioncurrency.ContractID, ha base.Address, templateID Uint256, id string) string {
	return fmt.Sprintf("%s-%s-%s-%s%s", StateKeyCredentialServicePrefix(ca, credentialServiceID), ha.String(), templateID.String(), id, CredentialSuffix)
}

func IsStateCredentialKey(key string) bool {
	return strings.HasPrefix(key, CredentialServicePrefix) && strings.HasSuffix(key, CredentialSuffix)
}

func StateCredentialValue(st base.State) (Credential, error) {
	v := st.Value()
	if v == nil {
		return Credential{}, util.ErrNotFound.Errorf("crednetial not found in State")
	}

	c, ok := v.(CredentialStateValue)
	if !ok {
		return Credential{}, errors.Errorf("invalid credential value found, %T", v)
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

func StateKeyHolderDID(ca base.Address, credentialServiceID extensioncurrency.ContractID, ha base.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyCredentialServicePrefix(ca, credentialServiceID), ha.String(), HolderDIDSuffix)
}

func checkExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return base.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkNotExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return base.NewBaseOperationProcessReasonError("state, %q exists", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	var st base.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, base.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = currency.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currency.CurrencyID, getStateFunc base.GetStateFunc) (extensioncurrency.CurrencyPolicy, error) {
	var policy extensioncurrency.CurrencyPolicy
	switch i, found, err := getStateFunc(extensioncurrency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return extensioncurrency.CurrencyPolicy{}, err
	case !found:
		return extensioncurrency.CurrencyPolicy{}, base.NewBaseOperationProcessReasonError("currency not found, %v", cid)
	default:
		currencydesign, ok := i.Value().(extensioncurrency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return extensioncurrency.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", i.Value())
		}
		policy = currencydesign.CurrencyDesign.Policy()
	}
	return policy, nil
}
