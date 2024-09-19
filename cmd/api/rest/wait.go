package rest

import "fmt"

func (a *Application) Background(fn func()) {
	a.Wg.Add(1)

	go func() {
		defer a.Wg.Done()

		defer func() {
			if err := recover(); err != nil {
				a.Logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}
