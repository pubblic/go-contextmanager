package contextmanager

import (
	"context"
	"reflect"
)

type Signal struct {
	o   interface{}
	err error
}

func SignalTx() chan Signal {
	return make(chan Signal, 1)
}

func (s *Signal) Object() interface{} { return s.o }

func (s *Signal) Err() error { return s.err }

func (s *Signal) As(v interface{}) bool {
	if s.err != nil {
		return false
	}
	ref := reflect.ValueOf(v)
	ref.Elem().Set(reflect.ValueOf(s.o))
	return true
}

func (s *Signal) ErrorAs(co *SignalContext, v interface{}) bool {
	if s.err != nil {
		return co.Error(s.err)
	}
	ref := reflect.ValueOf(v)
	ref.Elem().Set(reflect.ValueOf(s.o))
	return true
}

type SignalContext struct {
	ctx  context.Context
	sigc chan<- Signal
}

func NewSignalContext(ctx context.Context, sigc chan<- Signal) SignalContext {
	return SignalContext{
		ctx:  ctx,
		sigc: sigc,
	}
}

func (co *SignalContext) Err() error {
	return co.ctx.Err()
}

func (co *SignalContext) IsDone() bool {
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

func (co *SignalContext) Yield(o interface{}) bool {
	if co.IsDone() {
		return false
	}
	select {
	case co.sigc <- Signal{o: o}:
		return true
	case <-co.ctx.Done():
		return false
	}
}

func (co *SignalContext) Error(err error) bool {
	if co.IsDone() {
		return false
	}
	if err == nil {
		return true
	}
	select {
	case co.sigc <- Signal{err: err}:
		return true
	case <-co.ctx.Done():
		return false
	}
}

func (co *SignalContext) Fatal(err error) error {
	if co.IsDone() {
		return co.Err()
	}
	if err == nil {
		return err
	}
	select {
	case co.sigc <- Signal{err: err}:
		return err
	case <-co.ctx.Done():
		return co.Err()
	}
}

func (co *SignalContext) Exit(err error) error {
	return co.Fatal(err)
}

type ObjectContext struct {
	ctx  context.Context
	sigc chan<- interface{}
}

func ObjectTx() chan interface{} {
	return make(chan interface{}, 1)
}

func NewObjectContext(ctx context.Context, sigc chan<- interface{}) ObjectContext {
	return ObjectContext{
		ctx:  ctx,
		sigc: sigc,
	}
}

func (co *ObjectContext) Err() error {
	return co.ctx.Err()
}

func (co *ObjectContext) IsDone() bool {
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

func (co *ObjectContext) Yield(o interface{}) bool {
	if co.IsDone() {
		return false
	}
	select {
	case co.sigc <- o:
		return true
	case <-co.ctx.Done():
		return false
	}
}

type ErrorContext struct {
	ctx  context.Context
	sigc chan<- error
}

func ErrorTx() chan error {
	return make(chan error, 1)
}

func NewErrorContext(ctx context.Context, sigc chan<- error) ErrorContext {
	return ErrorContext{
		ctx:  ctx,
		sigc: sigc,
	}
}

func (co *ErrorContext) Err() error {
	return co.ctx.Err()
}

func (co *ErrorContext) IsDone() bool {
	select {
	case <-co.ctx.Done():
		return false
	default:
		return true
	}
}

func (co *ErrorContext) Error(err error) bool {
	if co.IsDone() {
		return false
	}
	if err == nil {
		return true
	}
	select {
	case co.sigc <- err:
		return true
	case <-co.ctx.Done():
		return false
	}
}

func (co *ErrorContext) Fatal(err error) error {
	if co.IsDone() {
		return co.Err()
	}
	if err == nil {
		return err
	}
	select {
	case co.sigc <- err:
		return err
	case <-co.ctx.Done():
		return co.Err()
	}
}

func (co *ErrorContext) Exit(err error) error {
	return co.Fatal(err)
}
