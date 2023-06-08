package credential

import (
	"encoding/json"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Credential currencybase.ContractID `json:"credential_service_id"`
	Policy     Policy                  `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		Credential: de.credentialServiceID,
		Policy:     de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	Credential string          `json:"credential_service_id"`
	Policy     json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Design")

	var ud DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ud.Hint, ud.Credential, ud.Policy)
}
