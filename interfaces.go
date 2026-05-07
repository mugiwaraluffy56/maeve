package maeve

import "context"

type Runner interface {
	Run(context.Context) error
}
