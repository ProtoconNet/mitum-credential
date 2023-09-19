package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *Credential) unpack(enc encoder.Encoder, ht hint.Hint,
	holder, tmplID string,
	id, v string,
	vFrom, vUntil uint64,
	did string,
) error {
	e := util.StringError("failed to unpack of Credential")

	t.BaseHinter = hint.NewBaseHinter(ht)
	t.id = id
	t.value = v
	t.did = did

	switch a, err := base.DecodeAddress(holder, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		t.holder = a
	}

	t.templateID = tmplID
	t.validFrom = vFrom
	t.validUntil = vUntil

	return nil
}
