package json

import (
	"encoding/json"
	opt "github.com/repeale/fp-go/option"
)

func one[T any]() func(input T) int {
	return func(input T) int {
		return 1
	}
}

//var message T
// https://qiita.com/nayuneko/items/2ec20ba69804e8bf7ca3
//if err := json.Unmarshal(input, &message); err != nil {
//	return opt.None[T]()
//} else {
//	return opt.Some[T](message)
//}
//fmt.Println("")

// https://text.baldanders.info/remark/2020/03/currying/
func TryParse[T any]() func(input []byte) opt.Option[T] {
	return func(input []byte) opt.Option[T] {
		var message T
		// https://qiita.com/nayuneko/items/2ec20ba69804e8bf7ca3
		if err := json.Unmarshal(input, &message); err != nil {
			return opt.None[T]()
		} else {
			return opt.Some[T](message)
		}
	}
}
