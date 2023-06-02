package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-credential/credential"
	"github.com/ProtoconNet/mitum2/base"
)

type RevokeCredentialsCommand struct {
	baseCommand
	OperationFlags
	Sender            AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract          AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	CredentialService ContractIDFlag `arg:"" name:"credential-service-id" help:"credential id" required:"true"`
	Holder            AddressFlag    `arg:"" name:"holder" help:"credential holder" required:"true"`
	TemplateID        string         `arg:"" name:"template-id" help:"template id" required:"true"`
	ID                string         `arg:"" name:"id" help:"credential id" required:"true"`
	Currency          CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender            base.Address
	contract          base.Address
	holder            base.Address
	tid               credential.Uint256
}

func NewRevokeCredentialsCommand() RevokeCredentialsCommand {
	cmd := NewbaseCommand()
	return RevokeCredentialsCommand{
		baseCommand: *cmd,
	}
}

func (cmd *RevokeCredentialsCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

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

	tid, err := credential.NewUint256FromString(cmd.TemplateID)
	if err != nil {
		return errors.Wrapf(err, "invalid template id format, %q", cmd.TemplateID)
	}
	cmd.tid = tid

	return nil
}

func (cmd *RevokeCredentialsCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []credential.RevokeCredentialsItem

	item := credential.NewRevokeCredentialsItem(
		cmd.contract,
		cmd.CredentialService.ID,
		cmd.holder,
		cmd.tid,
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
