{{define "list-last_h"}}
<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<meta http-equiv="refresh" content="{{.SecondsToReload}}; URL=list-last?eventid={{.Eventid}}&userno={{.userno}}&detail={{.Detail}}">
{{ $detail := .Detail }}
<html>
<body>
<button type="button" onclick="location.href='top'">イベント選択画面に戻る</button><br><br>
<button type="button" onclick="location.href='top?eventid={{.Eventid}}&userno={{.userno}}'">「SHOWROOMイベント結果表示」画面に戻る</button>
<br><br>
<p>直近の獲得ポイント一覧）　　<span style="color:red;">初めて使うときは表の後にある注意事項を読んでください！</span></p>
<p style="padding-left:2em">
{{.UpdateTime}}
</p>
<p style="padding-left:4em">
{{.NextTime}}
<br>
{{ if ne .NextTime "イベントは終了しています。" }}
{{.ReloadTime}}
{{ else }}
<span style="color: red;">最終結果の反映はイベント終了日翌日の13時30分頃です。<span>
{{ end }}
</p>
<table>
<tr><td align="center"><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.EventName}}</a>（{{.Eventid}}）</td></tr>
<tr><td align="center">{{.Period}}</td></tr>
</table>
{{ if eq .Detail "1" }}
<button type="button" onclick="location.href='list-last?eventid={{.Eventid}}&userno={{.userno}}&detail=0'">ルーム詳細情報を表示しない</button>
{{ else }}
<button type="button" onclick="location.href='list-last?eventid={{.Eventid}}&userno={{.userno}}&detail=1'">ルーム詳細情報を表示する</button>
{{ end }}
{{end}}
