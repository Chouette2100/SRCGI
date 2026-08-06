
<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
	{{/*
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}
	*/}}


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
	<td>{{.ListenerName}}</td>
	<td align="right">
		{{ if gt .Lsnid 0 }}
		<a href="list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid}}&name={{ .ListenerName }}">{{.Tlsnid}}</a>
		{{ end }}
	</td>
	</tr>
{{end}}
</table>
</body>
</html>
