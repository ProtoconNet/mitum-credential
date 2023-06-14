package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AssignCredentialsFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringErrorFunc("failed to unmarshal AssignCredentialsFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return e(err, "")
	}

	items := make([]AssignCredentialsItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(AssignCredentialsItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected AssignCredentialsItem, not %T", hit[i]), "")
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
