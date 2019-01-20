package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/xdimgg/cheater/services/quizlet"
)

const (
	quizletScore = 5
)

var (
	reQuizlet = regexp.MustCompile(`\d+`)

	errNoMatch = errors.New("no match found for the provided string")
)

type serviceFunc func(string) error

var serviceInitializers = map[string]func() (serviceFunc, error){
	"quizlet": func() (fn serviceFunc, err error) {
		s, err := quizlet.New()
		if err != nil {
			return
		}

		if err = s.Login(quizletUsername, quizletPassword); err != nil {
			return
		}

		return func(content string) error {
			id := reQuizlet.FindString(content)
			if id == "" {
				return errNoMatch
			}

			return s.UpdateLeaderboard(id, quizletScore)
		}, nil
	},
}

var services = make(map[string]serviceFunc)

func init() {
	var errsMu sync.Mutex
	var errs []error
	var wg sync.WaitGroup
	wg.Add(len(serviceInitializers))

	for name, init := range serviceInitializers {
		go func(name string, init func() (serviceFunc, error)) {
			defer wg.Done()

			fn, err := init()
			if err != nil {
				errsMu.Lock()
				errs = append(errs, fmt.Errorf("%s: %v", name, err))
				errsMu.Unlock()
				return
			}

			services[name] = fn
		}(name, init)
	}

	wg.Wait()

	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	if len(services) == 0 {
		os.Exit(1)
	}
}
