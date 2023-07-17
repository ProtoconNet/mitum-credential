package types

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type HolderJSONMarshaler struct {
	hint.BaseHinter
	Address         base.Address `json:"address"`
	CredentialCount uint64       `json:"credential_count"`
}

func (h Holder) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(HolderJSONMarshaler{
		BaseHinter:      h.BaseHinter,
		Address:         h.address,
		CredentialCount: h.credentialCount,
	})
}

type HolderJSONUnmarshaler struct {
	Hint            hint.Hint `json:"_hint"`
	Address         string    `json:"address"`
	CredentialCount uint64    `json:"credential_count"`
}

func (h *Holder) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Holder")

	var uho HolderJSONUnmarshaler
	if err := enc.Unmarshal(b, &uho); err != nil {
		return e.Wrap(err)
	}

	return h.unpack(enc, uho.Hint, uho.Address, uho.CredentialCount)
}

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Templates       []uint64 `json:"templates"`
	Holders         []Holder `json:"holders"`
	CredentialCount uint64   `json:"credential_count"`
}

func (po Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter:      po.BaseHinter,
		Templates:       po.templateIDs,
		Holders:         po.holders,
		CredentialCount: po.credentialCount,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint            hint.Hint       `json:"_hint"`
	Templates       []uint64        `json:"templates"`
	Holders         json.RawMessage `json:"holders"`
	CredentialCount uint64          `json:"credential_count"`
}

func (po *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Policy")

	var upo PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	return po.unpack(enc, upo.Hint, upo.Templates, upo.Holders, upo.CredentialCount)
}
