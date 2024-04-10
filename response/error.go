package response

import "fmt"

type CustomError struct {
	Code int
	Msg  string
}

func (err CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", 500, "Internal Server Error")
}
