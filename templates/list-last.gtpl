{{define "list-last"}}
{{ $Detail := .Detail }}
{{ $Isover := .Isover }}
{{ $r := .Roomid }}
{{/*
<table border="1" style="font-family: monospace,serif;">
*/}}
<table border="1">
{{/*
<tr><td>順位</td><td>配信者</td><td>獲得ポイント</td><td>ポイントの差</td><td colspan="2">現配信の獲得ポイント</td>
	<td>rank</td>
	<td>To the next rank</td>
	<td>level</td><td>followers</td><tr>
*/}}
<tr align="center">
	<td>順位</td>
	<td>配信者<br>(プロフィール)</td>
	<td>獲得<br>ポイント</td>
	<td>ポイントの差</td><td style="border-right-style:none;">現配信開始</td>
	<td style="border-left-style:none;">獲得<br>ポイント</td>
	<td style="border-right-style:none;">前配信期間</td>
	<td style="border-left-style:none;">獲得<br>ポイント</td>
	<td>Next Live</td>
	<td>LIVE(配信画面)<br>FC(ファンルーム)<br>Ctn1.(貢献）<br>Cnt2.(枠別貢献)</td>
	{{ if and (eq $Detail "1") (ne $Isover "1") }}
	<td style="border-right-style:none;">ジャンル</td>
	<td style="border-right-style:none;">ランク</td>
	<td style="border-left-style:none;">next</td>
	<td style="border-left-style:none;">prev.</td>
	<td>ルーム<br>レベル</td>
	<td>フォロワ<br>（人）</td>
	<td>ファン数<br>前月</td>
	<td>ファン数<br>（人）</td>
	<td>配信者<br>(プロフィール)</td>
	{{ end }}
<tr>
{{ range .Scorelist }}
	<tr {{ if and ( ne $r 0 ) ( eq .Userno $r ) }} class=bgct {{ end }}>
	<td align="right">{{.Srank}}</td>
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a></td>
	<td align="right">{{.Spoint}}</td>
	<td align="right">{{.Sdfr}}</td>
	<td align="left" style="border-right-style:none;">{{.Ptime}}</td>
	<td align="right" style="border-left-style:none;">{{.Pstatus}}</td>
	<td align="left" style="border-right-style:none;">{{.Qtime}}</td>
	<td align="right" style="border-left-style:none;">{{.Qstatus}}</td>
	<td align="center">{{.NextLive}}</td>
	<td>
		{{ if ne .Userno 0 }}
			<a href="https://www.showroom-live.com/{{.Shorturl}}">LIVE</a>/
			<a href="https://www.showroom-live.com/room/fan_club?room_id={{.Userno}}">FC</a>/
			<a href="https://www.showroom-live.com/event/contribution/{{ DelBlockID .Eventid }}?room_id={{.Userno}}">Cnt1.</a>/
			{{ if .Bcntrb }}
				<a href="list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}">Cnt2.</a>
			{{ else }}
				---
			{{ end }}
		{{ end }}
	</td>
	{{ if and (eq $Detail "1") (ne $Isover "1") }}
	<td align="center" style="border-right-style:none;">{{.Roomgenre}}</td>
	<td align="center" style="border-right-style:none;">{{.Roomrank}}</td>
	<td align="right" style="border-left-style:none;">{{.Roomnrank}}</td>
	<td align="right" style="border-left-style:none;">{{.Roomprank}}</td>
	<td align="right">{{.Roomlevel}}</td>
	<td align="right">{{.Followers}}</td>
	<td align="right">
		{{ if ne .Userno 0 }}
			{{.Fans_lst}}
		{{ end }}
	</td>
	<td align="right">
		{{ if ne .Userno 0 }}
			{{.Fans}}
		{{ end }}
	</td>
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a></td>
	{{ end }}
	</tr>
{{ end }}
</table>
<p style="padding-left:2em"><span style="color:red;">表に空白行があるときは、イベントに新たに参加したルームがあったか、順位に大きな変動があったとと思われます。</span><br>
「SHOWROOMイベント情報ページからDBへのイベント参加ルーム情報の追加と更新」を実行して参加ルームのリストを更新してください。<br>
</p>
<p style="padding-left:2em"><span style="color:red;">獲得ポイントを取得がイベント途中からはじまった場合、取得に中断がある場合「今回」、「前回」の配信期間、獲得ポイントは正しい値が表示されないことがあります。</span><br>
このような場合でも「配信枠毎の獲得ポイント」の機能では正しい値が表示されることがあります。<br>
</p>
<p style="padding-left:2em"><span style="color:red;">配信開始、前配信期間に表示される時刻には5分程度の誤差があります。</span>
</p>
<br>
<span style="font-weight: bold;">「現配信開始 獲得ポイント」欄の表示について</span><br>
・イベント終了前
<table>
<tr><td>　　　</td><td>=</td><td>配信中ではないと思われる（前回（ふつう5分前）のデータ取得から獲得ポイントは変化していない）</td></tr>
<tr><td>　　　</td><td>+##,###</td><td>配信中（前回（ふつう5分前）のデータ取得から獲得ポイントが増えている）</td></tr>
<tr><td>　　　</td><td>-##,###<br>　</td><td>減算が発生した（前回（ふつう5分前）のデータ取得から獲得ポイントが減っている）<BR>
過去の減算の有無は「配信枠毎の獲得ポイント」で確認できますが、減算が確実に確認できるのは減算の発表時（正午の数十分前）に配信が行われていない場合だけです。</td></tr>
<tr><td>　　　</td><td>n/a</td><td>獲得ポイントの増減が確認できない（データ取得プロセスの再起動直後、ルームを新しく追加したとき、データ取得に失敗したときなどに表示されます）</td></tr>
</table>
・イベント終了後
<table>
<tr><td>　　　</td><td>表示なし</td><td>2021年5月以前のイベント</td></tr>
<tr><td>　　　</td><td>Prov.</td><td>1. イベント終了後、獲得ポイント確定値が発表されるまでの期間。APIで最後に取得したポイントが表示されています。</td></tr>
<tr><td>　　　</td><td></td><td>2. イベント終了後、獲得ポイント確定値が発表されないイベントであったとき。APIで最後に取得したポイントが表示されています。</td></tr>
<tr><td>　　　</td><td></td><td>3. イベント終了後、獲得ポイント確定値が発表されない順位であったとき。APIで最後に取得したポイントが表示されています。</td></tr>
<tr><td>　　　</td><td></td><td>4. イベント期間中の獲得ポイントが0であったとき（本来はConf.と表示すべき）</td></tr>
<tr><td>　　　</td><td>Conf.</td><td>1. イベント終了後発表された獲得ポイント確定値が表示されています。</td></tr>
<tr><td>　　　</td><td></td><td>2. イベント終了後「配信枠毎の獲得ポイント」の機能を利用すれば、前配信期間/獲得ポイント欄に最終枠の結果（確定値）が表示されます</td></tr>
</table>
<p style="padding-left:2em">※　「(DB登録済み)イベント参加ルーム一覧（確認・編集）」の一覧にないルームでも確定値が発表されたルームはここに表示されます。<br>
※　「現配信開始 獲得ポイント」がProv.あるいはConf.のときはデータ取得時刻欄にはイベント終了直後の時刻が表示されます。
</p>
</body>
</html>
{{end}}
