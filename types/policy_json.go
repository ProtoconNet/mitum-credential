package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Templates       []Uint256 `json:"templates"`
	Holders         []Holder  `json:"holders"`
	CredentialCount uint64    `json:"credential_count"`
}

func (po Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter:      po.BaseHinter,
		Templates:       po.templates,
		Holders:         po.holders,
		CredentialCount: po.credentialCount,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint            hint.Hint       `json:"_hint"`
	Templates       []string        `json:"templates"`
	Holders         json.RawMessage `json:"holders"`
	CredentialCount uint64          `json:"credential_count"`
}

func (po *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Policy")

	var upo PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	return po.unpack(enc, upo.Hint, upo.Templates, upo.Holders, upo.CredentialCount)
}
