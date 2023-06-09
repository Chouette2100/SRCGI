<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<button type="button" onclick="location.href='top'">イベント選択画面に戻る</button>　
<button type="button" onclick="location.href='currentevent'">開催中イベント一覧表</button>　
<button type="button" onclick="location.href='list-last?eventid={{.Eventid}}&userno={{.Userno}}'">「直近の獲得ポイント一覧」画面に戻る</button>
<br><br>
<p>枠別貢献ポイント一覧表</p>
<p style="color:crimson;">本機能は実験的なものです。結果を100%信じないでください。<br>上位のリスナーの結果は比較的正確です。あくまで"比較的"にです。</p>
<p style="color:crimson;">イベント開始から終了までのすべてのデータが取得されていない場合、<br>（特に最初と最後の）データに不整合が発生していることがあります。</p>
<table>
<tr><td align="center"><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a>（{{.Eventid}}）</td></tr>
<tr><td align="center">{{.Period}}</td></tr>
<br>
<tr><td align="center"><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a>（{{.Userno}}）　<a href="https://www.showroom-live.com/event/contribution/{{ .Eventid}}?room_id={{.Userno}}">[公式]イベント貢献ランキング(100位まで)</a></td></tr>
</table>
<br>
<table>
    <tr>
        <td>
        {{ if ne .Nft -1 }}
            <button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{.Nft}}'">先頭に戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Npb -1 }}
            <button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{.Npb}}'">１ページ戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .N1b -1 }}
            <button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{.N1b}}'">一枠分戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .N1f -1 }}
            <button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{.N1f}}'">一枠分進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Npf -1 }}
            <button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{.Npf}}'">１ページ進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Nlt -1 }}
            <button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{.Nlt}}'">最後に進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
    </tr>
</table>