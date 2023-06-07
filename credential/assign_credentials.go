package credential

import (
	"fmt"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

type CredentialItem interface {
	util.Byter
	util.IsValider
	Currency() currencybase.CurrencyID
}

var (
	AssignCredentialsFactHint = hint.MustNewHint("mitum-credential-assign-credentials-operation-fact-v0.0.1")
	AssignCredentialsHint     = hint.MustNewHint("mitum-credential-assign-credentials-operation-v0.0.1")
)

var MaxAssignCredentialsItems uint = 10

type AssignCredentialsFact struct {
	base.BaseFact
	sender base.Address
	items  []AssignCredentialsItem
}

func NewAssignCredentialsFact(token []byte, sender base.Address, items []AssignCredentialsItem) AssignCredentialsFact {
	bf := base.NewBaseFact(AssignCredentialsFactHint, token)
	fact := AssignCredentialsFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AssignCredentialsFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AssignCredentialsFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AssignCredentialsFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact AssignCredentialsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currencybase.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxAssignCredentialsItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxAssignCredentialsItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return err
		}

		if it.contract.Equal(fact.sender) {
			return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
		}

		k := fmt.Sprintf("%s-%s-%s", it.contract, it.credentialServiceID, it.id)

		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicate credential id found, %s", k)
		}

		founds[k] = struct{}{}
	}

	return nil
}

func (fact AssignCredentialsFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AssignCredentialsFact) Sender() base.Address {
	return fact.sender
}

func (fact AssignCredentialsFact) Items() []AssignCredentialsItem {
	return fact.items
}

func (fact AssignCredentialsFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	adrMap := make(map[string]struct{})
	for i := range fact.items {
		for j := range fact.items[i].Addresses() {
			if _, found := adrMap[fact.items[i].Addresses()[j].String()]; !found {
				adrMap[fact.items[i].Addresses()[j].String()] = struct{}{}
				as = append(as, fact.items[i].Addresses()[j])
			}
		}
	}
	as = append(as, fact.sender)

	return as, nil
}

type AssignCredentials struct {
	currencybase.BaseOperation
}

func NewAssignCredentials(fact AssignCredentialsFact) (AssignCredentials, error) {
	return AssignCredentials{BaseOperation: currencybase.NewBaseOperation(AssignCredentialsHint, fact)}, nil
}

func (op *AssignCredentials) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
