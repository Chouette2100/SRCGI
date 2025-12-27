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
        {{ if ne .Eventid "" }}
        <tr>
            <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
            <td></td>
            <td><button type="button"
                    onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
            </td>
            <td></td>
        </tr>
        {{ if ne .Userno 0 }}
        <tr>
            <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
            <td><button type="button" onclick="location.href='list-cntrbex?eventid={{.Eventid}}&userno={{.Userno}}&ie={{ .Ie }}'">枠別貢献ポイント</button></td>
            <td></td>
            <td></td>
        </tr>
        {{ end}}
        {{ end}}
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
        <form method="POST" action="list-cntrbHEx">
            <input type="hidden" name="eventid" value="{{.Eventid}}">
            <input type="hidden" name="userno" value="{{.Userno}}">
            <input type="hidden" name="tlsnid" value="{{.Tlsnid}}">
            <input type="hidden" name="ie" value="{{.Ie}}">
            <input type="hidden" name="name" value="{{.Name}}">

            <input type="hidden" name="requestid" value="{{.RequestID}}">
            <div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
            <br>
            <button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
        </form>
    </div>
    <br><br>
{{else}}
    <p style="color: blue;">リスナーさんの過去のイベントでの貢献ポイントの履歴です<br>
    データは貢献100位までが発表されるイベント貢献ランキングから取得しています。<br>
    つまり貢献ポイントが1ptでもこのリストにあることもありますが、貢献ポイント10,000ptでもないことがあります。<br>
    </p>
    <p style="color: blue;">
    終了したイベントのこのデータは特定の日時（注１）以後では、イベント獲得ポイントの取得を指定したすべてのイベント・ユーザーついて表示されます。<br>
    通常獲得ポイントが上位・中位のルームは自動的にイベント獲得ポイント取得の対象となりますので、このデータも得られます。<br>
    この場合リスナー別/枠別貢献ポイントの取得を指定したか否かは無関係です。<br>
    開催中のイベントについてはこの機能は現時点ではリスナー別/枠別貢献ポイントの取得を指定しておく必要がありますが<br>
    今大急ぎで開催中イベントでも使えるように改修中です。<br>
    <BR>
    特定の日時（注１）以前では、リスナー別/枠別貢献ポイントの取得を指定したイベント・ルームのみが表示の対象となります。<br>
    注１　「特定の日時」とは現時点では2024年9月29日です。<br>
    この日時はなんらかの事情あるいはご要望により今後数ヶ月から数年さかのぼった日時に変更する可能性はあります。
    </p>
    <p style="color: blue;">
    なお、5月15日以後については最終枠の枠別貢献ポイントが取得できなくなっていましたが、これについては6月9日18時までにすべて修復しました。
    </p>
    <br>
    <div style="background: #ffe4e1; border: #ffe4e1 solid 2px; font-size: 100%; padding: 20px; width: 50em;">
    ※　イベント・ルームに対する貢献ポイントの月単位のランキングを作ってみました。<br><br>
    　　<a href="m-cntrbrank-listener">月別イベント・リスナー貢献ポイントランキング</a>（結果が表示されるまで数十秒かかります）<br><br>
    ※　これも貢献ランキングですがどのリスナーさんがどのルームを応援しているかわかるように作ったものです。<br><br>
    　　<a href="m-cntrbrank-Lg">月別貢献ポイントランキング（リスナー/ルーム）</a>（結果が表示されるまで数十秒かかります）
    </div>
    <br>
    <table>
    {{/*
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
    */}}
        <tr>
            <td align="center">“{{ .Listener }}” （ Tlsnid = {{ .Tlsnid }} ）</td>
        </tr>
    </table>
    <br>
    {{/*
    <table>
        <tr>
            <td width="400" align="left">
                {{ if ne .Tlsnid_b -1 }}
                <button type="button"
                    onclick="location.href='list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid_b}}'">{{
                    .Listener_b }}</button>
                {{ else }}
                -----------
                {{ end }}
            </td>
            <td width="400" align="right">
                {{ if ne .Tlsnid_f -1 }}
                <button type="button"
                    onclick="location.href='list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid_f}}'">{{
                    .Listener_f }}</button>
                {{ else }}
                -----------
                {{ end }}
            </td>
        </tr>
    </table>
    */}}

{{/*
    {{ $n := .Stnow }}
*/}}
   
{{/* ---------------------------------------------- */}}

 <button type="button" onclick="history.back()">「枠別貢献ポイント一覧表」画面に戻る</button><br>
<table border="1">
<tr style="text-align: center;">
<td>貢献ポイント</td>
<td>ルーム(「イベント獲得ポイントの履歴」へのリンク)</td>
<td style="border-right: none;">イベント(「直近の獲得ポイント一覧」と「グラフ」へのリンク)</td>
<td style="border-left: none;"></td>
<td>開始日時</td>
<td>終了日時</td>
</tr>

{{ range .CntrbhistoryEx }}
	{{ $e :=  FormatTime .Endtime "2006-01-02 15:04" }}
	{{ if lt $e "2024-09-29 00:00" }}
	<tr style="background-color: silver">
	{{ else if gt $e (FormatTime .Stnow "2006-01-02 15:04") }}
	<tr style="background-color: yellow">
	{{ else }}
	<tr>
	{{ end }}
	<td style="text-align: right;">{{ Comma .Point }}</td>
	{{/*
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{ .Roomno }}">{{ .Longname }}</a></td>
	*/}}
	<td><a href="/closedevents?userno={{ .Roomno }}&mode=0&path=5">{{ .Longname }}</a></td>
	{{/*
	<td><a href="https://www.showroom-live.com/event/{{ .Eventid }}">{{ .Eventname }}</a></td>
	*/}}
	<td style="border-right: none;"><a href="/list-last?eventid={{ .Eventid }}">{{ .Eventname }}</a></td>
	<td style="border-left: none;">
	  <a href="graph-total?eventid={{.Eventid}}">グラフ</a>
	  {{/*
	  <button type="button"
        onclick="location.href='graph-total?eventid={{.Eventid}}'">グラフ</button>
	  */}}
	</td>
	<td>{{ FormatTime .Starttime "2006-01-02 15:04" }}</td>
	<td>{{ $e }}</td>
	</tr>
{{end}}
</table>
 <button type="button" onclick="history.back()">「枠別貢献ポイント一覧表」画面に戻る</button><br>
{{ end }}
</body>
</html>
