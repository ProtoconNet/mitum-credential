package types

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var DesignHint = hint.MustNewHint("mitum-credential-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	credentialServiceID currencytypes.ContractID
	policy              Policy
}

func NewDesign(credentialServiceID currencytypes.ContractID, policy Policy) Design {
	return Design{
		BaseHinter:          hint.NewBaseHinter(DesignHint),
		credentialServiceID: credentialServiceID,
		policy:              policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.credentialServiceID,
		de.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %v", err)
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		de.credentialServiceID.Bytes(),
		de.policy.Bytes(),
	)
}

func (de Design) ServiceID() currencytypes.ContractID {
	return de.credentialServiceID
}

func (de Design) Policy() Policy {
	return de.policy
}
