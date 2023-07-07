package cmds

import (
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/digest"
	digestisaac "github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: common.BaseStateHint, Instance: common.BaseState{}},
	{Hint: common.NodeHint, Instance: common.BaseNode{}},
	{Hint: currencytypes.AccountHint, Instance: currencytypes.Account{}},
	{Hint: currencytypes.AddressHint, Instance: currencytypes.Address{}},
	{Hint: currencytypes.AmountHint, Instance: currencytypes.Amount{}},
	{Hint: currencytypes.AccountKeysHint, Instance: currencytypes.BaseAccountKeys{}},
	{Hint: currencytypes.AccountKeyHint, Instance: currencytypes.BaseAccountKey{}},
	{Hint: currencytypes.CurrencyDesignHint, Instance: currencytypes.CurrencyDesign{}},
	{Hint: currencytypes.CurrencyPolicyHint, Instance: currencytypes.CurrencyPolicy{}},
	{Hint: currencytypes.ContractAccountKeysHint, Instance: currencytypes.ContractAccountKeys{}},
	{Hint: currencytypes.NilFeeerHint, Instance: currencytypes.NilFeeer{}},
	{Hint: currencytypes.FixedFeeerHint, Instance: currencytypes.FixedFeeer{}},
	{Hint: currencytypes.RatioFeeerHint, Instance: currencytypes.RatioFeeer{}},

	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.CredentialHint, Instance: types.Credential{}},
	{Hint: types.TemplateHint, Instance: types.Template{}},
	{Hint: types.HolderHint, Instance: types.Holder{}},

	{Hint: currency.CreateAccountsItemMultiAmountsHint, Instance: currency.CreateAccountsItemMultiAmounts{}},
	{Hint: currency.CreateAccountsItemSingleAmountHint, Instance: currency.CreateAccountsItemSingleAmount{}},
	{Hint: currency.CreateAccountsHint, Instance: currency.CreateAccounts{}},
	{Hint: currency.KeyUpdaterHint, Instance: currency.KeyUpdater{}},
	{Hint: currency.TransfersItemMultiAmountsHint, Instance: currency.TransfersItemMultiAmounts{}},
	{Hint: currency.TransfersItemSingleAmountHint, Instance: currency.TransfersItemSingleAmount{}},
	{Hint: currency.TransfersHint, Instance: currency.Transfers{}},
	{Hint: currency.SuffrageInflationHint, Instance: currency.SuffrageInflation{}},
	{Hint: currency.CurrencyPolicyUpdaterHint, Instance: currency.CurrencyPolicyUpdater{}},
	{Hint: currency.CurrencyRegisterHint, Instance: currency.CurrencyRegister{}},
	{Hint: currency.GenesisCurrenciesFactHint, Instance: currency.GenesisCurrenciesFact{}},
	{Hint: currency.GenesisCurrenciesHint, Instance: currency.GenesisCurrencies{}},

	{Hint: extension.CreateContractAccountsItemMultiAmountsHint, Instance: extension.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extension.CreateContractAccountsItemSingleAmountHint, Instance: extension.CreateContractAccountsItemSingleAmount{}},
	{Hint: extension.CreateContractAccountsHint, Instance: extension.CreateContractAccounts{}},
	{Hint: extension.WithdrawsItemMultiAmountsHint, Instance: extension.WithdrawsItemMultiAmounts{}},
	{Hint: extension.WithdrawsItemSingleAmountHint, Instance: extension.WithdrawsItemSingleAmount{}},
	{Hint: extension.WithdrawsHint, Instance: extension.Withdraws{}},

	{Hint: credential.CreateCredentialServiceHint, Instance: credential.CreateCredentialService{}},
	{Hint: credential.AddTemplateHint, Instance: credential.AddTemplate{}},
	{Hint: credential.AssignCredentialsItemHint, Instance: credential.AssignCredentialsItem{}},
	{Hint: credential.AssignCredentialsHint, Instance: credential.AssignCredentials{}},
	{Hint: credential.RevokeCredentialsItemHint, Instance: credential.RevokeCredentialsItem{}},
	{Hint: credential.RevokeCredentialsHint, Instance: credential.RevokeCredentials{}},

	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.CredentialStateValueHint, Instance: state.CredentialStateValue{}},
	{Hint: state.HolderDIDStateValueHint, Instance: state.HolderDIDStateValue{}},
	{Hint: state.TemplateStateValueHint, Instance: state.TemplateStateValue{}},

	{Hint: statecurrency.AccountStateValueHint, Instance: statecurrency.AccountStateValue{}},
	{Hint: statecurrency.BalanceStateValueHint, Instance: statecurrency.BalanceStateValue{}},
	{Hint: statecurrency.CurrencyDesignStateValueHint, Instance: statecurrency.CurrencyDesignStateValue{}},
	{Hint: stateextension.ContractAccountStateValueHint, Instance: stateextension.ContractAccountStateValue{}},

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: currency.CreateAccountsFactHint, Instance: currency.CreateAccountsFact{}},
	{Hint: currency.KeyUpdaterFactHint, Instance: currency.KeyUpdaterFact{}},
	{Hint: currency.TransfersFactHint, Instance: currency.TransfersFact{}},
	{Hint: currency.SuffrageInflationFactHint, Instance: currency.SuffrageInflationFact{}},

	{Hint: currency.CurrencyRegisterFactHint, Instance: currency.CurrencyRegisterFact{}},
	{Hint: currency.CurrencyPolicyUpdaterFactHint, Instance: currency.CurrencyPolicyUpdaterFact{}},
	{Hint: extension.CreateContractAccountsFactHint, Instance: extension.CreateContractAccountsFact{}},
	{Hint: extension.WithdrawsFactHint, Instance: extension.WithdrawsFact{}},

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
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
