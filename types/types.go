package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"regexp"
	"time"
)

var (
	ReValidDate = regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`)
	DateLayout  = "2006-01-02"
)

type Date string

func (s Date) Bytes() []byte {
	return []byte(s)
}

func (s Date) String() string {
	return string(s)
}

func (s Date) IsValid([]byte) error {
	if !ReValidDate.Match([]byte(s)) {
		return common.ErrValueInvalid.Errorf("wrong date, %v", s)
	}

	return nil
}

func (s Date) Parse() (time.Time, error) {
	return time.Parse(DateLayout, string(s))
}

type Bool bool

func (b Bool) Bytes() []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}
