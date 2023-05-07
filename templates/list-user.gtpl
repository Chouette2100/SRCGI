<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
{{/*	<meta http-equiv="refresh" content="{{.SecondsToReload}}; URL=list-last?eventid={{.Eventid}}">	*/}}
<html>
<body>
<p>
<button type="button" onclick="history.back()">結果表示選択画面に戻る</button><br>
</p>
<table border="1">
<tr align="center">
	<td>ルーム名</td>
	<td>アカウント</td>
	<td>ユーザーID</td>
	<td>ジャンル</td>
	<td>ランク</td>
	<td>レベル</td>
	<td>フォロワー数</td>
	<td>ポイント</td>
<tr>
{{ range . }}
<tr>
	<td>{{.Name}}</td>
	<td>{{.Account}}</td>
	<td>{{.ID}}</td>
	<td>{{.Genre}}</td>
	<td>{{.Rank}}</td>
	<td align="right">{{.Level}}</td>
	<td align="right">{{.Followers}}</td>
	<td align="right">{{.Point}}</td>
</tr>
{{ end }}
</table>
</body>
</html>
