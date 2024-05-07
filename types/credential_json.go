package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CredentialJSONMarshaler struct {
	hint.BaseHinter
	Holder     base.Address `json:"holder"`
	TemplateID string       `json:"template_id"`
	ID         string       `json:"id"`
	Value      string       `json:"value"`
	ValidFrom  uint64       `json:"valid_from"`
	ValidUntil uint64       `json:"valid_until"`
	DID        string       `json:"did"`
}

func (c Credential) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CredentialJSONMarshaler{
		BaseHinter: c.BaseHinter,
		Holder:     c.holder,
		TemplateID: c.templateID,
		ID:         c.id,
		Value:      c.value,
		ValidFrom:  c.validFrom,
		ValidUntil: c.validUntil,
		DID:        c.did,
	})
}

type CredentialJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Holder     string    `json:"holder"`
	TemplateID string    `json:"template_id"`
	ID         string    `json:"id"`
	Value      string    `json:"value"`
	ValidFrom  uint64    `json:"valid_from"`
	ValidUntil uint64    `json:"valid_until"`
	DID        string    `json:"did"`
}

func (c *Credential) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("decode json of Credential")

	var u CredentialJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return c.unpack(enc, u.Hint,
		u.Holder,
		u.TemplateID,
		u.ID,
		u.Value,
		u.ValidFrom,
		u.ValidUntil,
		u.DID,
	)
}
