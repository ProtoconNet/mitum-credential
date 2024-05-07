package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (de *Design) unpack(enc encoder.Encoder, ht hint.Hint, bPcy []byte) error {
	e := util.StringError("unpack Design")

	de.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(bPcy); err != nil {
		return e.Wrap(err)
	} else if po, ok := hinter.(Policy); !ok {
		return e.Wrap(common.ErrTypeMismatch.Wrap(errors.Errorf("expected Policy, not %T", hinter)))
	} else {
		de.policy = po
	}
	if err := de.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}
