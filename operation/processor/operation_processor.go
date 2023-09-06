package processor

import (
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender             currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency           currencytypes.DuplicationType = "currency"
	DuplicationTypeContractCredential currencytypes.DuplicationType = "contract-credential"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var did string
	var didtype currencytypes.DuplicationType
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.CreateAccount:
		fact, ok := t.Fact().(currency.CreateAccountFact)
		if !ok {
			return errors.Errorf("expected CreateAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		did = fact.Target().String()
		didtype = DuplicationTypeSender
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case extensioncurrency.CreateContractAccount:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
	case extensioncurrency.Withdraw:
		fact, ok := t.Fact().(extensioncurrency.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case credential.CreateService:
		fact, ok := t.Fact().(credential.CreateServiceFact)
		if !ok {
			return errors.Errorf("expected CreateServiceFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case credential.AddTemplate:
		fact, ok := t.Fact().(credential.AddTemplateFact)
		if !ok {
			return errors.Errorf("expected AddTemplateFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case credential.Assign:
		fact, ok := t.Fact().(credential.AssignFact)
		if !ok {
			return errors.Errorf("expected AssignFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case credential.Revoke:
		fact, ok := t.Fact().(credential.AssignFact)
		if !ok {
			return errors.Errorf("expected Revoke, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected CurrencyRegisterFact, not %T", t.Fact())
		}
		did = fact.Currency().Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		did = fact.Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.Mint:
	default:
		return nil
	}

	if len(did) > 0 {
		if _, found := opr.Duplicated[did]; found {
			switch didtype {
			case DuplicationTypeSender:
				return errors.Errorf("violates only one sender in proposal")
			case DuplicationTypeCurrency:
				return errors.Errorf("duplicate currency id, %q found in proposal", did)
			default:
				return errors.Errorf("violates duplication in proposal")
			}
		}

		opr.Duplicated[did] = didtype
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func GetNewProcessor(opr *currencyprocessor.OperationProcessor, op base.Operation) (base.OperationProcessor, bool, error) {
	switch i, err := opr.GetNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccount,
		currency.UpdateKey,
		currency.Transfer,
		extensioncurrency.CreateContractAccount,
		extensioncurrency.Withdraw,
		currency.RegisterCurrency,
		currency.UpdateCurrency,
		currency.Mint,
		credential.CreateService,
		credential.AddTemplate,
		credential.Assign,
		credential.Revoke:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
