package credential

import (
	"github.com/ProtoconNet/mitum-credential/types"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AddTemplateFact) unpack(enc encoder.Encoder,
	sAdr, cAdr, tmplID string,
	tmplName, svcDate, expDate string,
	tmplShr, ma bool,
	dpName, subjKey, desc, crAdr, cid string,
) error {
	fact.templateName = tmplName
	fact.serviceDate = types.Date(svcDate)
	fact.expirationDate = types.Date(expDate)
	fact.templateShare = types.Bool(tmplShr)
	fact.multiAudit = types.Bool(ma)
	fact.displayName = dpName
	fact.subjectKey = subjKey
	fact.description = desc
	fact.currency = ctypes.CurrencyID(cid)
	fact.templateID = tmplID

	switch a, err := base.DecodeAddress(sAdr, enc); {
	case err != nil:
		return err
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return err
	default:
		fact.contract = a
	}

	switch a, err := base.DecodeAddress(crAdr, enc); {
	case err != nil:
		return err
	default:
		fact.creator = a
	}

	return nil
}
