package credential

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateCredentialServiceFact) unpack(enc encoder.Encoder, sa, ca, csid string, cid string) error {
	e := util.StringError("failed to unmarshal CreateCredentialServiceFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.contract = a
	}

	fact.credentialServiceID = currencytypes.ContractID(csid)
	fact.currency = currencytypes.CurrencyID(cid)

	return nil
}
