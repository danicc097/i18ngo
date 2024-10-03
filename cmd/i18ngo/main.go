package main

import (
	"fmt"
	"os"

	"github.com/danicc097/i18ngo"
)

func main() {
	// create fs.FS from cli arg of directory --> first arg.
	fs := os.DirFS(os.Args[1])
	pkgName := os.Args[2]

	data, err := i18ngo.GetTranslationData(fs, ".", pkgName)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "data: %v\n", data)

	src, err := i18ngo.Generate(data)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(src))
}
