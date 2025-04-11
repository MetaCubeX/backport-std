package context

import "context"

var Cause = context.Cause
var WithCancel = context.WithCancel
var WithCancelCause = context.WithCancelCause
var WithDeadline = context.WithDeadline
var WithTimeout = context.WithTimeout

//var WithTimeoutCause = context.WithTimeoutCause
//var WithDeadlineCause = context.WithDeadlineCause

type CancelCauseFunc = context.CancelCauseFunc
type CancelFunc = context.CancelFunc
type Context = context.Context

var Background = context.Background
var TODO = context.TODO
var WithValue = context.WithValue

//var WithoutCancel = context.WithoutCancel
