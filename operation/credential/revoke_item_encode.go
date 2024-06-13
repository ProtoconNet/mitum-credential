package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RevokeItem) unpack(enc encoder.Encoder, ht hint.Hint,
	cAdr, hAdr, tmplID string,
	id, cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	it.credentialID = id
	it.currency = types.CurrencyID(cid)

	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(hAdr, enc); {
	case err != nil:
		return err
	default:
		it.holder = a
	}

	it.templateID = tmplID

	return nil
}
