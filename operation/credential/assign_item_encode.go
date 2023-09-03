package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *AssignItem) unpack(enc encoder.Encoder, ht hint.Hint,
	ca, sid, hd, tid string,
	id string,
	v string,
	vf, vu uint64,
	did, cid string,
) error {
	e := util.StringError("failed to unmarshal AssignItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.serviceID = types.ServiceID(sid)
	it.id = id
	it.value = v
	it.did = did
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
	it.validfrom = vf
	it.validuntil = vu

	return nil
}
