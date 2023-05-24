package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *Template) unpack(enc encoder.Encoder, ht hint.Hint,
	id, name, service, expire string,
	share, audit bool,
	dname, sjk, desc, creator string,
) error {
	e := util.StringErrorFunc("failed to decode bson of Template")

	t.BaseHinter = hint.NewBaseHinter(ht)
	t.templateName = name
	t.serviceDate = Date(service)
	t.expirationDate = Date(expire)
	t.templateShare = Bool(share)
	t.multiAudit = Bool(audit)
	t.displayName = dname
	t.subjectKey = sjk
	t.description = desc

	tid, err := NewUint256FromString(id)
	if err != nil {
		return e(err, "")
	}
	t.templateID = tid

	switch a, err := base.DecodeAddress(creator, enc); {
	case err != nil:
		return e(err, "")
	default:
		t.creator = a
	}

	return nil
}
