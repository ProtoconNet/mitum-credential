package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AddTemplateFact) unpack(enc encoder.Encoder,
	sa, ca, csid, tid, tname, sd, ed string,
	ts, ma bool,
	dn, sk, desc, cr, cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal AddTemplateFact")

	fact.credentialServiceID = currencybase.ContractID(csid)
	fact.templateName = tname
	fact.serviceDate = Date(sd)
	fact.expirationDate = Date(ed)
	fact.templateShare = Bool(ts)
	fact.multiAudit = Bool(ma)
	fact.displayName = dn
	fact.subjectKey = sk
	fact.description = desc
	fact.currency = currencybase.CurrencyID(cid)

	templateid, err := NewUint256FromString(tid)
	if err != nil {
		return e(err, "")
	}
	fact.templateID = templateid

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

	switch a, err := base.DecodeAddress(cr, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.creator = a
	}

	return nil
}
