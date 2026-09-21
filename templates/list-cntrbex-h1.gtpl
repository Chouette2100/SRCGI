<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<head>
    {{if .TurnstileSiteKey}}
    <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    {{end}}
</head>
<body>
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}

{{if .TurnstileSiteKey}}
    <div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9;">
        <h3>セキュリティチェック</h3>
        {{if .TurnstileError}}
        <p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
        {{end}}
        <p>枠別貢献ポイント一覧表を表示するには、セキュリティチェックを完了してください。</p>
        <p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
        <form method="POST" action="list-cntrbex">
            <input type="hidden" name="eventid" value="{{.Eventid}}">
            <input type="hidden" name="userno" value="{{.Userno}}">
            <input type="hidden" name="requestid" value="{{.RequestID}}">
            <div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
            <br>
            <button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
        </form>
    </div>
    <br><br>
</body>
</html>
{{else}}

<table>
    <tr>
  <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
  <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
  <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
  <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
    </tr>
    <tr>
  <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
  <td></td>
  <td><button type="button" onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button></td>
  <td></td>
    </tr>
    <tr>
  <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
  <td></td>
  <td></td>
  <td></td>
    </tr>
  </table>

<p>枠別貢献ポイント一覧表</p>
{{/*
<p style="color:crimson;">本機能は実験的なものです。結果を100%信じないでください。<br>上位のリスナーの結果は比較的正確です。あくまで"比較的"にです。</p>
*/}}
<p style="color:crimson;">イベント開始から終了までのすべてのデータが取得されていない場合、<br>（特に最初と最後の）データに不整合が発生していることがあります。</p>
<p style="color:green;">2014-03-11以後に開始されてイベントについてはtlsnidはリスナーさんのユーザーIDです。</p>
<table>
<tr><td align="center"><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a>（{{.Eventid}}）</td></tr>
<tr><td align="center">{{.Period}}</td></tr>
<br>
<tr><td align="center"><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a>（{{.Userno}}）　<a href="https://www.showroom-live.com/event/contribution/{{ .Eventid}}?room_id={{.Userno}}">[公式]イベント貢献ランキング(100位まで)</a></td></tr>
</table>
<br>
{{/*
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
*/}}
{{end}}
