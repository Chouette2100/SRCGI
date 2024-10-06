
{{ $i := .Ncr }}
{{ $g := .Grid }}
{{ $c := .Cntrblst }}
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
			{{ else if eq . 0  }}
			---
			{{ else if lt . -1  }}
			（ {{ add . 10000 }} ）
			{{ else }}
			{{ Comma . }}
			{{ end }}
		</td>
	{{end}}
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.User_name}}</a>（{{ .Rank }}）</td>
	<td align="right">
		{{ if eq $c 1 }}
		<a href="/listgsc?giftid={{$g}}&userno={{.Userno}}">貢献ランキング</a>・
		{{ end }}
		<a href="https://www.showroom-live.com/r/{{.Url}}">Live</a>・
		<a href="/closedevents?userno={{.Userno}}&mode=0&path=5">イベント履歴</a>
	</td>
	</tr>
{{end}}
</table>
</body>
</html>
