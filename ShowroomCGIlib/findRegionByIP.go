package ShowroomCGIlib

import (
	// "bufio"
	// "fmt"
	// "log"
	"net"
	// "os"
	// "strings"
	// "sync" // 並行処理のためのsyncパッケージを追加
)

// findRegionByIP はIPアドレス文字列から地域コードを検索する
func FindRegionByIP(ipStr string, ipRegionList []IPRegion) string {
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return "" // 無効なIPアドレス
	}

	// IPv4-mapped IPv6アドレス (例: ::ffff:192.0.2.1) をIPv4に変換
	// これにより、IPv4アドレスがIPv6形式でログに出力されていても正しく判定できる
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		parsedIP = ipv4
	}

	for _, ipRegion := range ipRegionList {
		if ipRegion.Network.Contains(parsedIP) {
			return ipRegion.Region
		}
	}
	return "" // 見つからなかった場合
}
