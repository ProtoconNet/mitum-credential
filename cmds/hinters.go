package cmds

import (
	"github.com/ProtoconNet/mitum-credential/credential"
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum-currency/v3/digest"
	digestisaac "github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensionstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: currencybase.BaseStateHint, Instance: currencybase.BaseState{}},
	{Hint: currencybase.NodeHint, Instance: currencybase.BaseNode{}},
	{Hint: currencybase.AccountHint, Instance: currencybase.Account{}},
	{Hint: currencybase.AddressHint, Instance: currencybase.Address{}},
	{Hint: currencybase.AmountHint, Instance: currencybase.Amount{}},
	{Hint: currencybase.AccountKeysHint, Instance: currencybase.BaseAccountKeys{}},
	{Hint: currencybase.AccountKeyHint, Instance: currencybase.BaseAccountKey{}},
	{Hint: mitumcurrency.CreateAccountsItemMultiAmountsHint, Instance: mitumcurrency.CreateAccountsItemMultiAmounts{}},
	{Hint: mitumcurrency.CreateAccountsItemSingleAmountHint, Instance: mitumcurrency.CreateAccountsItemSingleAmount{}},
	{Hint: mitumcurrency.CreateAccountsHint, Instance: mitumcurrency.CreateAccounts{}},
	{Hint: mitumcurrency.KeyUpdaterHint, Instance: mitumcurrency.KeyUpdater{}},
	{Hint: mitumcurrency.TransfersItemMultiAmountsHint, Instance: mitumcurrency.TransfersItemMultiAmounts{}},
	{Hint: mitumcurrency.TransfersItemSingleAmountHint, Instance: mitumcurrency.TransfersItemSingleAmount{}},
	{Hint: mitumcurrency.TransfersHint, Instance: mitumcurrency.Transfers{}},
	{Hint: mitumcurrency.SuffrageInflationHint, Instance: mitumcurrency.SuffrageInflation{}},
	{Hint: currencystate.AccountStateValueHint, Instance: currencystate.AccountStateValue{}},
	{Hint: currencystate.BalanceStateValueHint, Instance: currencystate.BalanceStateValue{}},

	{Hint: currencybase.CurrencyDesignHint, Instance: currencybase.CurrencyDesign{}},
	{Hint: currencybase.CurrencyPolicyHint, Instance: currencybase.CurrencyPolicy{}},
	{Hint: mitumcurrency.CurrencyRegisterHint, Instance: mitumcurrency.CurrencyRegister{}},
	{Hint: mitumcurrency.CurrencyPolicyUpdaterHint, Instance: mitumcurrency.CurrencyPolicyUpdater{}},
	{Hint: currencybase.ContractAccountKeysHint, Instance: currencybase.ContractAccountKeys{}},
	{Hint: extensioncurrency.CreateContractAccountsItemMultiAmountsHint, Instance: extensioncurrency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extensioncurrency.CreateContractAccountsItemSingleAmountHint, Instance: extensioncurrency.CreateContractAccountsItemSingleAmount{}},
	{Hint: extensioncurrency.CreateContractAccountsHint, Instance: extensioncurrency.CreateContractAccounts{}},
	{Hint: extensioncurrency.WithdrawsItemMultiAmountsHint, Instance: extensioncurrency.WithdrawsItemMultiAmounts{}},
	{Hint: extensioncurrency.WithdrawsItemSingleAmountHint, Instance: extensioncurrency.WithdrawsItemSingleAmount{}},
	{Hint: extensioncurrency.WithdrawsHint, Instance: extensioncurrency.Withdraws{}},
	{Hint: mitumcurrency.GenesisCurrenciesFactHint, Instance: mitumcurrency.GenesisCurrenciesFact{}},
	{Hint: mitumcurrency.GenesisCurrenciesHint, Instance: mitumcurrency.GenesisCurrencies{}},
	{Hint: currencybase.NilFeeerHint, Instance: currencybase.NilFeeer{}},
	{Hint: currencybase.FixedFeeerHint, Instance: currencybase.FixedFeeer{}},
	{Hint: currencybase.RatioFeeerHint, Instance: currencybase.RatioFeeer{}},
	{Hint: extensionstate.ContractAccountStateValueHint, Instance: extensionstate.ContractAccountStateValue{}},
	{Hint: currencystate.CurrencyDesignStateValueHint, Instance: currencystate.CurrencyDesignStateValue{}},

	{Hint: credential.DesignHint, Instance: credential.Design{}},
	{Hint: credential.DesignStateValueHint, Instance: credential.DesignStateValue{}},
	{Hint: credential.PolicyHint, Instance: credential.Policy{}},
	{Hint: credential.CredentialHint, Instance: credential.Credential{}},
	{Hint: credential.CredentialStateValueHint, Instance: credential.CredentialStateValue{}},
	{Hint: credential.HolderDIDStateValueHint, Instance: credential.HolderDIDStateValue{}},
	{Hint: credential.TemplateHint, Instance: credential.Template{}},
	{Hint: credential.TemplateStateValueHint, Instance: credential.TemplateStateValue{}},
	{Hint: credential.CreateCredentialServiceHint, Instance: credential.CreateCredentialService{}},
	{Hint: credential.AddTemplateHint, Instance: credential.AddTemplate{}},
	{Hint: credential.AssignCredentialsItemHint, Instance: credential.AssignCredentialsItem{}},
	{Hint: credential.AssignCredentialsHint, Instance: credential.AssignCredentials{}},
	{Hint: credential.RevokeCredentialsItemHint, Instance: credential.RevokeCredentialsItem{}},
	{Hint: credential.RevokeCredentialsHint, Instance: credential.RevokeCredentials{}},

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},

	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: mitumcurrency.CreateAccountsFactHint, Instance: mitumcurrency.CreateAccountsFact{}},
	{Hint: mitumcurrency.KeyUpdaterFactHint, Instance: mitumcurrency.KeyUpdaterFact{}},
	{Hint: mitumcurrency.TransfersFactHint, Instance: mitumcurrency.TransfersFact{}},
	{Hint: mitumcurrency.SuffrageInflationFactHint, Instance: mitumcurrency.SuffrageInflationFact{}},

	{Hint: mitumcurrency.CurrencyRegisterFactHint, Instance: mitumcurrency.CurrencyRegisterFact{}},
	{Hint: mitumcurrency.CurrencyPolicyUpdaterFactHint, Instance: mitumcurrency.CurrencyPolicyUpdaterFact{}},
	{Hint: extensioncurrency.CreateContractAccountsFactHint, Instance: extensioncurrency.CreateContractAccountsFact{}},
	{Hint: extensioncurrency.WithdrawsFactHint, Instance: extensioncurrency.WithdrawsFact{}},

	{Hint: credential.CreateCredentialServiceFactHint, Instance: credential.CreateCredentialServiceFact{}},
	{Hint: credential.AddTemplateFactHint, Instance: credential.AddTemplateFact{}},
	{Hint: credential.AssignCredentialsFactHint, Instance: credential.AssignCredentialsFact{}},
	{Hint: credential.RevokeCredentialsFactHint, Instance: credential.RevokeCredentialsFact{}},

	{Hint: isaacoperation.GenesisNetworkPolicyFactHint, Instance: isaacoperation.GenesisNetworkPolicyFact{}},
	{Hint: isaacoperation.SuffrageCandidateFactHint, Instance: isaacoperation.SuffrageCandidateFact{}},
	{Hint: isaacoperation.SuffrageDisjoinFactHint, Instance: isaacoperation.SuffrageDisjoinFact{}},
	{Hint: isaacoperation.SuffrageJoinFactHint, Instance: isaacoperation.SuffrageJoinFact{}},
	{Hint: isaacoperation.SuffrageGenesisJoinFactHint, Instance: isaacoperation.SuffrageGenesisJoinFact{}},
}

func init() {
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for i := range Hinters {
		if err := enc.Add(Hinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := enc.Add(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
