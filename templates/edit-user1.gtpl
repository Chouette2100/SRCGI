<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<button type="button" onclick="location.href='top?eventid={{.Eventid}}'">「SHOWROOMイベント結果表示」画面に戻る</button>
<br><br>
<p>(DB登録済)イベント参加者ルーム一覧（確認・編集）　　<span style="color:red;">初めて使うときは表の後にある説明をよく読んでください！</span></p>
<table style="text-align: center">
<tr><td style="width:2em"></td><td><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a></td></tr>
<tr><td style="width:2em"></td><td>{{.Period}}</td></tr>
</table>
