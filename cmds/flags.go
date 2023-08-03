package cmds

import "github.com/ProtoconNet/mitum-credential/types"

type ServiceIDFlag struct {
	ID types.ServiceID
}

func (v *ServiceIDFlag) UnmarshalText(b []byte) error {
	id := types.ServiceID(string(b))
	if err := id.IsValid(nil); err != nil {
		return err
	}
	v.ID = id

	return nil
}

func (v *ServiceIDFlag) String() string {
	return v.ID.String()
}
