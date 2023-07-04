package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-credential/operation/credential"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type CreateCredentialServiceCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender            currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract          currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	CredentialService currencycmds.ContractIDFlag `arg:"" name:"credential-service-id" help:"credential id" required:"true"`
	Currency          currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender            base.Address
	contract          base.Address
}

func (cmd *CreateCredentialServiceCommand) Run(pctx context.Context) error { // nolint:dupl
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
	e := util.StringError("failed to create create-credential-service operation")

	fact := credential.NewCreateCredentialServiceFact([]byte(cmd.Token), cmd.sender, cmd.contract, cmd.CredentialService.ID, cmd.Currency.CID)

	op, err := credential.NewCreateCredentialService(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
