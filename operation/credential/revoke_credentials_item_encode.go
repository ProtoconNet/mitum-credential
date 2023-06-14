package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
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
	it.credentialServiceID = currencytypes.ContractID(csid)
	it.id = id
	it.currency = currencytypes.CurrencyID(cid)

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

	templateid, err := types.NewUint256FromString(tid)
	if err != nil {
		return e(err, "")
	}
	it.templateID = templateid

	return nil
}
