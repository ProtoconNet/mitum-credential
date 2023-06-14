package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *Credential) unpack(enc encoder.Encoder, ht hint.Hint,
	hd, tid, id, v, vf, vu, did string,
) error {
	e := util.StringErrorFunc("failed to decode bson of Credential")

	t.BaseHinter = hint.NewBaseHinter(ht)
	t.id = id
	t.value = v
	t.did = did

	switch a, err := base.DecodeAddress(hd, enc); {
	case err != nil:
		return e(err, "")
	default:
		t.holder = a
	}

	templateid, err := NewUint256FromString(tid)
	if err != nil {
		return e(err, "")
	}
	t.templateID = templateid

	validfrom, err := NewUint256FromString(vf)
	if err != nil {
		return e(err, "")
	}
	t.validfrom = validfrom

	validuntil, err := NewUint256FromString(vu)
	if err != nil {
		return e(err, "")
	}
	t.validuntil = validuntil

	return nil
}
