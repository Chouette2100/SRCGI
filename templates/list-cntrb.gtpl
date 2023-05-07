{{ range . }}
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
		<a href="list-cntrbH?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid}}">{{.Tlsnid}}</a>
	</td>
	</tr>
{{end}}
</table>
</body>
</html>
