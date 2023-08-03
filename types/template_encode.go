package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *Template) unpack(enc encoder.Encoder, ht hint.Hint,
	tid string,
	name, service, expire string,
	share, audit bool,
	dname, sjk, desc, creator string,
) error {
	e := util.StringError("failed to unpack of Template")

	t.BaseHinter = hint.NewBaseHinter(ht)
	t.templateID = tid
	t.templateName = name
	t.serviceDate = Date(service)
	t.expirationDate = Date(expire)
	t.templateShare = Bool(share)
	t.multiAudit = Bool(audit)
	t.displayName = dname
	t.subjectKey = sjk
	t.description = desc

	switch a, err := base.DecodeAddress(creator, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		t.creator = a
	}

	return nil
}
