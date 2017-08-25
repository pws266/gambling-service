// command line arguments processing
package aux

import (
	"regexp"
	"strconv"
	"strings"
)

// command line parameters for
type InputTraits struct {
	PortNumber string

	DBName       string
	UserLogin    string
	UserPassword string
}

var InParam InputTraits = InputTraits{"8080", "", "", ""}

// getting arguments from specified command line
// cmdLine - command line without program name stored in 0th argument
//           input parameters should be in format: "key0 = value0 key1 = value1"
func (p *InputTraits) ReadArgs(cmdLine string) error {
	res := regexp.MustCompile("[ =\t]+").Split(cmdLine, -1)

	if res == nil {
		var err *Error = CreateError("ReadArgs", "No arguments found in command line")
		return err
	}

	if len(res)%2 != 0 {
		var err *Error = CreateError("ReadArgs", "Illegal command line format")
		return err
	}

	for i := 0; i < len(res); i++ {
		if strings.Compare("port", res[i]) == 0 {
			i++
			p.PortNumber = res[i]
		}

		if strings.Compare("login", res[i]) == 0 {
			i++
			p.UserLogin = res[i]
		}

		if strings.Compare("password", res[i]) == 0 {
			i++
			p.UserPassword = res[i]
		}

		if strings.Compare("db_name", res[i]) == 0 {
			i++
			p.DBName = res[i]
		}
	}

	_, err := strconv.Atoi(p.PortNumber)

	if err != nil {
		var err *Error = CreateError("ReadArgs", "Illegal port number value")
		return err
	}

	return nil
}
