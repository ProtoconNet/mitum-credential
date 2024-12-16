package credential

import (
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *IssueItem) unpack(enc encoder.Encoder, ht hint.Hint,
	cAdr, hAdr, tmplID string,
	id string,
	val string,
	vFrom, vUntil uint64,
	did, cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	it.credentialID = id
	it.value = val
	it.did = did
	it.currency = ctypes.CurrencyID(cid)

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
	it.validFrom = vFrom
	it.validUntil = vUntil

	return nil
}
