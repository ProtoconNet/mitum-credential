package cmds

import (
	"context"
	"github.com/ProtoconNet/mitum2/util"

	"github.com/ProtoconNet/mitum-credential/operation/credential"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type IssueCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender     currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	Holder     currencycmds.AddressFlag    `arg:"" name:"holder" help:"credential holder" required:"true"`
	TemplateID string                      `arg:"" name:"template-id" help:"template id" required:"true"`
	ID         string                      `arg:"" name:"id" help:"credential id" required:"true"`
	Value      string                      `arg:"" name:"value" help:"credential value" required:"true"`
	ValidFrom  uint64                      `arg:"" name:"valid-from" help:"valid from" required:"true"`
	ValidUntil uint64                      `arg:"" name:"valid-until" help:"valid until" required:"true"`
	DID        string                      `arg:"" name:"did" help:"did" required:"true"`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	holder     base.Address
}

func (cmd *IssueCommand) Run(pctx context.Context) error {
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

func (cmd *IssueCommand) parseFlags() error {
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

	holder, err := cmd.Holder.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid holder account format, %q", cmd.Holder.String())
	}
	cmd.holder = holder

	return nil
}

func (cmd *IssueCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create issue operation")

	var items []credential.IssueItem
	item := credential.NewIssueItem(
		cmd.contract,
		cmd.holder,
		cmd.TemplateID,
		cmd.ID,
		cmd.Value,
		cmd.ValidFrom,
		cmd.ValidUntil,
		cmd.DID,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := credential.NewIssueFact([]byte(cmd.Token), cmd.sender, items)

	op := credential.NewIssue(fact)

	err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
