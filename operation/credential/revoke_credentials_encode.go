package credential

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RevokeCredentialsFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringErrorFunc("failed to unmarshal RevokeCredentialsFact")

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

	items := make([]RevokeCredentialsItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(RevokeCredentialsItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected RevokeCredentialsItem, not %T", hit[i]), "")
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
