<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<p>
<button type="button" onclick="location.href='top?eventid={{.Eventid}}'">「SHOWROOMイベント結果表示」画面に戻る</button><br>
</p>
<h2>日々の獲得ポイント</h2>
<p style="padding-left:2em">
イベント全期間に渡ってデータを取得していない（＝欠測がある）ケースで日毎の獲得ポイントを算出できないときは数値が表示されません。
（欠測があっても日毎の獲得ポイントを算出できるときもあります）
</p>

<p style="padding-left:2em">
<a href="https://www.showroom-live.com/event/{{.Eventid}}">{{ .Eventname }}</a>（{{.Eventid}}）<br>
{{ .Period }}
</p>
日々の獲得ポイント
<table border="1">
<tr align="center"><td></td>
{{ range .Longnamelist }}
<td>{{.Name}}</td>
{{ end }}
</tr>
{{ range .Pointrecordlist }}
<tr align="right">
<td>{{.Day}}</td>
{{ range .Pointlist }}
<td>{{.Spnt}}</td>
{{ end }}
</tr>
{{ end }}
</table>
<br>
累積獲得ポイント
<table border="1">
<tr align="center"><td></td>
{{ range .Longnamelist }}
<td>{{.Name}}</td>
{{ end }}
</tr>
{{ range .Pointrecordlist }}
<tr align="right">
<td>{{.Day}}</td>
{{ range .Pointlist }}
<td>{{.Tpnt}}</td>
{{ end }}
</tr>
{{ end }}
</table>
</body>
</html>
