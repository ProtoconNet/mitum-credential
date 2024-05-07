package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *RevokeFact) unpack(enc encoder.Encoder, sAdr string, bItm []byte) error {
	switch a, err := base.DecodeAddress(sAdr, enc); {
	case err != nil:
		return err
	default:
		fact.sender = a
	}

	hItm, err := enc.DecodeSlice(bItm)
	if err != nil {
		return err
	}

	items := make([]RevokeItem, len(hItm))
	for i := range hItm {
		j, ok := hItm[i].(RevokeItem)
		if !ok {
			return common.ErrTypeMismatch.Wrap(errors.Errorf("expected RevokeItem, not %T", hItm[i]))
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
