package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var DesignHint = hint.MustNewHint("mitum-credential-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	serviceID ServiceID
	policy    Policy
}

func NewDesign(serviceID ServiceID, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		serviceID:  serviceID,
		policy:     policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.serviceID,
		de.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %v", err)
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		de.serviceID.Bytes(),
		de.policy.Bytes(),
	)
}

func (de Design) ServiceID() ServiceID {
	return de.serviceID
}

func (de Design) Policy() Policy {
	return de.policy
}
