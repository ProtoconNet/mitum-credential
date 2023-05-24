package cmds

import (
	extensioncurrencycmds "github.com/ProtoconNet/mitum-currency-extension/v2/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v2/cmds"
)

type OperationCommand struct {
	CreateAccount           currencycmds.CreateAccountCommand                  `cmd:"" name:"create-account" help:"create new account"`
	KeyUpdater              currencycmds.KeyUpdaterCommand                     `cmd:"" name:"key-updater" help:"update account keys"`
	Transfer                currencycmds.TransferCommand                       `cmd:"" name:"transfer" help:"transfer amounts to receiver"`
	CreateContractAccount   extensioncurrencycmds.CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw                extensioncurrencycmds.WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw amounts from target contract account"`
	CreateCredentialService CreateCredentialServiceCommand                     `cmd:"" name:"create-credential-service" help:"register credential service to contract account"`
	CurrencyRegister        currencycmds.CurrencyRegisterCommand               `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater   currencycmds.CurrencyPolicyUpdaterCommand          `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	SuffrageInflation       currencycmds.SuffrageInflationCommand              `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"`
	SuffrageCandidate       currencycmds.SuffrageCandidateCommand              `cmd:"" name:"suffrage-candidate" help:"suffrage candidate operation"`
	SuffrageJoin            currencycmds.SuffrageJoinCommand                   `cmd:"" name:"suffrage-join" help:"suffrage join operation"`
	SuffrageDisjoin         currencycmds.SuffrageDisjoinCommand                `cmd:"" name:"suffrage-disjoin" help:"suffrage disjoin operation"` // revive:disable-line:line-length-limit
}

func NewOperationCommand() OperationCommand {
	return OperationCommand{
		CreateAccount:           currencycmds.NewCreateAccountCommand(),
		KeyUpdater:              currencycmds.NewKeyUpdaterCommand(),
		Transfer:                currencycmds.NewTransferCommand(),
		CreateContractAccount:   extensioncurrencycmds.NewCreateContractAccountCommand(),
		Withdraw:                extensioncurrencycmds.NewWithdrawCommand(),
		CreateCredentialService: NewCreateCredentialServiceCommand(),
		CurrencyRegister:        currencycmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater:   currencycmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:       currencycmds.NewSuffrageInflationCommand(),
		SuffrageCandidate:       currencycmds.NewSuffrageCandidateCommand(),
		SuffrageJoin:            currencycmds.NewSuffrageJoinCommand(),
		SuffrageDisjoin:         currencycmds.NewSuffrageDisjoinCommand(),
	}
}
