
{{ $i := .Ncr }}
{{ range .Cntrbinflist }}
	<tr>
	<td align="right">
		{{ if eq .Ranking 999 }}
			---
		{{ else if ne .Ranking -1 }}
			{{ .Ranking }}
		{{ end }}
	</td>
	<td align="right">
		{{ if eq .Point -1 }}
			n/a
		{{ else }}
			{{ Comma .Point}}
		{{ end }}
	</td>
	{{ range .Incremental }}
		<td align="right">
			{{ if eq . -1 }}
			n/a
			{{ else }}
			{{ Comma . }}
			{{ end }}
		</td>
	{{end}}
	<td>{{.ListenerName}}</td>
	<td align="right">
		<a href="list-cntrbH?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid}}&ie={{ $i }}">{{.Tlsnid}}</a>
	</td>
	<td align="right">
		{{ if gt .Lsnid 0 }}
		<a href="list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid}}&ie={{ $i }}">{{.Tlsnid}}</a>
		{{ end }}
	</td>
	</tr>
{{end}}
</table>
</body>
</html>
