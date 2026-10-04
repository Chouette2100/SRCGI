<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
{{if .TurnstileSiteKey}}
    <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    <div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9; margin: 20px auto;">
        <h3>セキュリティチェック</h3>
        {{if .TurnstileError}}
        <p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
        {{end}}
        <p>処理を続行するには、セキュリティチェックを完了してください。</p>
        <form method="POST" action="list-perday">
            <input type="hidden" name="eventid" value="{{.Eventid}}">
            <input type="hidden" name="requestid" value="{{.RequestID}}">
            <div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
            <br>
            <button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
        </form>
    </div>
{{else}}
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}

<p>
<table>
    <tr>
  <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
  <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
  <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
  <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
    </tr>
    <tr>
  <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
  <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
  <td><button type="button" onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button></td>
  <td></td>
    </tr>
  </table>


</p>
<h2>日々の獲得ポイント</h2>
<p style="padding-left:2em">
イベント全期間に渡ってデータを取得していない（＝欠測がある）ケースで日毎の獲得ポイントを算出できないときは数値が表示されません。
（欠測があっても日毎の獲得ポイントを算出できるときもあります）<br>
</p>
<p style="padding-left:2em;color: red;">
日々の獲得ポイントでは日のはじまりは「イベントパラメータの設定」にある「日々の獲得ポイントのリセット時刻」で指定した時刻とします。
デフォルトでは午前4時となります。グラフも同様です。
</p>


<p style="padding-left:2em">
<a href="https://www.showroom-live.com/event/{{.Eventid}}">{{ .Eventname }}</a>（{{.Eventid}}）<br>
{{ .Period }}
</p>
日々の獲得ポイント
<table border="1">
<tr align="center"><td></td>
{{ range .Longnamelist }}
<td>{{.Name}}</td>
{{ end }}
</tr>
{{ range .Pointrecordlist }}
<tr align="right">
<td>{{.Day}}</td>
{{ range .Pointlist }}
<td>{{.Spnt}}</td>
{{ end }}
</tr>
{{ end }}
</table>
<br>
累積獲得ポイント
<table border="1">
<tr align="center"><td></td>
{{ range .Longnamelist }}
<td>{{.Name}}</td>
{{ end }}
</tr>
{{ range .Pointrecordlist }}
<tr align="right">
<td>{{.Day}}</td>
{{ range .Pointlist }}
<td>{{.Tpnt}}</td>
{{ end }}
</tr>
{{ end }}
</table>
{{end}}
</body>
</html>
