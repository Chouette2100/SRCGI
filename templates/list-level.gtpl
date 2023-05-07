<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<p>
{{/*
<button type="button" onclick="location.href='top'">サイトトップ画面に戻る</button><br>
*/}}
<button type="button" onclick="location.href='top?userno={{.Userno}}'">イベント選択画面に戻る</button>
</p>
<h2>ルーム情報（レベルとフォロワーの推移）</h2>
<p style="padding-left:2em">
この機能は現在作成中です<br>
下表の「日時」はデータを取得した日時です。レベルが上がったり、フォロワー数が変化した時刻ではありません。</p>

<p style="padding-left:2em">
<a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{ .User_name }}</a>（{{.Userno}}）<br>
</p>
<table border="1">
<tr align="center">
    <td>日時</td>
    <td>レベル</td>
    <td>フォロワー数</td>
    <td>ランク</td>
    <td>To the next</td>
    <td>From the prev.</td>
    <td>ジャンル</td>
    <td>ルーム名</td>
 <tr>
{{ range .RoomLevelList }}
<tr align="right">
    <td>{{.Sts}}</td>
    <td>{{.Level}}</td>
    <td>{{.Followers}}</td>
    <td>{{.Rank}}</td>
    <td>{{.Nrank}}</td>
    <td>{{.Prank}}</td>
    <td>{{.Genre}}</td>
    <td align="left">{{.User_name}}</td>
 <tr>
{{ end }}
</table>
</body>
</html>
