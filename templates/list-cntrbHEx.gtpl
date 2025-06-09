 <button type="button" onclick="history.back()">「枠別貢献ポイント一覧表」画面に戻る</button><br>
<table border="1">
<tr style="text-align: center;">
<td>貢献ポイント</td>
<td>ルーム</td>
<td>イベント</td>
<td>開始日時</td>
<td>終了日時</td>
</tr>

{{ range . }}
	<tr>
	<td style="text-align: right;">{{ Comma .Point }}</td>
	{{/*
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{ .Roomno }}">{{ .Longname }}</a></td>
	*/}}
	<td><a href="/closedevents?userno={{ .Roomno }}&mode=0&path=5">{{ .Longname }}</a></td>
	{{/*
	<td><a href="https://www.showroom-live.com/event/{{ .Eventid }}">{{ .Eventname }}</a></td>
	*/}}
	<td><a href="/list-last?eventid={{ .Eventid }}">{{ .Eventname }}</a></td>
	<td>{{ FormatTime .Starttime "2006-01-02 15:04" }}</td>
	<td>{{ FormatTime .Endtime "2006-01-02 15:04" }}</td>
	</tr>
{{end}}
</table>
 <button type="button" onclick="history.back()">「枠別貢献ポイント一覧表」画面に戻る</button><br>
</body>
</html>
