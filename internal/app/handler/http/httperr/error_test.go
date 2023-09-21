package httperr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"web-studio-backend/internal/app/domain/apperror"
)

func TestUnwrapAppError(t *testing.T) {
	tempAppError := apperror.NewInvalidRequest("msg", "field")

	tests := []struct {
		err      error
		appError bool
	}{
		{
			err:      fmt.Errorf("err3: %w", fmt.Errorf("err2: %w", tempAppError)),
			appError: true,
		},
		{
			err:      fmt.Errorf("err2: %w", fmt.Errorf("err1: message")),
			appError: false,
		},
		{
			err:      tempAppError,
			appError: true,
		},
		{
			err:      fmt.Errorf("err5: %w", fmt.Errorf("err4: %w", fmt.Errorf("err3: %w", fmt.Errorf("err2: %w", tempAppError)))),
			appError: true,
		},
		{
			err:      fmt.Errorf("err"),
			appError: false,
		},
		{
			err:      nil,
			appError: false,
		},
	}

	for _, tt := range tests {
		ae := UnwrapAppError(tt.err)

		if tt.appError {
			require.NotNil(t, ae)
		} else {
			require.Nil(t, ae)
		}
	}
}
