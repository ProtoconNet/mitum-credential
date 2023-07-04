package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-credential/operation/credential"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type RevokeCredentialsCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender            currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract          currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	CredentialService currencycmds.ContractIDFlag `arg:"" name:"credential-service-id" help:"credential id" required:"true"`
	Holder            currencycmds.AddressFlag    `arg:"" name:"holder" help:"credential holder" required:"true"`
	TemplateID        uint64                      `arg:"" name:"template-id" help:"template id" required:"true"`
	ID                string                      `arg:"" name:"id" help:"credential id" required:"true"`
	Currency          currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender            base.Address
	contract          base.Address
	holder            base.Address
}

func (cmd *RevokeCredentialsCommand) Run(pctx context.Context) error {
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

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *RevokeCredentialsCommand) parseFlags() error {
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

	holder, err := cmd.Holder.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid holder account format, %q", cmd.Holder.String())
	}
	cmd.holder = holder

	return nil
}

func (cmd *RevokeCredentialsCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []credential.RevokeCredentialsItem

	item := credential.NewRevokeCredentialsItem(
		cmd.contract,
		cmd.CredentialService.ID,
		cmd.holder,
		cmd.TemplateID,
		cmd.ID,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := credential.NewRevokeCredentialsFact([]byte(cmd.Token), cmd.sender, items)

	op, err := credential.NewRevokeCredentials(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to revoke credentials operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to revoke credentials operation")
	}

	return op, nil
}
