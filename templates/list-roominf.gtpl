<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<p>
<button type="button" onclick="location.href='top'">「イベント選択」画面に戻る</button><br>
</p>
<h2>配信者レベルとフォロワー数</h2>
<p style="padding-left:2em">
下記のデータは不定期あるいは定期的に取得したものです。レベルやフォロワー数の変化をすべて捕捉しているわけではありません。
</p>

<p style="padding-left:2em">
<a href="https://www.showroom-live.com/event/{{.Eventid}}">{{ .Eventname }}</a>（{{.Eventid}}）<br>
{{ .Period }}
</p>
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
</body>
</html>
