package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/pkg/errors"
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
		return common.ErrValOOR.Wrap(errors.Errorf(
			"%d <= length of service id <= %d",
			MinLengthContractID,
			MaxLengthContractID,
		))
	}
	if !REServiceIDExp.Match([]byte(sid)) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("service ID %s, must match regex `^[A-Za-z0-9가-힣][A-Za-z0-9가-힣-_]*$`", sid))
	}

	return nil
}
