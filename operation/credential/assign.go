package credential

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

type CredentialItem interface {
	util.Byter
	util.IsValider
	Currency() currencytypes.CurrencyID
}

var (
	AssignFactHint = hint.MustNewHint("mitum-credential-assign-operation-fact-v0.0.1")
	AssignHint     = hint.MustNewHint("mitum-credential-assign-operation-v0.0.1")
)

var MaxAssignItems uint = 10

type AssignFact struct {
	base.BaseFact
	sender base.Address
	items  []AssignItem
}

func NewAssignFact(token []byte, sender base.Address, items []AssignItem) AssignFact {
	bf := base.NewBaseFact(AssignFactHint, token)
	fact := AssignFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AssignFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AssignFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AssignFact) Bytes() []byte {
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

func (fact AssignFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxAssignItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxAssignItems)
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

		k := fmt.Sprintf("%s-%s-%s", it.contract, it.serviceID, it.id)

		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicate credential id found, %s", k)
		}

		founds[k] = struct{}{}
	}

	return nil
}

func (fact AssignFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AssignFact) Sender() base.Address {
	return fact.sender
}

func (fact AssignFact) Items() []AssignItem {
	return fact.items
}

func (fact AssignFact) Addresses() ([]base.Address, error) {
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

type Assign struct {
	common.BaseOperation
}

func NewAssign(fact AssignFact) (Assign, error) {
	return Assign{BaseOperation: common.NewBaseOperation(AssignHint, fact)}, nil
}

func (op *Assign) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
