
{{ $i := .Ncr }}
{{ range .Vgslist }}
	<tr>
	<td align="right">
		{{ if eq .Orderno 0 }}
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
	<td>{{ .Viewername }}</td>
	<td align="right">{{ .Viewerid }}</td>
	</tr>
{{end}}
</table>
</body>
</html>
