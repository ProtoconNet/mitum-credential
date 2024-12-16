package credential

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

type CredentialItem interface {
	util.Byter
	util.IsValider
	Currency() types.CurrencyID
}

var (
	IssueFactHint = hint.MustNewHint("mitum-credential-issue-operation-fact-v0.0.1")
	IssueHint     = hint.MustNewHint("mitum-credential-issue-operation-v0.0.1")
)

var MaxIssueItems uint = 1000

type IssueFact struct {
	base.BaseFact
	sender base.Address
	items  []IssueItem
}

func NewIssueFact(token []byte, sender base.Address, items []IssueItem) IssueFact {
	bf := base.NewBaseFact(IssueFactHint, token)
	fact := IssueFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact IssueFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact IssueFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact IssueFact) Bytes() []byte {
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

func (fact IssueFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if n := len(fact.items); n < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items")))
	} else if n > int(MaxIssueItems) {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("items, %d over max, %d", n, MaxIssueItems)))
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if it.contract.Equal(fact.sender) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		k := fmt.Sprintf("%s-%s", it.contract, it.credentialID)

		if _, found := founds[k]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("credential id %v for template %v in contract account %v", it.CredentialID(), it.TemplateID(), it.Contract())))
		}

		founds[k] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact IssueFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact IssueFact) Sender() base.Address {
	return fact.sender
}

func (fact IssueFact) Items() []IssueItem {
	return fact.items
}

func (fact IssueFact) Addresses() ([]base.Address, error) {
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

func (fact IssueFact) FeeBase() map[types.CurrencyID][]common.Big {
	required := make(map[types.CurrencyID][]common.Big)

	for i := range fact.items {
		zeroBig := common.ZeroBig
		cid := fact.items[i].Currency()
		var amsTemp []common.Big
		if ams, found := required[cid]; found {
			ams = append(ams, zeroBig)
			required[cid] = ams
		} else {
			amsTemp = append(amsTemp, zeroBig)
			required[cid] = amsTemp
		}
	}

	return required
}

func (fact IssueFact) FeePayer() base.Address {
	return fact.sender
}

func (fact IssueFact) FactUser() base.Address {
	return fact.sender
}

func (fact IssueFact) Signer() base.Address {
	return fact.sender
}

func (fact IssueFact) ActiveContractOwnerHandlerOnly() [][2]base.Address {
	var arr [][2]base.Address
	for i := range fact.items {
		arr = append(arr, [2]base.Address{fact.items[i].contract, fact.sender})
	}
	return arr
}

type Issue struct {
	extras.ExtendedOperation
}

func NewIssue(fact IssueFact) Issue {
	return Issue{
		ExtendedOperation: extras.NewExtendedOperation(IssueHint, fact),
	}
}
