package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CredentialJSONMarshaler struct {
	hint.BaseHinter
	Holder     base.Address `json:"holder"`
	TemplateID Uint256      `json:"templateid"`
	ID         string       `json:"id"`
	Value      string       `json:"value"`
	ValidFrom  Uint256      `json:"valid_from"`
	ValidUntil Uint256      `json:"valid_until"`
	DID        string       `json:"did"`
}

func (cd Credential) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CredentialJSONMarshaler{
		BaseHinter: cd.BaseHinter,
		Holder:     cd.holder,
		TemplateID: cd.templateID,
		ID:         cd.id,
		Value:      cd.value,
		ValidFrom:  cd.validfrom,
		ValidUntil: cd.validuntil,
	})
}

type CredentialJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Holder     string    `json:"holder"`
	TemplateID string    `json:"templateid"`
	ID         string    `json:"id"`
	Value      string    `json:"value"`
	ValidFrom  string    `json:"valid_from"`
	ValidUntil string    `json:"valid_until"`
	DID        string    `json:"did"`
}

func (cd *Credential) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Credential")

	var u CredentialJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return cd.unpack(enc, u.Hint,
		u.Holder,
		u.TemplateID,
		u.ID,
		u.Value,
		u.ValidFrom,
		u.ValidUntil,
		u.DID,
	)
}
