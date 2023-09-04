package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AddTemplateFact) unpack(enc encoder.Encoder,
	sa, ca, sid, tid string,
	tname, sd, ed string,
	ts, ma bool,
	dn, sk, desc, cr, cid string,
) error {
	e := util.StringError("failed to unmarshal AddTemplateFact")

	fact.serviceID = types.ServiceID(sid)
	fact.templateName = tname
	fact.serviceDate = types.Date(sd)
	fact.expirationDate = types.Date(ed)
	fact.templateShare = types.Bool(ts)
	fact.multiAudit = types.Bool(ma)
	fact.displayName = dn
	fact.subjectKey = sk
	fact.description = desc
	fact.currency = currencytypes.CurrencyID(cid)
	fact.templateID = tid

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

	switch a, err := base.DecodeAddress(cr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.creator = a
	}

	return nil
}
