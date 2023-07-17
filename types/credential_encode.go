package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *Credential) unpack(enc encoder.Encoder, ht hint.Hint,
	hd string,
	tid uint64,
	id, v string,
	vf, vu uint64,
	did string,
) error {
	e := util.StringError("failed to unpack of Credential")

	t.BaseHinter = hint.NewBaseHinter(ht)
	t.id = id
	t.value = v
	t.did = did

	switch a, err := base.DecodeAddress(hd, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		t.holder = a
	}

	t.templateID = tid
	t.validFrom = vf
	t.validUntil = vu

	return nil
}
