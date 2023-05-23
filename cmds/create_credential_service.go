package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-credential/credential"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type CreateCredentialServiceCommand struct {
	baseCommand
	OperationFlags
	Sender      AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract    AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	Credential  ContractIDFlag `arg:"" name:"credential-id" help:"credential id" required:"true"`
	Currency    CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender      base.Address
	contract    base.Address
	controllers []base.Address
}

func NewCreateCredentialServiceCommand() CreateCredentialServiceCommand {
	cmd := NewbaseCommand()
	return CreateCredentialServiceCommand{
		baseCommand: *cmd,
	}
}

func (cmd *CreateCredentialServiceCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *CreateCredentialServiceCommand) parseFlags() error {
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

	return nil
}

func (cmd *CreateCredentialServiceCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringErrorFunc("failed to create create-credential-service operation")

	fact := credential.NewCreateCredentialServiceFact([]byte(cmd.Token), cmd.sender, cmd.contract, cmd.Credential.ID, cmd.Currency.CID)
	if err := fact.IsValid(nil); err != nil {
		return nil, e(err, "")
	}

	op, err := credential.NewCreateCredentialService(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
