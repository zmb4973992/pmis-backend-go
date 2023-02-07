package main

import "fmt"

func main() {
	test1 := map[string]string{
		"a": "aaaa",
		"b": "bbbb",
	}

	test2 := map[string][]string{
		"a": {"aaaa", "bbbb", "cccc"},
		"b": {"aaaa", "bbbb", "cccc"},
	}

	fmt.Println(test1, test2)
}
