package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateCredentialServiceFact) unpack(enc encoder.Encoder, sa, ca, csid string, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal CreateCredentialServiceFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.contract = a
	}

	fact.credentialServiceID = currencybase.ContractID(csid)
	fact.currency = currencybase.CurrencyID(cid)

	return nil
}
