package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	ServiceID ServiceID `json:"service_id"`
	Policy    Policy    `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		ServiceID:  de.serviceID,
		Policy:     de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint      hint.Hint       `json:"_hint"`
	ServiceID string          `json:"service_id"`
	Policy    json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Design")

	var ud DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e.Wrap(err)
	}

	return de.unpack(enc, ud.Hint, ud.ServiceID, ud.Policy)
}
