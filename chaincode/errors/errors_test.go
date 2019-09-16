package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCases(t *testing.T) {
	strBassic := "this is an error msg"
	strFormat := "this is a complex msg %s %s"
	param1 := "Yeah"
	param2 := ";*)"
	strFinal := fmt.Sprintf(strFormat, param1, param2)
	testCases := []struct {
		desc           string
		errorFunc      func(...interface{}) Error
		args           []interface{}
		expectedMsg    string
		expectedStatus int
	}{
		{
			desc:           "E basic",
			errorFunc:      E,
			args:           []interface{}{strBassic},
			expectedMsg:    strBassic,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			desc:           "NotFound basic",
			errorFunc:      NotFound,
			args:           []interface{}{strBassic},
			expectedMsg:    strBassic,
			expectedStatus: http.StatusNotFound,
		},
		{
			desc:           "Conflict basic",
			errorFunc:      Conflict,
			args:           []interface{}{strBassic},
			expectedMsg:    strBassic,
			expectedStatus: http.StatusConflict,
		},
		{
			desc:           "BadRequest basic",
			errorFunc:      BadRequest,
			args:           []interface{}{strBassic},
			expectedMsg:    strBassic,
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc:           "E formating",
			errorFunc:      E,
			args:           []interface{}{strFormat, param1, param2},
			expectedMsg:    strFinal,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			desc:           "Conflict formating",
			errorFunc:      Conflict,
			args:           []interface{}{strFormat, param1, param2},
			expectedMsg:    strFinal,
			expectedStatus: http.StatusConflict,
		},
		{
			desc:      "E complex call",
			errorFunc: E,
			args: []interface{}{
				notFound,
				internal,
				E(strBassic),
				conflict,
				strFormat,
				param2,
				notFound,
				param1,
			},
			expectedMsg:    fmt.Sprintf(strFormat, param2, notFound, param1) + " " + strBassic,
			expectedStatus: http.StatusConflict,
		},
		{
			desc:           "NotFound simple wrapping",
			errorFunc:      NotFound,
			args:           []interface{}{fmt.Errorf(strBassic)},
			expectedMsg:    strBassic,
			expectedStatus: http.StatusNotFound,
		},
		{
			desc:           "BadRequest empty",
			errorFunc:      BadRequest,
			args:           []interface{}{},
			expectedMsg:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc:           "Forbidden overides",
			errorFunc:      Forbidden,
			args:           []interface{}{Internal()},
			expectedMsg:    "",
			expectedStatus: http.StatusForbidden,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			e := tC.errorFunc(tC.args...)
			assert.Error(t, e)
			assert.Equal(t, tC.expectedMsg, e.Error())
			assert.Equal(t, tC.expectedStatus, Wrap(e).HTTPStatusCode())
		})
	}
}
