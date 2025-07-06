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

    （<input type="month" size="8" min="2025-01" max="2025-06" name="yearmonth" value='{{ .Year }}-{{ FormatInt .Month "%02d" }}'>）
    貢献ポイントランキング（全イベント・全ルーム）表示範囲:
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

    <input type="submit" value="再表示" formaction="m-cntrbrank-listener" formmethod="GET" style="background-color: khaki">

    </p>
    </form>
    <p style="color: blue;">
    ※　該当するデータの一部が表示されていない可能性があります<br>
    ※　関連するページの作成に連動してこのページのレイアウトを変更していきます<br>
    ※　レスポンスは今後追加する機能のレスポンスを考慮しながら改善する予定です<br>
    ※　リスナーさんの名前は昔の名前になっていることがあります
    </p>
    <br>
    <table border="1">
        <tr allign="center">
            <td>貢献ポイント</td>
            <td>リスナー（ユーザーID）</td>
            <td>ルーム名・ルームIDと配信ページへのリンク</td>
            <td>イベント名・イベントIDとイベントページへのリンク・終了日時</td>
        </tr>
        {{ range .MonthlyCntrbRankList }}
            <tr>
            <td allign="right">{{ Comma .Point }}</td>
            <td allign="left">{{ .Listener }}（{{ .Lsnid }}）</td>
            <td allign="left">{{ .Longname }} <a href="https://www.showroom-live.com/{{ .Userid }}">{{ .Roomno }}</a></td>
            <td allign="left">
                {{ .Eventname }}
                <a href="https://showroom-live.com/event/{{ .Eventid }}">{{ .Ieventid }}</a>
                〜{{ FormatTime .Endtime "2006-01-02 15:04" }}
            </td>
            </tr>
        {{ end }}
    </table>
</body>
</html>
