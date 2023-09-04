package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RevokeItem) unpack(enc encoder.Encoder, ht hint.Hint,
	ca, sid, hd, tid string,
	id, cid string,
) error {
	e := util.StringError("failed to unmarshal RevokeItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.serviceID = types.ServiceID(sid)
	it.id = id
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(hd, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.holder = a
	}

	it.templateID = tid

	return nil
}
