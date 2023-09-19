package credential

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RevokeItem) unpack(enc encoder.Encoder, ht hint.Hint,
	cAdr, hAdr, tmplID string,
	id, cid string,
) error {
	e := util.StringError("failed to unmarshal RevokeItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.id = id
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(hAdr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.holder = a
	}

	it.templateID = tmplID

	return nil
}
