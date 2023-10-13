package credential

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CreateServiceFactHint = hint.MustNewHint("mitum-credential-create-service-operation-fact-v0.0.1")
	CreateServiceHint     = hint.MustNewHint("mitum-credential-create-service-operation-v0.0.1")
)

type CreateServiceFact struct {
	base.BaseFact
	sender   base.Address
	contract base.Address
	currency currencytypes.CurrencyID
}

func NewCreateServiceFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	currency currencytypes.CurrencyID,
) CreateServiceFact {
	bf := base.NewBaseFact(CreateServiceFactHint, token)
	fact := CreateServiceFact{
		BaseFact: bf,
		sender:   sender,
		contract: contract,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateServiceFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateServiceFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateServiceFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact CreateServiceFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(
		nil,
		false,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact CreateServiceFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateServiceFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateServiceFact) Contract() base.Address {
	return fact.contract
}

func (fact CreateServiceFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CreateServiceFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type CreateService struct {
	common.BaseOperation
}

func NewCreateService(fact CreateServiceFact) (CreateService, error) {
	return CreateService{BaseOperation: common.NewBaseOperation(CreateServiceHint, fact)}, nil
}
