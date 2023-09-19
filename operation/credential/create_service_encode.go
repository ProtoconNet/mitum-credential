package credential

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateServiceFact) unpack(enc encoder.Encoder, sAdr, cAdr, cid string) error {
	e := util.StringError("failed to unmarshal CreateServiceFact")

	switch a, err := base.DecodeAddress(sAdr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.contract = a
	}

	fact.currency = currencytypes.CurrencyID(cid)

	return nil
}
