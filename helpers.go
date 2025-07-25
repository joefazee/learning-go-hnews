package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) getUserFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(contextUserKey).(*User)
	if !ok {
		panic("no user found in context")
	}
	return u
}
