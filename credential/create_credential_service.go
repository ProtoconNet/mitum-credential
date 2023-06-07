package credential

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CreateCredentialServiceFactHint = hint.MustNewHint("mitum-credential-create-credential-service-operation-fact-v0.0.1")
	CreateCredentialServiceHint     = hint.MustNewHint("mitum-credential-create-credential-service-operation-v0.0.1")
)

type CreateCredentialServiceFact struct {
	base.BaseFact
	sender              base.Address
	contract            base.Address
	credentialServiceID currencybase.ContractID
	currency            currencybase.CurrencyID
}

func NewCreateCredentialServiceFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	credentialServiceID currencybase.ContractID,
	currency currencybase.CurrencyID,
) CreateCredentialServiceFact {
	bf := base.NewBaseFact(CreateCredentialServiceFactHint, token)
	fact := CreateCredentialServiceFact{
		BaseFact:            bf,
		sender:              sender,
		contract:            contract,
		credentialServiceID: credentialServiceID,
		currency:            currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateCredentialServiceFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateCredentialServiceFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateCredentialServiceFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.credentialServiceID.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact CreateCredentialServiceFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currencybase.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.credentialServiceID, fact.contract, fact.currency); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	return nil
}

func (fact CreateCredentialServiceFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateCredentialServiceFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateCredentialServiceFact) Contract() base.Address {
	return fact.contract
}

func (fact CreateCredentialServiceFact) CredentialServiceID() currencybase.ContractID {
	return fact.credentialServiceID
}

func (fact CreateCredentialServiceFact) Currency() currencybase.CurrencyID {
	return fact.currency
}

func (fact CreateCredentialServiceFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type CreateCredentialService struct {
	currencybase.BaseOperation
}

func NewCreateCredentialService(fact CreateCredentialServiceFact) (CreateCredentialService, error) {
	return CreateCredentialService{BaseOperation: currencybase.NewBaseOperation(CreateCredentialServiceHint, fact)}, nil
}

func (op *CreateCredentialService) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
