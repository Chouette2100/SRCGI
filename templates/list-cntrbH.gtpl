<table border="1">
<tr>
<td>配信開始時刻</td>
<td>配信終了時刻</td>
<td>目標値(推定)</td>
<td>貢献ポイント</td>
<td>達成状況</td>
<td>累計ポイント</td>
<td>リスナー名（変更履歴）</td>
<td>突き合わせ状況</td>
</tr>

{{ range . }}
	<tr>
	<td>{{.S_stime}}</td>
	<td>{{.S_etime}}</td>

	{{/*
	<td align="right">
		{{ if lt .Target 0 }}
			n/a
		{{ else }}
			{{ Comma .Target }}
		{{ end }}
	</td>
	*/}}
	<td align="right">---</td>

	<td align="right">
		{{ if eq .Incremental -1 }}
			---
		{{ else }}
			{{ Comma .Incremental }}
		{{ end }}
	</td>
	
	{{/*
	<td align="right">
		{{ if or ( eq .Incremental -1) (lt .Target 0 ) }}
			---
		{{ else }}
			{{ Comma (sub .Incremental .Target) }}
		{{ end }}
	</td>
	*/}}
	<td align="right">---</td>

	<td align="right">
		{{ if lt .Point 0 }}
			---
		{{ else }}
			{{ Comma .Point }}
		{{ end }}
	</td>
	<td>{{.Listener}}</td>
	<td>{{.Lastname}}</td>
	</tr>
{{end}}
</table>
</body>
</html>
