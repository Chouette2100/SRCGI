<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>

<br>
このページはブックマーク可能です。
<br>
<br>
{{/*
<a href="/fanlevel">トップページへ</a>
<br>
<br>
*/}}
イベント：　{{ .Eventname }}　（{{.Eventid}}）
<br>
開催期間：　{{ .Period}}
<br>
<br>
{{/*
データ取得時刻：　{{ .Ts_lst}}　（次回データ取得予定時刻：　{{ .Ts_nxt }}
*/}}
データ取得時刻：　{{ .Ts_lst}}
<br>
{{/*
<br>
（リスナー名をクリックするとそのリスナーの配信者ごとのレベルが表示されます）
<br>
*/}}
<table border="1">
    <tr align="center">
        <td>ランク内<br>順位</td>
        <td>ルーム名（room_id）</td>
        <td>ランク</td>
        <td>ファン数</td>
        <td>獲得pt<br>順位</td>
    </tr>
    {{ range .RankingInfList }}
    {{ if eq .Room_id 164614 }}
    <tr style="color:red">
    {{ else }}
    {{ if eq .Irank 2 }}
    <tr style="color:blue">
    {{ else }}
    <tr>
    {{ end }}
    {{ end }}
        {{ if eq .Irorder -1 }}
        <td></td>
        {{ else }}
        <td align="right">{{ .Irorder }}</td>
        {{ end }}
        <td>{{ .Room_name}} ({{ .Room_id }})</td>
        <td>{{ .Srank}}</td>
        <td align="right">{{ .Fans}}</td>
        <td align="right">{{ .Iorder}}</td>
    </tr>
    {{end}}
 </table>
<br>
<br>
</body>
</html>
