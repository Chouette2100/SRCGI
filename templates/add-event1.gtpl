<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
{{/*	<meta http-equiv="refresh" content="{{.SecondsToReload}}; URL=list-last?eventid={{.Eventid}}">	*/}}
<html>
<body>
<p>
<button type="button" onclick="location.href='top?eventid={{.Event_ID}}'">「SHOWROOMイベント結果表示」画面に戻る</button>
</p>
<table>
<tr><td>イベントID</td><td>{{.Event_ID}}</td></tr>
<tr><td>イベント名</td><td>{{.Event_name}}</td></tr>
<tr><td>期間</td><td>{{.Period}}</td></tr>
<tr><td>開始日時</td><td>{{.Start_time}}</td></tr>
<tr><td>終了日時</td><td>{{.End_time}}</td></tr>
</table>

