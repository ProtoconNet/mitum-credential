package credential

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RevokeFactHint = hint.MustNewHint("mitum-credential-revoke-operation-fact-v0.0.1")
	RevokeHint     = hint.MustNewHint("mitum-credential-revoke-operation-v0.0.1")
)

var MaxRevokeItems uint = 10

type RevokeFact struct {
	base.BaseFact
	sender base.Address
	items  []RevokeItem
}

func NewRevokeFact(token []byte, sender base.Address, items []RevokeItem) RevokeFact {
	bf := base.NewBaseFact(RevokeFactHint, token)
	fact := RevokeFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RevokeFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RevokeFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RevokeFact) Bytes() []byte {
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

func (fact RevokeFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("empty items")))
	} else if n > int(MaxRevokeItems) {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("items, %d over max, %d", n, MaxRevokeItems)))
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
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		k := fmt.Sprintf("%s-%s", it.contract, it.id)

		if _, found := founds[k]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("credential id %v for template %v in contract account %v", it.ID(), it.TemplateID(), it.Contract())))
		}

		founds[it.ID()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact RevokeFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RevokeFact) Sender() base.Address {
	return fact.sender
}

func (fact RevokeFact) Items() []RevokeItem {
	return fact.items
}

func (fact RevokeFact) Addresses() ([]base.Address, error) {
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

type Revoke struct {
	common.BaseOperation
}

func NewRevoke(fact RevokeFact) Revoke {
	return Revoke{BaseOperation: common.NewBaseOperation(RevokeHint, fact)}
}
