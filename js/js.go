package js

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/StackExchange/dnscontrol/models"

	"github.com/robertkrimen/otto"
	//load underscore js into vm by default
	_ "github.com/robertkrimen/otto/underscore"
)

//ExecuteJavascript accepts a javascript string and runs it, returning the resulting dnsConfig.
func ExecuteJavascript(script string, devMode bool) (*models.DNSConfig, error) {
	vm := otto.New()

	vm.Set("require", require)

	helperJs := GetHelpers(devMode)
	// run helper script to prime vm and initialize variables
	if _, err := vm.Run(helperJs); err != nil {
		return nil, err
	}

	// run user script
	if _, err := vm.Run(script); err != nil {
		return nil, err
	}

	// export conf as string and unmarshal
	value, err := vm.Run(`JSON.stringify(conf)`)
	if err != nil {
		return nil, err
	}
	str, err := value.ToString()
	if err != nil {
		return nil, err
	}
	conf := &models.DNSConfig{}
	if err = json.Unmarshal([]byte(str), conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func GetHelpers(devMode bool) string {
	return _escFSMustString(devMode, "/helpers.js")
}

func require(call otto.FunctionCall) otto.Value {
	file := call.Argument(0).String()
	fmt.Printf("requiring: %s\n", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	_, err = call.Otto.Run(string(data))
	if err != nil {
		panic(err)
	}
	return otto.TrueValue()
}
