package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (de *Design) unpack(enc encoder.Encoder, ht hint.Hint, id string, bpo []byte) error {
	e := util.StringError("failed to unpack of Design")

	de.BaseHinter = hint.NewBaseHinter(ht)
	de.serviceID = ServiceID(id)

	if hinter, err := enc.Decode(bpo); err != nil {
		return e.Wrap(err)
	} else if po, ok := hinter.(Policy); !ok {
		return e.Wrap(errors.Errorf("expected Policy, not %T", hinter))
	} else {
		de.policy = po
	}

	return nil
}
