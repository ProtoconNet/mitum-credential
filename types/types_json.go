package types

import "github.com/ProtoconNet/mitum2/util"

func (u Uint256) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(u.String())
}

func (u *Uint256) UnmarshalJSON(b []byte) error {
	var s string
	if err := util.UnmarshalJSON(b, &s); err != nil {
		return err
	}

	ui, err := NewUint256FromString(s)
	if err != nil {
		return err
	}
	*u = ui

	return nil
}
