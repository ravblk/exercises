package service

import (
	"context"
	"ravblk/exercises/service/services/brackets"
)

type Brackets interface {
	Validate(ctx context.Context, in *brackets.Brackets) (*brackets.ResultValidateBrackets, error)
	Fix(ctx context.Context, in *brackets.Brackets) (*brackets.ResultFixBrackets, error)
}

type BracketsMiddleware func(Brackets) Brackets
