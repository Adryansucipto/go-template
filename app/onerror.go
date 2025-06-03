package app

import "go.uber.org/dig"

func OnError(container *dig.Container, err error) {
	err = container.Invoke(func() error {
		if err != nil {
			return err
		}
		return nil
	})
}
