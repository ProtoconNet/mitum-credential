package processor

import (
	"fmt"
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender     ctypes.DuplicationType = "sender"
	DuplicationTypeCurrency   ctypes.DuplicationType = "currency"
	DuplicationTypeContract   ctypes.DuplicationType = "contract"
	DuplicationTypeCredential ctypes.DuplicationType = "credential"
)

func CheckDuplication(opr *processor.OperationProcessor, op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var duplicationTypeSenderID string
	var duplicationTypeCurrencyID string
	var duplicationTypeCredentialID []string
	var duplicationTypeContractID string
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
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected RegisterCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = processor.DuplicationKey(fact.Currency().Currency().String(), DuplicationTypeCurrency)
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = processor.DuplicationKey(fact.Currency().String(), DuplicationTypeCurrency)
	case currency.Mint:
	case extension.CreateContractAccount:
		fact, ok := t.Fact().(extension.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeContract)
	case extension.Withdraw:
		fact, ok := t.Fact().(extension.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case credential.RegisterModel:
		fact, ok := t.Fact().(credential.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected CreateServiceFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = processor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case credential.AddTemplate:
		fact, ok := t.Fact().(credential.AddTemplateFact)
		if !ok {
			return errors.Errorf("expected AddTemplateFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case credential.Issue:
		fact, ok := t.Fact().(credential.IssueFact)
		if !ok {
			return errors.Errorf("expected IssueFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = processor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		var credentials []string
		for _, v := range fact.Items() {
			key := processor.DuplicationKey(fmt.Sprintf("%s-%s-%s", v.Contract().String(), v.TemplateID(), v.CredentialID()), DuplicationTypeCredential)
			credentials = append(credentials, key)
		}
		duplicationTypeCredentialID = credentials
	default:
		return nil
	}

	if len(duplicationTypeSenderID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			return errors.Errorf("proposal cannot have duplicated sender, %v", duplicationTypeSenderID)
		}

		opr.Duplicated[duplicationTypeSenderID] = struct{}{}
	}

	if len(duplicationTypeCurrencyID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeCurrencyID]; found {
			return errors.Errorf(
				"cannot register duplicated currency id, %v within a proposal",
				duplicationTypeCurrencyID,
			)
		}

		opr.Duplicated[duplicationTypeCurrencyID] = struct{}{}
	}
	if len(duplicationTypeContractID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContractID]; found {
			return errors.Errorf(
				"cannot use a duplicated contract for registering in contract model , %v within a proposal",
				duplicationTypeSenderID,
			)
		}

		opr.Duplicated[duplicationTypeContractID] = struct{}{}
	}
	if len(duplicationTypeCredentialID) > 0 {
		for _, v := range duplicationTypeCredentialID {
			if _, found := opr.Duplicated[v]; found {
				return errors.Errorf(
					"cannot use a duplicated contract-template-credential for credential model , %v within a proposal",
					v,
				)
			}
			opr.Duplicated[v] = struct{}{}
		}
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func GetNewProcessor(opr *processor.OperationProcessor, op base.Operation) (base.OperationProcessor, bool, error) {
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
		extension.CreateContractAccount,
		extension.Withdraw,
		currency.RegisterCurrency,
		currency.UpdateCurrency,
		currency.Mint,
		credential.RegisterModel,
		credential.AddTemplate,
		credential.Issue,
		credential.Revoke:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
