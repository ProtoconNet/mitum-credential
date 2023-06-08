package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RevokeCredentialsItem) unpack(enc encoder.Encoder, ht hint.Hint,
	ca, csid, hd, tid, id, cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal RevokeCredentialsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.credentialServiceID = currencybase.ContractID(csid)
	it.id = id
	it.currency = currencybase.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(hd, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.holder = a
	}

	templateid, err := NewUint256FromString(tid)
	if err != nil {
		return e(err, "")
	}
	it.templateID = templateid

	return nil
}
