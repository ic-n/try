package try_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/ic-n/try"
)

var errSyntetic = errors.New("syntetic error")

func Test(t *testing.T) {
	t.Run("try_without_error", func(t *testing.T) {
		tc, cancel := try.New(context.Background())
		defer cancel()

		tc.Try(func() error {
			return nil
		})
		tc.Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("try_with_error", func(t *testing.T) {
		tc, cancel := try.New(context.Background())
		defer cancel()

		tc.Try(func() error {
			return errSyntetic
		})
		tc.Try(func() error {
			panic("must not be called")
		})
		tc.Catch(func(err error) {
			if !errors.Is(err, errSyntetic) {
				t.Error(err)
			}
		})
	})

	t.Run("try_panic", func(t *testing.T) {
		tc, cancel := try.New(context.Background())
		defer cancel()

		tc.Try(func() error {
			panic("hello")
		})
		var err error
		tc.Catch(func(e error) {
			err = e
		})

		if err == nil || err.Error() != "panic: hello" {
			t.Errorf("error expected to be \"panic: hello\", got %v", err)
		}
	})

	t.Run("try_exact_error", func(t *testing.T) {
		tc, cancel := try.New(context.Background())
		defer cancel()

		tc.Try(func() error {
			return errSyntetic
		})
		tc.CatchError(os.ErrInvalid, func(err error) {
			t.Error(err)
		})
		tc.CatchError(errSyntetic, func(err error) {})
		tc.Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("try_pass_error", func(t *testing.T) {
		tc, cancel := try.New(context.Background())
		defer cancel()

		tc.Try(func() error {
			return errSyntetic
		})
		tc.CatchError(os.ErrInvalid, func(err error) {
			t.Error(err)
		})
		tc.PassError(errSyntetic)
		tc.Catch(func(err error) {
			t.Error(err)
		})
	})
}
