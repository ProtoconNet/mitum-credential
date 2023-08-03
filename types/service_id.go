package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"regexp"
	"unicode/utf8"
)

var (
	MinLengthContractID = 3
	MaxLengthContractID = 20
	REServiceIDString   = `^[A-Za-z0-9가-힣][A-Za-z0-9가-힣-_]*$`
	REServiceIDExp      = regexp.MustCompile(REServiceIDString)
)

type ServiceID string

func (sid ServiceID) Bytes() []byte {
	return []byte(sid)
}

func (sid ServiceID) String() string {
	return string(sid)
}

func (sid ServiceID) IsValid([]byte) error {
	if l := utf8.RuneCountInString(sid.String()); l < MinLengthContractID || l > MaxLengthContractID {
		return util.ErrInvalid.Errorf(
			"invalid length of service id, %d <= length <= %d",
			MinLengthContractID,
			MaxLengthContractID,
		)
	}
	if !REServiceIDExp.Match([]byte(sid)) {
		return util.ErrInvalid.Errorf("wrong service id, %v", sid)
	}

	return nil
}
