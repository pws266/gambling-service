// errors processing
package aux

import (
	"fmt"
)

type Error struct {
	objectName  string
	description string
}

func (p *Error) Error() string {
	return fmt.Sprintf("\nOoops! Error occured!\n   In: %s\n   Problem: %s\n",
		p.objectName, p.description)
}

func CreateError(where string, description string) *Error {
	return &Error{where, description}
}

func CreateExternalError(where string, description string, external error) *Error {
	externalMsg := fmt.Sprintf("\n   External: %s\n", external)
	return &Error{where, description + externalMsg}
}
