
{{ $i := .Ncr }}
{{ range .Gslist }}
	<tr>
	<td align="right">
		{{ if eq .Orderno 999 }}
			---
		{{ else if ne .Orderno -1 }}
			{{ .Orderno }}
		{{ end }}
	</td>
	{{ range .Score }}
		<td align="right">
			{{ if eq . -1 }}
			n/a
			{{ else }}
			{{ Comma . }}
			{{ end }}
		</td>
	{{end}}
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.User_name}}</a>（{{ .Rank }}）</td>
	<td align="right">
		<a href="https://www.showroom-live.com/r/{{.Url}}">Live</a> <a href="/closedevents?userno={{.Userno}}&mode=0&path=5">イベント履歴</a>
	</td>
	</tr>
{{end}}
</table>
</body>
</html>
