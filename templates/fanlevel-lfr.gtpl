<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>

<br>
このページはブックマーク可能です。
<br>
<br>
<a href="/fanlevel">トップページへ</a>
<br>
<br>
{{ .Roomname }}　（{{.Roomid}}）
<br>
<br>
（リスナー名をクリックするとそのリスナーの配信者ごとのレベルが表示されます）
<br>
<table border="1">
    <tr align="center">
        <td>草うま王</td>
        <td>ファン<br>レベル</td>
        <td>前月<br>レベル</td>
    </tr>
    {{ range .Lfr }}
    <tr>
        <td><a href="/fanlevel?userid={{ .User_id}}">{{ .User_name }}</a></td>

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
