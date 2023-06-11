package types

import (
	"regexp"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

type Uint256 struct {
	n uint256.Int
}

func NewUint256FromString(s string) (Uint256, error) {
	b, err := common.NewBigFromString(s)
	if err != nil {
		return Uint256{}, err
	}

	i := b.Int

	n, ok := uint256.FromBig(i)
	if !ok {
		return Uint256{}, errors.Errorf("not big string")
	}

	return Uint256{
		n: *n,
	}, nil
}

func (id Uint256) Bytes() []byte {
	return id.n.Bytes()
}

func (id Uint256) String() string {
	return id.n.ToBig().String()
}

func (id Uint256) IsValid([]byte) error {
	return nil
}

func (id Uint256) N() *uint256.Int {
	return &id.n
}

var (
	ReValidDate = regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`)
	DateLayout  = "yyyy-MM-dd"
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
		return util.ErrInvalid.Errorf("wrong date, %q", s)
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
