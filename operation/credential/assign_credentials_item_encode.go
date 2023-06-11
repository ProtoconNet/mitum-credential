package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *AssignCredentialsItem) unpack(enc encoder.Encoder, ht hint.Hint,
	ca, csid, hd, tid, id, v, vf, vu, did, cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal AssignCredentialsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.credentialServiceID = currencytypes.ContractID(csid)
	it.id = id
	it.value = v
	it.did = did
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

	validfrom, err := types.NewUint256FromString(vf)
	if err != nil {
		return e(err, "")
	}
	it.validfrom = validfrom

	validuntil, err := types.NewUint256FromString(vu)
	if err != nil {
		return e(err, "")
	}
	it.validuntil = validuntil

	return nil
}
