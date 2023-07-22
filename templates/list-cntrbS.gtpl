<table border="1">
	<tr align="center">
		<td>順位</td>
		<td>累計貢献ポイント</td>
		<td>貢献ポイント増分</td>
		<td>リスナー名</td>
		<td>前枠リスナー名</td>
	</tr>
	{{ range . }}
	<tr>
		<td align="right">{{ .Ranking }}</td>
		<td align="right">
			{{ if ne .Point -1 }}
			{{ Comma .Point }}
			{{ else }}
			n/a
			{{ end }}
		</td>
		<td align="right">
			{{ if ne .Point -1 }}
			{{ Comma .Incremental }}
			{{ else }}
			n/a
			{{ end }}
		</td>
		<td>{{.ListenerName}}</td>
		<td>{{.LastName}}</td>
	</tr>
	{{end}}
</table>
</body>

</html>