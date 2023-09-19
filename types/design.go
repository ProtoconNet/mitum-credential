package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var DesignHint = hint.MustNewHint("mitum-credential-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	policy Policy
}

func NewDesign(policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		policy:     policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %v", err)
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		de.policy.Bytes(),
	)
}

func (de Design) Policy() Policy {
	return de.policy
}
