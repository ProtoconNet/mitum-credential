package cmds

import (
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.CredentialHint, Instance: types.Credential{}},
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.HolderHint, Instance: types.Holder{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.TemplateHint, Instance: types.Template{}},

	{Hint: credential.CreateCredentialServiceHint, Instance: credential.CreateCredentialService{}},
	{Hint: credential.AddTemplateHint, Instance: credential.AddTemplate{}},
	{Hint: credential.AssignCredentialsItemHint, Instance: credential.AssignCredentialsItem{}},
	{Hint: credential.AssignCredentialsHint, Instance: credential.AssignCredentials{}},
	{Hint: credential.RevokeCredentialsItemHint, Instance: credential.RevokeCredentialsItem{}},
	{Hint: credential.RevokeCredentialsHint, Instance: credential.RevokeCredentials{}},

	{Hint: state.CredentialStateValueHint, Instance: state.CredentialStateValue{}},
	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.HolderDIDStateValueHint, Instance: state.HolderDIDStateValue{}},
	{Hint: state.TemplateStateValueHint, Instance: state.TemplateStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: credential.AddTemplateFactHint, Instance: credential.AddTemplateFact{}},
	{Hint: credential.AssignCredentialsFactHint, Instance: credential.AssignCredentialsFact{}},
	{Hint: credential.CreateCredentialServiceFactHint, Instance: credential.CreateCredentialServiceFact{}},
	{Hint: credential.RevokeCredentialsFactHint, Instance: credential.RevokeCredentialsFact{}},
}

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)
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
