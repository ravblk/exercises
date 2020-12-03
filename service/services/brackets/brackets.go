package brackets

import (
	"bytes"
	"context"

	"go.uber.org/zap"
)

var bracketsOpening = map[byte]byte{
	'}': '{',
	')': '(',
	']': '[',
}

var bracketsClosing = map[byte]byte{
	'{': '}',
	'(': ')',
	'[': ']',
}

const (
	OutputsBalanced    = "Balanced"
	OutputsNotBalanced = "Not Balanced"
)

type service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *service {
	return &service{
		logger: logger,
	}
}

func (s *service) Validate(ctx context.Context, in *Brackets) (*ResultValidateBrackets, error) {
	res := &ResultValidateBrackets{
		ResultValidate: OutputsNotBalanced,
	}

	bs := []byte(in.Brackets)

	stack := make([]byte, 0)

	for i := 0; i < len(bs); i++ {
		if v, ok := bracketsOpening[bs[i]]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != v {
				return res, nil
			}

			stack = stack[0 : len(stack)-1]
			continue
		}

		stack = append(stack, bs[i])
	}

	if len(stack) == 0 {
		res.ResultValidate = OutputsBalanced
	}

	return res, nil
}

func (s *service) Fix(ctx context.Context, in *Brackets) (*ResultFixBrackets, error) {
	res := &ResultFixBrackets{}

	valid, err := s.Validate(ctx, in)
	if err != nil {
		s.logger.Warn("", zap.Error(err))
		return nil, err
	}

	if valid.ResultValidate == OutputsBalanced {
		res.ResultFix = in.Brackets

		return res, nil
	}

	bs := []byte(in.Brackets)

	buf := bytes.Buffer{}

	stack := make([]byte, 0)

	for i := 0; i < len(bs); i++ {
		if v, ok := bracketsOpening[bs[i]]; ok {
			if len(stack) == 0 {
				continue
			}

			if stack[len(stack)-1] != v {
				if err := buf.WriteByte(bracketsClosing[stack[len(stack)-1]]); err != nil {
					s.logger.Error("", zap.Error(err))
					return nil, err
				}
			} else {
				if err := buf.WriteByte(bs[i]); err != nil {
					s.logger.Error("", zap.Error(err))
					return nil, err
				}
			}

			stack = stack[0 : len(stack)-1]
			continue
		}

		if err := buf.WriteByte(bs[i]); err != nil {
			s.logger.Error("", zap.Error(err))
			return nil, err
		}

		stack = append(stack, bs[i])
	}

	for len(stack) > 0 {
		if err := buf.WriteByte(bracketsClosing[stack[len(stack)-1]]); err != nil {
			s.logger.Error("", zap.Error(err))
			return nil, err
		}

		stack = stack[0 : len(stack)-1]
	}

	res.ResultFix = buf.String()

	return res, nil
}
