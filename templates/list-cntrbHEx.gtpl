 <button type="button" onclick="history.back()">「枠別貢献ポイント一覧表」画面に戻る</button><br>
<table border="1">
<tr style="text-align: center;">
<td>貢献ポイント</td>
<td>ルーム(「イベント獲得ポイントの履歴」へのリンク)</td>
<td style="border-right: none;">イベント(「直近の獲得ポイント一覧」と「グラフ」へのリンク)</td>
<td style="border-left: none;"></td>
<td>開始日時</td>
<td>終了日時</td>
</tr>

{{ range . }}
	{{ $e :=  FormatTime .Endtime "2006-01-02 15:04" }}
	{{ if lt $e "2024-09-29 00:00" }}
	<tr style="background-color: silver">
	{{ else if gt $e (FormatTime .Stnow "2006-01-02 15:04") }}
	<tr style="background-color: yellow">
	{{ else }}
	<tr>
	{{ end }}
	<td style="text-align: right;">{{ Comma .Point }}</td>
	{{/*
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{ .Roomno }}">{{ .Longname }}</a></td>
	*/}}
	<td><a href="/closedevents?userno={{ .Roomno }}&mode=0&path=5">{{ .Longname }}</a></td>
	{{/*
	<td><a href="https://www.showroom-live.com/event/{{ .Eventid }}">{{ .Eventname }}</a></td>
	*/}}
	<td style="border-right: none;"><a href="/list-last?eventid={{ .Eventid }}">{{ .Eventname }}</a></td>
	<td style="border-left: none;">
	  <a href="graph-total?eventid={{.Eventid}}">グラフ</a>
	  {{/*
	  <button type="button"
        onclick="location.href='graph-total?eventid={{.Eventid}}'">グラフ</button>
	  */}}
	</td>
	<td>{{ FormatTime .Starttime "2006-01-02 15:04" }}</td>
	<td>{{ $e }}</td>
	</tr>
{{end}}
</table>
 <button type="button" onclick="history.back()">「枠別貢献ポイント一覧表」画面に戻る</button><br>
</body>
</html>
