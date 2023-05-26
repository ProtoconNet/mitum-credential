package credential

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func (u Uint256) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, u.String()), nil
}

func (u *Uint256) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	if t != bsontype.String {
		return errors.Errorf("invalid marshaled type for Uint256, %v", t)
	}

	s, _, ok := bsoncore.ReadString(b)
	if !ok {
		return errors.Errorf("can not read string")
	}

	ui, err := NewUint256FromString(s)
	if err != nil {
		return err
	}
	*u = ui

	return nil
}
