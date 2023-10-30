package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (h *Holder) unpack(enc encoder.Encoder, ht hint.Hint, adr string, count uint64) error {
	e := util.StringError("failed to unpack of Holder")

	h.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(adr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		h.address = a
	}
	if err := h.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	h.credentialCount = count

	if err := h.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint, tmplIDs []string, bHolders []byte, count uint64) error {
	e := util.StringError("failed to unpack of Policy")

	po.BaseHinter = hint.NewBaseHinter(ht)
	po.templateIDs = tmplIDs

	hds, err := enc.DecodeSlice(bHolders)
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
	po.credentialCount = count
	if err := po.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}
