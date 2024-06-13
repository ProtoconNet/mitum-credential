package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (c *Credential) unpack(enc encoder.Encoder, ht hint.Hint,
	holder, tmplID string,
	id, v string,
	vFrom, vUntil uint64,
	did string,
) error {
	e := util.StringError("unpack Credential")

	c.BaseHinter = hint.NewBaseHinter(ht)
	c.credentialID = id
	c.value = v
	c.did = did

	switch a, err := base.DecodeAddress(holder, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		c.holder = a
	}

	c.templateID = tmplID
	c.validFrom = vFrom
	c.validUntil = vUntil
	if err := c.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}
