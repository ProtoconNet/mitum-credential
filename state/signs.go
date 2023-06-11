package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

func CheckFactSignsByState(
	address base.Address,
	fs []base.Sign,
	getState base.GetStateFunc,
) error {
	st, err := state.ExistsState(currency.StateKeyAccount(address), "keys of account", getState)
	if err != nil {
		return err
	}
	keys, err := currency.StateKeysValue(st)
	switch {
	case err != nil:
		return base.NewBaseOperationProcessReasonError("failed to get Keys %w", err)
	case keys == nil:
		return base.NewBaseOperationProcessReasonError("empty keys found")
	}

	if err := checkThreshold(fs, keys); err != nil {
		return base.NewBaseOperationProcessReasonError("failed to check threshold %w", err)
	}

	return nil
}

func checkThreshold(fs []base.Sign, keys types.AccountKeys) error {
	var sum uint
	for i := range fs {
		ky, found := keys.Key(fs[i].Signer())
		if !found {
			return errors.Errorf("unknown key found, %q", fs[i].Signer())
		}
		sum += ky.Weight()
	}

	if sum < keys.Threshold() {
		return errors.Errorf("not passed threshold, sum=%d < threshold=%d", sum, keys.Threshold())
	}

	return nil
}
