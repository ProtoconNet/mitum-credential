package credential

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	DesignHint = hint.MustNewHint("mitum-credential-design-v0.0.1")
)

type Design struct {
	hint.BaseHinter
	creditID extensioncurrency.ContractID
	policy   Policy
}

func NewDesign(creditID extensioncurrency.ContractID, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		creditID:   creditID,
		policy:     policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.creditID,
		de.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %w", err)
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		de.creditID.Bytes(),
		de.policy.Bytes(),
	)
}

func (de Design) CreditID() extensioncurrency.ContractID {
	return de.creditID
}

func (de Design) Policy() Policy {
	return de.policy
}
