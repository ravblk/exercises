package brackets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testValidate = map[string]string{
	"[()]{}{[()()]()}":        OutputsBalanced,
	"[(])":                    OutputsNotBalanced,
	"[({":                     OutputsNotBalanced,
	")}]":                     OutputsNotBalanced,
	"[()]{}{[()()]()}}}}}}}}": OutputsNotBalanced,
}

var testFix = map[string]string{
	"[()]{}{[()()]()}":        "[()]{}{[()()]()}",
	"[(])":                    "[()]",
	"[({":                     "[({})]",
	")}]":                     "",
	"[()]{}{[()()]()}}}}}}}}": "[()]{}{[()()]()}",
}

func newService() *service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}

	defer logger.Sync()

	return NewService(logger)
}

func TestValidate(t *testing.T) {
	svc := newService()
	for val, expected := range testValidate {
		testBracket := &Brackets{val}

		res, err := svc.Validate(context.Background(), testBracket)
		assert.Nil(t, err)

		assert.Equal(t, expected, res.ResultValidate)
	}
}

func TestFix(t *testing.T) {
	svc := newService()
	for val, expected := range testFix {
		testBracket := &Brackets{val}

		res, err := svc.Fix(context.Background(), testBracket)
		assert.Nil(t, err)

		assert.Equal(t, expected, res.ResultFix)
	}
}
