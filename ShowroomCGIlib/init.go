package ShowroomCGIlib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	// "testing"
)

// IPRegion は地域コードとネットワークアドレスのペアを表す
type IPRegion struct {
	Region  string     // 地域コード (例: "JO", "JP")
	Network *net.IPNet // ネットワークアドレスとマスク (例: 193.203.110.0/23)
}

// テスト用のIPRegionListをロードするヘルパー関数
// ベンチマーク実行ごとにファイルを読み込むのは非効率なので、一度だけロードして使い回す
var IPRegionList []IPRegion

// loadIPRegionList は指定されたファイルからIPRegionのリストを読み込む
func loadIPRegionList(filename string) ([]IPRegion, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ipRegions []IPRegion
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // スペースで分割
		if len(parts) != 2 {
			// log.Printf("Warning: Skipping malformed line: %s", line)
			continue // 不正な行はスキップ
		}
		region := parts[0]
		cidrStr := parts[1]

		_, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			log.Printf("Warning: Failed to parse CIDR %s in line '%s': %v", cidrStr, line, err)
			continue
		}
		ipRegions = append(ipRegions, IPRegion{Region: region, Network: ipNet})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ipRegions, nil
}

func init() {
	// ベンチマーク実行前に一度だけIPRegionListをロード
	// 実際のファイル名に合わせて変更してください
	list, err := loadIPRegionList("cidr.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load IP region list: %v\n", err)
		os.Exit(1)
	}
	IPRegionList = list
	fmt.Printf("Loaded %d IP region entries.\n", len(IPRegionList))
}
