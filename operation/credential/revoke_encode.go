package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *RevokeFact) unpack(enc encoder.Encoder, sAdr string, bItm []byte) error {
	e := util.StringError("failed to unmarshal RevokeFact")

	switch a, err := base.DecodeAddress(sAdr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	hItm, err := enc.DecodeSlice(bItm)
	if err != nil {
		return e.Wrap(err)
	}

	items := make([]RevokeItem, len(hItm))
	for i := range hItm {
		j, ok := hItm[i].(RevokeItem)
		if !ok {
			return e.Wrap(errors.Errorf("expected RevokeItem, not %T", hItm[i]))
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
