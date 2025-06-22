package main

import (
	"SRCGI/ShowroomCGIlib"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadEntry() {
	// YAMLファイルパス
	yamlFilePath := "nontargetentry.yml"

	// YAMLファイルを読み込む
	yamlFile, err := os.ReadFile(yamlFilePath) // Go 1.16以降
	// Go 1.15以前の場合は以下を使用:
	// yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		log.Fatalf("YAMLファイルの読み込みに失敗しました: %v", err)
	}

	// 読み込んだYAMLデータを map[string]int にデコードする
	// デコード先の変数 NontargetHandler を指定
	err = yaml.Unmarshal(yamlFile, &ShowroomCGIlib.NontargetEntry)
	// ShowroomCGIlib.NontargetHandler は事前に定義されて
	if err != nil {
		log.Fatalf("YAMLデータのデコードに失敗しました: %v", err)
	}

	// 読み込み結果を確認
	fmt.Println("NontargetHandler がYAMLから初期化されました:")
	for key, value := range ShowroomCGIlib.NontargetEntry {
		fmt.Printf("  %s: %d\n", key, value)
	}

	// // ここから NontargetHandler を使った処理を記述...
	// // 例: 特定のキーの値を取得
	// if val, ok := ShowroomCGIlib.NontargetHandler["path1"]; ok {
	// 	fmt.Printf("path1 の値は %d です\n", val)
	// } else {
	// 	fmt.Println("path1 は NontargetHandler に含まれていません")
	// }
}
