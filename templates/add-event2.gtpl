<br> 
<table border="1">
	<tr align="center">
	<td>ルーム名</td>
	<td>プロフィール</td>
	<td>ジャンル</td>
	<td align="center" style="border-right-style:none;">ランク</td>
	<td align="right" style="border-left-style:none;"></td>
	<td>レベル</td>
	<td>フォロワー数</td>
	<td>ポイント</td>
	<td>状態</td>
<tr>
{{ range . }}
<tr style="color:{{.Statuscolor}}">
	<td><a href="https://www.showroom-live.com/{{.Account}}">{{.Name}}</a></td>
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.ID}}">プロフィール</a></td>
	<td>{{.Genre}}</td>
	<td align="center" style="border-right-style:none;">{{.Rank}}</td>
	<td align="right" style="border-left-style:none;">{{.Nrank}}</td>
	<td align="right">{{.Level}}</td>
	<td align="right">{{.Followers}}</td>
	<td align="right">{{.Spoint}}</td>
	<td>{{.Status}}</td>
</tr>
{{ end }}
</table>
<br>
ここで状態が「新規」（緑色で表示）となっているルームについては「(DB登録済み)イベント参加ルーム一覧（確認・編集）」で表・グラフに表示するルーム名、データ取得・表示の有無、グラフ表示色を設定してください。
</body>
</html>
