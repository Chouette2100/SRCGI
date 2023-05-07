<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>

<button type="button" onclick="location.href='list-last?eventid={{.Eventid}}&userno={{.Userno}}'">「直近の獲得ポイント一覧」画面に戻る</button>
<br>
<button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}'">「枠別貢献ポイント一覧」画面に戻る</button>
<br><br>
<p>枠別貢献ポイント一覧表</p>
<table>
<tr><td align="center"><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a>（{{.Eventid}}）</td></tr>
<tr><td align="center">{{.Period}}</td></tr>
<br>
<tr><td align="center"><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a>（{{.Userno}}）</td></tr>
<br>
<tr><td align="center">{{ .S_stime }} ----- {{ .S_etime }} （ {{ .Ifrm1 }} ）</td></tr>
</table>
<br>
<table>
    <tr>
        <td width="400" align="left">
        {{ if ne .Ifrm_b -1 }}
            {{ if eq .Srt 1 }}
            <button type="button" onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_b}}&sort=D'">一つ前の配信枠へ</button>
            {{ else }}
            <button type="button" onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_b}}'">一つ前の配信枠へ</button>
            {{ end }}
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td width="400" align="right">
        {{ if ne .Ifrm_f -1 }}
            {{ if eq .Srt 1 }}
            <button type="button" onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_f}}&sort=D'">次の配信枠へ</button>
            {{ else }}
            <button type="button" onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_f}}'">次の配信枠へ</button>
            {{ end }}
        {{ else }}
            -----------
        {{ end }}
        </td>
    </tr>
</table>