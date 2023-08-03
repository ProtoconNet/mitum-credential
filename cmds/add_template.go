package cmds

import (
	"context"
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-credential/types"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type AddTemplateCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender            currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract          currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	CredentialService ServiceIDFlag               `arg:"" name:"credential-service-id" help:"credential id" required:"true"`
	TemplateID        string                      `arg:"" name:"template-id" help:"template id" required:"true"`
	TemplateName      string                      `arg:"" name:"template-name" help:"template name"  required:"true"`
	ServiceDate       string                      `arg:"" name:"service-date" help:"service date; yyyy-MM-dd" required:"true"`
	ExpirationDate    string                      `arg:"" name:"expiration-date" help:"expiration date; yyyy-MM-dd" required:"true"`
	TemplateShare     bool                        `name:"template-share" help:"template share; true | false" required:"true"`
	MultiAudit        bool                        `name:"multi-audit" help:"multi audit; true | false" required:"true"`
	DisplayName       string                      `arg:"" name:"display-name" help:"display name" required:"true"`
	SubjectKey        string                      `arg:"" name:"subject-key" help:"subject key" required:"true"`
	Description       string                      `arg:"" name:"description" help:"description"  required:"true"`
	Creator           currencycmds.AddressFlag    `arg:"" name:"creator" help:"creator address"  required:"true"`
	Currency          currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender            base.Address
	contract          base.Address
	service           types.Date
	expiration        types.Date
	creator           base.Address
}

func (cmd *AddTemplateCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.Encoders
	enc = cmd.Encoder

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *AddTemplateCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	creator, err := cmd.Creator.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid creator account format, %q", cmd.Creator.String())
	}
	cmd.creator = creator

	service, expiration := types.Date(cmd.ServiceDate), types.Date(cmd.ExpirationDate)
	if err := service.IsValid(nil); err != nil {
		return errors.Wrapf(err, "invalid service date format, %q", cmd.ServiceDate)
	}
	if err := expiration.IsValid(nil); err != nil {
		return errors.Wrapf(err, "invalid expiration date format, %q", cmd.ExpirationDate)
	}
	cmd.service = service
	cmd.expiration = expiration

	return nil
}

func (cmd *AddTemplateCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError("failed to create add-template operation")

	fact := credential.NewAddTemplateFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.CredentialService.ID,
		cmd.TemplateID,
		cmd.TemplateName,
		cmd.service,
		cmd.expiration,
		types.Bool(cmd.TemplateShare),
		types.Bool(cmd.MultiAudit),
		cmd.DisplayName,
		cmd.SubjectKey,
		cmd.Description,
		cmd.creator,
		"MCC",
		//cmd.Currency.CID,
	)

	op, err := credential.NewAddTemplate(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}

	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
