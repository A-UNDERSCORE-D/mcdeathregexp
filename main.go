package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	targetJar    = flag.String("mcjar", "./minecraft.jar", "Sets the minecraft jar to extract messages from")
	assetPath    = flag.String("langpath", "assets/minecraft/lang/en_us.json", "sets the path to the lang file in the jar")
	noRegexp     = flag.Bool("noregex", false, "output the raw lang entries, dont turn them into a regexp")
	targetPrefix = flag.String("prefix", "death.", "the prefix to filter for in the lang files")
	ignoreKeys   = flag.String("ignorekeys", "death.attack.badRespawnPoint.link", "keys in the lang files to ignore (comma sep.)")
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	flag.Parse()
	langJson, err := getLangFile(*targetJar)
	checkErr(err)

	res, err := extractDeathMessages(langJson)
	checkErr(err)

	if *noRegexp {
		for _, v := range res {
			fmt.Println(v)
		}
	} else {
		fmt.Println(regexpify(res))
	}

}

func getLangFile(filePath string) ([]byte, error) {
	zipped, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read jar: %s", err)
	}

	defer zipped.Close()

	for _, f := range zipped.File {
		if f.Name == *assetPath {
			reader, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("could not open lang file for reading %w", err)
			}

			defer reader.Close()

			res, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, fmt.Errorf("error while reading lang file: %w", err)
			}

			return res, nil
		}
	}

	return nil, errors.New("could not find target asset file in jar")
}

func matchesAny(x string, y []string) bool {
	for _, v := range y {
		if v == x {
			return true
		}
	}
	return false
}

func extractDeathMessages(langFile []byte) ([]string, error) {
	var data map[string]string
	if err := json.Unmarshal(langFile, &data); err != nil {
		return nil, fmt.Errorf("could not parse json: %w", err)
	}

	var out []string

	for k, v := range data {
		if !strings.HasPrefix(k, *targetPrefix) || matchesAny(k, strings.Split(*ignoreKeys, ",")) {
			continue
		}

		out = append(out, v)
	}

	sort.StringSlice(out).Sort()

	return out, nil
}

var replacer = strings.NewReplacer(
	`%1$s`, `\S+`,
	`%2$s`, `\S+`,
	`%3$s`, `.*`,
	`%s`, `.*`,
)

func regexpify(in []string) string {
	// One day Ill dedupe these programmatically. Today isnt that day.

	out := strings.Builder{}
	out.WriteString("(")
	first := true
	for _, msg := range in {
		if first {
			first = false
		} else {
			out.WriteString("|")
		}

		toAdd := replacer.Replace(msg)
		out.WriteString(toAdd)
	}
	out.WriteString(")")
	return out.String()
}
