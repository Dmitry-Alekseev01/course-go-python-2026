package handler

import (
	"context"
	"errors"

	myerrors "github.com/tcarzverey/course-go-python/homeworks/hw1/handler/errors"
)

type Handler struct {
	db UsersDB
}

func NewHandler(db UsersDB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) UpdateUserBalance(ctx context.Context, userID, balance int64) error {
	err := h.db.UpdateBalance(ctx, userID, balance)
	var re *myerrors.RetryableError
	var nfe *myerrors.NotFoundError
	var ame *myerrors.AdditionalMessageError
	switch {
	case errors.As(err, &re):
		for i := 0; i < re.RetryCount(); i++ {
			err := h.db.UpdateBalance(ctx, userID, balance)
			if err == nil {
				return nil
			}
		}
		return err
	case errors.As(err, &nfe):
		return myerrors.NewAdditionalMessageError(err, "not found")
	case errors.As(err, &ame):
		return err
	case err == nil:
		return nil
	default:
		panic("wrong error type")
	}
}
