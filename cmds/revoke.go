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
	Sender     currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	Holder     currencycmds.AddressFlag    `arg:"" name:"holder" help:"credential holder" required:"true"`
	TemplateID string                      `arg:"" name:"template-id" help:"template id" required:"true"`
	ID         string                      `arg:"" name:"id" help:"credential id" required:"true"`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	holder     base.Address
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

	PrettyPrint(cmd.Out, op)

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
	var items []credential.RevokeItem

	item := credential.NewRevokeItem(
		cmd.contract,
		cmd.holder,
		cmd.TemplateID,
		cmd.ID,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := credential.NewRevokeFact([]byte(cmd.Token), cmd.sender, items)

	op := credential.NewRevoke(fact)
	err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to revoke operation")
	}

	return op, nil
}
