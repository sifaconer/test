package common

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"
)

func CheckDBErrorType(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return nil
	}

	pgErr, ok := err.(pgdriver.Error)
	if !ok {
		return err
	}
	switch pgErr.Field('C') {
	case pgerrcode.UniqueViolation:
		return errors.New("a resource with this data already exists, please provide valid information")
	case pgerrcode.NotNullViolation:
		return errors.New("the value cannot be null, please provide a valid value")
	case pgerrcode.ForeignKeyViolation:
		return errors.New("the value does not exist in the referenced table, please provide a valid reference")
	case pgerrcode.CheckViolation:
		return errors.New("the value does not meet the required conditions, please provide a valid value")
	case pgerrcode.ExclusionViolation:
		return errors.New("the value is excluded from the valid range, please provide a valid value")
	case pgerrcode.IntegrityConstraintViolation:
		return errors.New("the value does not meet the integrity constraints of the referenced table, please provide a valid value")
	case pgerrcode.RestrictViolation:
		return errors.New("the value does not meet the required conditions, please provide a valid value")
	default:
		return err
	}
}
