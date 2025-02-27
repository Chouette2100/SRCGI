package main

import (
    "fmt"
    // "io/ioutil"
    "os"
    "log"
    "regexp"

    "gopkg.in/yaml.v2"
)

type BotsList struct {
    Bots []string `yaml:"bots"`
}

func ReadBots() (re *regexp.Regexp) {
    // YAMLファイルを読み込む
    // data, err := ioutil.ReadFile("bots.yml")
    data, err := os.ReadFile("bots.yml")
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    // 構造体にアンマーシャル
    var bl BotsList
    err = yaml.Unmarshal(data, &bl)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    // 結果を表示
    fmt.Println("Bots:")
    for _, s := range bl.Bots {
        fmt.Println("-", s)
    }

    	pattern := ""
	for i, b := range bl.Bots {
		if i != 0 {
			pattern += "|"
		}
		pattern += b
	}
	re = regexp.MustCompile(pattern)


    return
}

