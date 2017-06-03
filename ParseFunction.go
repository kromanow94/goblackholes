package goblackholes

import (
	//"errors"
	"github.com/Knetic/govaluate"
	"log"
	. "math"
	"reflect"
	"strconv"
	"strings"
	"fmt"
)

var functions map[string]govaluate.ExpressionFunction

func init() {
	fmt.Println("parseFunction init")
	functions = make(map[string]govaluate.ExpressionFunction)
	functions["strlen"] = func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return (float64)(length), nil
	}
	functions["sin"] = func(args ...interface{}) (interface{}, error) {
		return Sin(args[0].(float64)), nil
	}
	functions["cos"] = func(args ...interface{}) (interface{}, error) {
		return Cos(args[0].(float64)), nil
	}
	functions["pow"] = func(args ...interface{}) (interface{}, error) {
		return Pow(args[0].(float64), args[1].(float64)), nil
	}
	functions["PI"] = func(args ...interface{}) (interface{}, error) {
		return Pi, nil
	}
}

/* statistics of functions:
n of iter:  314160
EvalWithParams2:
7534605279
EvalWithParams:
9022396575
ParseFunction:
6428727151
NewEval with const:
8585686337
Normal function:
45168146
*/

/// @funcStr is string. It may have maximum of three parameters. Parameters must be signed as "x", "y", "z"
/// @args - at least two arguments: 1.: channel; 2.: float64
/// example:
/// go ParseFunction(funcStr, channelString, i, i)
/// go EvaluateFunction(<-channelString, channelEval)
func ParseFunction(str string, args ...interface{}) {
	if (reflect.TypeOf(args[0]).String()) != "chan string" {
		log.Fatal("First argument must be a channel.")
		return
	}
	if len(args) < 2 {
		log.Fatal("You must have at least one channel and one value argument.")
		return
	}
	if len(args) > 4 {
		log.Fatal("You can't have more than 3 value arguments.")
		return
	}

	for i := 1; i < len(args); i++ {
		if reflect.TypeOf(args[i]).String() != "float64" {
			log.Fatal("Argument " + strconv.FormatInt(int64(i+1), 10) + " is not float64")
			return
		}
	}

	valMap := []string{"x", "y", "z"}
	var result string
	for i := 1; i < len(args); i++ {
		value := args[i].(float64)
		result = ""
		sliceStr := strings.Split(str, valMap[i-1])
		for j := 0; j < len(sliceStr); j++ {
			result += sliceStr[j]
			if j+1 < len(sliceStr) {
				result += strconv.FormatFloat(value, 'f', -1, 64)
			}
		}
		str = result
	}
	args[0].(chan string) <- result
	return
}

func EvaluateFunction(str string, out chan float64) {
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(str, functions)
	val, _ := expression.Evaluate(nil)
	out <- val.(float64)
}

/// args must be a pair of value name (string) and value (float64)
/// example:
/// go EvaluateWithParameters(funcStr, channelEval, "x", i, "y", i)
func EvaluateWithParameters(str string, out chan float64, args ...interface{}) {
	pairLen := len(args) / 2
	params := make(map[string]interface{}, pairLen)
	for i := 0; i < pairLen*2; i += 2 {
		params[args[i].(string)] = args[i+1].(float64)
	}
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(str, functions)
	val, _ := expression.Evaluate(params)
	out <- val.(float64)
}

/// example:
/// parameters := make(map[string]interface{}, 1)
/// parameters["x"] = i;
/// parameters["y"] = i;
/// go EvaluateWithParameters2(funcStr, parameters, channelEval)
/// go FlushChannel(channelEval)
func EvaluateWithParameters2(str string, params map[string]interface{}, out chan float64) {
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(str, functions)
	val, _ := expression.Evaluate(params)
	out <- val.(float64)
}

func FlushChannel(c chan float64) {
	<-c

}
