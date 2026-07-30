<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}


<br>
このページはブックマーク可能です。
<br>
<br>
<a href="/fanlevel">トップページへ</a>
<br>
<br>
{{ .Username }}　（{{ .Userid }}）
<br>
<br>
（ルーム名をクリックするとそのルームにいる草うま王のレベルが表示されます）
<br>
<table border="1">
    <tr align="center">
        <td>ルーム名</td>
        <td>ファン<br>レベル</td>
        <td>前月<br>レベル</td>
    </tr>
    {{ range .Levellist }}
    <tr>
        <td><a href="/fanlevel?roomid={{ .Room_id }}">{{ .Room_name }}</a></td>
        <td align="right">
    	    {{ if eq .Level -1 }}
			    n/a
		    {{ else }}
			    {{ .Level }}
		    {{ end }}
        </td>
        <td align="right">
    	    {{ if eq .Level_lst -1 }}
			    n/a
		    {{ else }}
			    {{ .Level_lst }}
		    {{ end }}
        </td>
    </tr>
    {{end}}
 </table>
<br>
<br>
</body>
</html>
