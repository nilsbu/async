package async

import (
	"fmt"
	"time"
)

// T executes a function and measures the execution time.
func T(name string, f func()) {
	start := time.Now()
	f()
	fmt.Printf("%v took %s\n", name, time.Since(start))
}

// P executes fs in parallel
func P(fs []func()) {
	back := make(chan bool, len(fs))

	for _, f := range fs {
		go func(f func()) {
			f()
			back <- true
		}(f)
	}

	for i := 0; i < len(fs); i++ {
		<-back
	}
}

// Pi executes f in parallel n times
func Pi(n int, f func(int)) {
	back := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			f(i)
			back <- true
		}(i)
	}

	for i := 0; i < n; i++ {
		<-back
	}
}

// Pie executes f in parallel n times and collects errors
func Pie(n int, f func(int) error) error {
	errs := make([]error, n)
	hasError := false

	back := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			if !hasError {
				if err := f(i); err != nil {
					hasError = true
					errs[i] = err
				}
			}
			back <- true
		}(i)
	}

	for i := 0; i < n; i++ {
		<-back
	}

	if !hasError {
		return nil
	}

	merr := &MultiError{
		Msg:  "error while executing in parallel",
		Errs: []error{},
	}
	for _, err := range errs {
		if err != nil {
			merr.Errs = append(merr.Errs, err)
		}
	}
	return merr
}

// MultiError captures multiple errors.
type MultiError struct {
	Msg  string
	Errs []error
}

func (err *MultiError) Error() string {
	msg := err.Msg + ":"

	for _, e := range err.Errs {
		msg += "\n  " + e.Error()
	}

	return msg
}
