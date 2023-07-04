package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint, ts []uint64, bhd []byte, ccount uint64) error {
	e := util.StringError("failed to decode bson of Policy")

	po.BaseHinter = hint.NewBaseHinter(ht)
	po.templateIDs = ts

	hds, err := enc.DecodeSlice(bhd)
	if err != nil {
		return e.Wrap(err)
	}

	holders := make([]Holder, len(hds))
	for i := range hds {
		j, ok := hds[i].(Holder)
		if !ok {
			return e.Wrap(errors.Errorf("expected Holder, not %T", hds[i]))
		}

		holders[i] = j
	}
	po.holders = holders
	po.credentialCount = ccount

	return nil
}
