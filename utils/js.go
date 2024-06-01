package utils

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
)

func GetAbogus(params string, userAgent string) (bogus string) {
	vm := goja.New()
	jsData, err := os.ReadFile("./javascript/a_bogus.js")
	if err != nil {
		return
	}

	_, err = vm.RunString(string(jsData))
	if err != nil {
		fmt.Println(1)
		return
	}

	genAbogus, ok := goja.AssertFunction(vm.Get("generate_a_bogus"))
	if !ok {
		fmt.Println(1)
		return
	}

	// query params  user-agent
	result, err := genAbogus(goja.Undefined(), vm.ToValue(params), vm.ToValue(userAgent))

	if err != nil {
		fmt.Println(1)
		return
	}
	bogus = result.String()

	fmt.Println(bogus)
	return
}

func GetXbogus(params string, userAgent string) {

	vm := goja.New()
	jsData, err := os.ReadFile("./javascript/x_bogus.js")
	if err != nil {
		fmt.Println(1)
		return
	}

	_, err = vm.RunString(string(jsData))
	if err != nil {
		fmt.Println(1)
		return
	}

	sign, ok := goja.AssertFunction(vm.Get("sign"))
	if !ok {
		fmt.Println(1)
		return
	}

	// query params  user-agent
	result, err := sign(goja.Undefined(), vm.ToValue(params), vm.ToValue(userAgent))

	if err != nil {
		fmt.Println(1)
		return
	}

	fmt.Println(result)
}
