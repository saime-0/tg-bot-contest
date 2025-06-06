package ue

import (
	"database/sql"
	"errors"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
)

type Err struct {
	text string
}

func (u *Err) Error() string {
	return u.text
}

func New(s string) *Err {
	return &Err{text: s}
}

func Sql(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return New(l10n.RequestedDataNotFound)
	}

	return err
}
