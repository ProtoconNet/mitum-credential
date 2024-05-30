package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-credential/operation/credential"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type CreateServiceCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
}

func (cmd *CreateServiceCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

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

func (cmd *CreateServiceCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	return nil
}

func (cmd *CreateServiceCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError("failed to create create-service operation")

	fact := credential.NewCreateServiceFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.Currency.CID,
	)

	op := credential.NewCreateService(fact)
	err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
