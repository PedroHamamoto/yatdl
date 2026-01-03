package httperr

import "net/http"

var (
	BadRequest = Error{
		http.StatusBadRequest,
		400001,
		"Invalid request. Check your data and try again",
	}

	Unauthorized = Error{
		http.StatusUnauthorized,
		401001,
		"Invalid credentials",
	}

	ExpiredToken = Error{
		http.StatusUnauthorized,
		401002,
		"Session expired. Please, log in again",
	}

	Forbidden = Error{
		http.StatusForbidden,
		403001,
		"Your are not allowed to perform this action",
	}

	TaskNotFound = Error{
		http.StatusNotFound,
		404001,
		"Task not found. Check the ID and try again",
	}

	Unknown = Error{
		http.StatusInternalServerError,
		500001,
		"An unexpected error happened. If the problem persists, contact our support",
	}
)
