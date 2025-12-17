<!DOCTYPE html>
<html>

<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    {{/* Turnstile 1 */}}
    {{if .TurnstileSiteKey}}
        <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    {{end}}
    {{/* ----------- */}}
</head>

<body>
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
            <td><button type="button"
                    onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
            </td>
            <td></td>
        </tr>
        <tr>
            <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
            <td><button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{ .Ie }}'">枠別貢献ポイント</button></td>
            <td></td>
            <td></td>
        </tr>
        </table>
    </table>
    <br>

{{if .TurnstileSiteKey }}
                <!-- Turnstileチャレンジ表示 -->
                        <div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9;">
                                <h3>セキュリティチェック</h3>
                                {{if .TurnstileError}}
                                <p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
                                {{end}}
                                <p>処理を続行するには、セキュリティチェックを完了してください。</p>
                                <p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
                                <form method="POST" action="list-cntrbH">
                                        <input type="hidden" name="eventid" value="{{.Eventid}}">
                                        <input type="hidden" name="userno" value="{{.Userno}}">
                                        <input type="hidden" name="tlsnid" value="{{.Tlsnid}}">
                                        <input type="hidden" name="ie" value="{{.Ie}}">

                                        <input type="hidden" name="requestid" value="{{.RequestID}}">
                                        <div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
                                        <br>
                                        <button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
                                </form>
                        </div>
                        <br><br>

{{else}}

    <p>枠別貢献ポイント一覧表</p>
    <table>
        <tr>
            <td align="center"><a
                    href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a>（{{.Eventid}}）</td>
        </tr>
        <tr>
            <td align="center">{{.Period}}</td>
        </tr>
        <br>
        <tr>
            <td align="center"><a
                    href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a>（{{.Userno}}）
            </td>
        </tr>
        <br>
        <tr>
            <td align="center">“{{ .Listener }}” （ Tlsnid = {{ .Tlsnid }} ）</td>
        </tr>
    </table>
    <br>
    <table>
        <tr>
            <td width="400" align="left">
                {{ if ne .Tlsnid_b -1 }}
                <button type="button"
                    onclick="location.href='list-cntrbH?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid_b}}'">{{
                    .Listener_b }}</button>
                {{ else }}
                -----------
                {{ end }}
            </td>
            <td width="400" align="right">
                {{ if ne .Tlsnid_f -1 }}
                <button type="button"
                    onclick="location.href='list-cntrbH?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid_f}}'">{{
                    .Listener_f }}</button>
                {{ else }}
                -----------
                {{ end }}
            </td>
        </tr>
    </table>

{{/* ---------------------------------------------------------- */}}
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

{{ range .CntrbHistory }}
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
{{ end }}
</body>
</html>
