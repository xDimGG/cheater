package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/xdimgg/cheater/services/quizlet"
)

const (
	matchScore   = 5
	gravityScore = math.MaxUint32
	spellScore   = 0
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

		return func(content string) (err error) {
			content = strings.ToLower(content)

			if strings.Contains(content, "refresh") {
				ns, err := quizlet.New()
				if err != nil {
					return err
				}

				if err = s.Login(quizletUsername, quizletPassword); err != nil {
					return err
				}

				s = ns
				return nil
			}

			id := reQuizlet.FindString(content)
			if id == "" {
				return errNoMatch
			}

			methods := [...]struct {
				mode string
				run  func() error
			}{
				{"match", func() error { return s.UpdateHighScore(id, "match", matchScore) }},
				{"gravity", func() error { return s.UpdateHighScore(id, "gravity", gravityScore) }},
				// {"learn", nil},
				// {"write", nil},
				// {"spell", func() error { return s.EndSpellGame(id, spellScore) }},
			}

			for _, m := range methods {
				if strings.Contains(content, m.mode) {
					return m.run()
				}
			}

			for i, m := range methods {
				if i != 0 {
					time.Sleep(time.Second)
				}

				if err = m.run(); err != nil {
					return
				}
			}

			return nil
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
