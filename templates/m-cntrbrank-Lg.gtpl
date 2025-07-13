<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
{{/*
type MonthlyCntrbRank struct {
	YearMonth            string
	Year                 int
	Month                int
	Thpoint              int
	Thlist               []int
	MonthlyCntrbRankList []MonthlyCntrbRankData
}

type MonthlyCntrbRankData struct {
	Lsnid     int
	Listener  string
	Eventid   string
	Ieventid  int
	Eventname string
	Starttime time.Time
	Endtime   time.Time
	Roomno    int
	Userid    string
	Longname  string
	Point     int
}
*/}}
<html>
<body>
    <table>
        <tr>
            <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
            <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
        </tr>
    </table>
    
    <br>
    <form>
    <p style="color: black;">

    {{ .Year }}年 {{ .Month }}月 貢献ポイントランキング（リスナー/ルーム）<br>
    <br>
    　　年月　<input type="month" size="8" min="2025-01" max="2025-06" name="yearmonth" value='{{ .Year }}-{{ FormatInt .Month "%02d" }}'>
    　　最小pt　<input type="number" size="4" min="30000" max="300000" name="thpoint" value='{{ .Thpoint }}'> pt
    　　上位　<input type="number" size="4" min="3" max="50" name="limit" value='{{ .Limit }}'> 位まで
    {{/*
    貢献ポイント:
    <select name="thpoint">
        {{ $t := .Thpoint }}
        {{ range .Thlist }}
        <option
        {{ if eq . $t }}
            selected
        {{ end }}
        value="{{ . }}">{{ Comma . }}</option>
        {{ end }}
    </select>
    pt 以上
    */}}

    <input type="submit" value="再表示" formaction="m-cntrbrank-Lg" formmethod="GET" style="background-color: khaki">

    </p>
    </form>
    <p style="color: blue;">
    ※　指定できる年月は2024年10月から前月までです。<br>
    ※　該当するデータの一部が表示されていない可能性があります<br>
    ※　関連するページの作成に連動してこのページのレイアウトを変更していきます<br>
    ※　レスポンスは今後追加する機能のレスポンスを考慮しながら改善する予定です<br>
    ※　リスナーさんの名前は昔の名前になっていることがあります
    </p>
    <p style="color: red;">
    ※　「貢献pt合計」はあくまでこのリストにある貢献ptの合計です。<br>
    ※　ですので、データ取得の対象とする「最小pt」の設定によって「貢献pt合計」は変化します。<br>
    ※　実際の「貢献pt合計」はこの表にある値より大きくなります。
    </p>
    <br>
    <table border="1">
        <tr align="center">
            <td>リスナー（ユーザーID）</td>
            <td>貢献pt合計</td>
            <td>ルーム名・ルームIDと配信ページへのリンク</td>
            <td>貢献pt</td>
        </tr>
        {{ $sp := 0 }}
        {{ $id := 0 }}
        {{ range .MonthlyCntrbRankList }}
            <tr>
            {{ if and (eq .Spoint $sp) (eq .Lsnid $id ) }}
                <td></td>
                <td></td>
            {{ else }}
                <td align="left">{{ .Listener }}（{{ .Lsnid }}）</td>
                <td align="right">{{ Comma .Spoint }}</td>
                {{ $id = .Lsnid }}
                {{ $sp = .Spoint }}
            {{ end }}
            <td align="left">{{ .Longname }} <a href="https://www.showroom-live.com/{{ .Userid }}">{{ .Roomno }}</a></td>
            <td align="right">{{ Comma .Point }}</td>
            </tr>
        {{ end }}
    </table>
</body>
</html>
