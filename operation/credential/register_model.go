package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RegisterModelFactHint = hint.MustNewHint("mitum-credential-register-model-operation-fact-v0.0.1")
	RegisterModelHint     = hint.MustNewHint("mitum-credential-register-model-operation-v0.0.1")
)

type RegisterModelFact struct {
	base.BaseFact
	sender   base.Address
	contract base.Address
	currency types.CurrencyID
}

func NewRegisterModelFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	currency types.CurrencyID,
) RegisterModelFact {
	bf := base.NewBaseFact(RegisterModelFactHint, token)
	fact := RegisterModelFact{
		BaseFact: bf,
		sender:   sender,
		contract: contract,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RegisterModelFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RegisterModelFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RegisterModelFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact RegisterModelFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := util.CheckIsValiders(
		nil,
		false,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("contract address is same with sender, %q", fact.sender)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact RegisterModelFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RegisterModelFact) Sender() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) Contract() base.Address {
	return fact.contract
}

func (fact RegisterModelFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact RegisterModelFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

func (fact RegisterModelFact) FeeBase() map[types.CurrencyID][]common.Big {
	required := make(map[types.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact RegisterModelFact) FeePayer() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) FactUser() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) Signer() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) InActiveContractOwnerHandlerOnly() [][2]base.Address {
	return [][2]base.Address{{fact.contract, fact.sender}}
}

type RegisterModel struct {
	extras.ExtendedOperation
}

func NewRegisterModel(fact RegisterModelFact) RegisterModel {
	return RegisterModel{
		ExtendedOperation: extras.NewExtendedOperation(RegisterModelHint, fact),
	}
}
