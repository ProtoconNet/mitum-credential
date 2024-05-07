package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *Template) unpack(enc encoder.Encoder, ht hint.Hint,
	tmplID string,
	tmplName, svcDate, expDate string,
	share, audit bool,
	dpName, subjKey, desc, creator string,
) error {
	e := util.StringError("unpack Template")

	t.BaseHinter = hint.NewBaseHinter(ht)
	t.templateID = tmplID
	t.templateName = tmplName
	t.serviceDate = Date(svcDate)
	t.expirationDate = Date(expDate)
	t.templateShare = Bool(share)
	t.multiAudit = Bool(audit)
	t.displayName = dpName
	t.subjectKey = subjKey
	t.description = desc

	switch a, err := base.DecodeAddress(creator, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		t.creator = a
	}
	if err := t.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}
