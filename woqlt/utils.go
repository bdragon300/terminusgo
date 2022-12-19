package woqlt

import (
	"fmt"
	"strconv"
	"strings"
)

func paramVariadic[T any](p []callParam) []T {
	var res []T
	for _, v := range p {
		res = append(res, v.(T))
	}
	return res
}

func param1[T0 any](p []callParam) T0 {
	checkParamCount(p, []int{1})
	return p[0].(T0)
}

func param2[T0, T1 any](p []callParam) (T0, T1) {
	checkParamCount(p, []int{2})
	return p[0].(T0), p[1].(T1)
}

func param3[T0, T1, T2 any](p []callParam) (T0, T1, T2) {
	checkParamCount(p, []int{3})
	return p[0].(T0), p[1].(T1), p[2].(T2)
}

func param4[T0, T1, T2, T3 any](p []callParam) (T0, T1, T2, T3) {
	checkParamCount(p, []int{4})
	return p[0].(T0), p[1].(T1), p[2].(T2), p[3].(T3)
}

func param5[T0, T1, T2, T3, T4 any](p []callParam) (T0, T1, T2, T3, T4) {
	checkParamCount(p, []int{5})
	return p[0].(T0), p[1].(T1), p[2].(T2), p[3].(T3), p[4].(T4)
}

func checkParamCount(p []callParam, choices []int) {
	choicesStr := make([]string, 0)
	for _, v := range choices {
		if len(p) == v {
			return
		}
		choicesStr = append(choicesStr, strconv.Itoa(v))
	}
	msg := strings.Join(choicesStr, " or ")
	panic(fmt.Sprintf("Expected %s parameters, got %d", msg, len(p)))
}
