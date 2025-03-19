package util

import "strings"

func CheckError(e any, de any) {
	if e != nil {
		if de != nil {
			panic(de)
		}
		panic(e)
	}
}

type ContextKey struct {
	Name string
}

func (k *ContextKey) String() string {
	return "context value " + k.Name
}

func ContainString(list []string, a string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}
