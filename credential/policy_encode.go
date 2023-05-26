package credential

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint, bts []string, bhd []byte, ccount uint64) error {
	e := util.StringErrorFunc("failed to decode bson of Policy")

	po.BaseHinter = hint.NewBaseHinter(ht)

	templates := make([]Uint256, len(bts))
	for i := range bts {
		t, err := NewUint256FromString(bts[i])
		if err != nil {
			return e(err, "")
		}
		templates[i] = t
	}
	po.templates = templates

	hds, err := enc.DecodeSlice(bhd)
	if err != nil {
		return e(err, "")
	}

	holders := make([]Holder, len(hds))
	for i := range hds {
		j, ok := hds[i].(Holder)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected Holder, not %T", hds[i]), "")
		}

		holders[i] = j
	}
	po.holders = holders

	return nil
}
