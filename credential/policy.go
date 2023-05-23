package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	PolicyHint = hint.MustNewHint("mitum-credential-policy-v0.0.1")
)

type Policy struct {
	hint.BaseHinter
	templates       []Uint256
	holders         []base.Address
	credentialCount uint64
}

func NewPolicy(templates []Uint256, holders []base.Address, credentialCount uint64) Policy {
	return Policy{
		BaseHinter:      hint.NewBaseHinter(PolicyHint),
		templates:       templates,
		holders:         holders,
		credentialCount: credentialCount,
	}
}

func (po Policy) Bytes() []byte {
	ts := make([][]byte, len(po.templates))
	for i, t := range po.templates {
		ts[i] = t.Bytes()
	}

	hs := make([][]byte, len(po.holders))
	for i, h := range po.holders {
		hs[i] = h.Bytes()
	}

	return util.ConcatBytesSlice(
		util.ConcatBytesSlice(ts...),
		util.ConcatBytesSlice(hs...),
		util.Uint64ToBytes(po.credentialCount),
	)
}

func (po Policy) IsValid([]byte) error {
	e := util.StringErrorFunc("invalid credential policy")

	if err := util.CheckIsValiders(nil, false, po.BaseHinter); err != nil {
		return e(err, "")
	}

	for _, t := range po.templates {
		if err := t.IsValid(nil); err != nil {
			return e(err, "")
		}
	}

	for _, h := range po.holders {
		if err := h.IsValid(nil); err != nil {
			return e(err, "")
		}
	}

	return nil
}

func (po Policy) Templates() []Uint256 {
	return po.templates
}

func (po Policy) Holders() []base.Address {
	return po.holders
}

func (po Policy) CredentialCount() uint64 {
	return po.credentialCount
}
