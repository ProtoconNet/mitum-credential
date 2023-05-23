package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/holiman/uint256"
)

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint, bts, bhs []string, ccount uint64) error {
	e := util.StringErrorFunc("failed to decode bson of Policy")

	po.BaseHinter = hint.NewBaseHinter(ht)

	templates := make([]Uint256, len(bts))
	for i := range bts {
		t, err := uint256.FromHex(bts[i])
		if err != nil {
			return e(err, "")
		}
		templates[i] = Uint256{
			n: *t,
		}
	}
	po.templates = templates

	holders := make([]base.Address, len(bhs))
	for i := range bhs {
		a, err := base.DecodeAddress(bhs[i], enc)
		if err != nil {
			return e(err, "")
		}
		holders[i] = a
	}
	po.holders = holders

	return nil
}
