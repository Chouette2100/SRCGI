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
	<td>順位/<br>レベル</td>
	<td>配信者<br>(プロフィール)<br>
	<a  href="/edit-user?eventid={{.Eventid}}">一覧にないルームの追加<br>表示名・枠別貢献pt取得の設定</a></td>
	<td>獲得<br>ポイント</td>
	<td>ポイントの差</td><td style="border-right-style:none;">現配信開始</td>
	<td style="border-left-style:none;">獲得<br>ポイント</td>
	<td style="border-right-style:none;">前配信期間</td>
	<td style="border-left-style:none;">獲得<br>ポイント</td>
	<td>Next Live</td>
	<td>LIVE(配信画面)<br>FC(ファンルーム)<br>Ctn1.(貢献）<br>Graph2(グラフ)<br>Graph（グラフ旧)<br>LPS(枠別獲得pt.)
		<br>Cnt2.(枠別貢献)(<a href="/edit-user?eventid={{.Eventid}}">取得設定</a>)
		<br>※ イベント貢献履歴</td>
	{{ if and (eq $Detail "1") (ne $Isover "1") }}
	<td style="border-right-style:none;">ジャンル</td>
	<td style="border-right-style:none;">ランク</td>
	{{/*
	<td style="border-left-style:none;">next</td>
	<td style="border-left-style:none;">prev.</td>
	*/}}
	<td style="border-left-style:none;">ルーム<br>レベル</td>
	<td style="border-left-style:none;">フォロワ<br>（人）</td>
	{{/*
	<td>ファン数<br>前月</td>
	<td>ファン数<br>（人）</td>
	*/}}
	<td>配信者<br>(プロフィール)</td>
	{{ end }}
<tr>
{{ range .Scorelist }}
	<tr {{ if and ( ne $r 0 ) ( eq .Userno $r ) }} class=bgct {{ end }}>

	{{/* <td align="right">{{.Srank}}</td> */}}
       {{ if eq .Srank "9999" }}
        <td align="right">n/a</td>
        {{ else if eq .Srank "-1" }}
        <td align="right">-</td>
        {{ else if lt .Rank -1}}
        <td align="right">{{ Add .Rank 10000}}</td>
		{{ else }}
        <td align="right">{{.Srank}}</td>
        {{ end }}




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
			{{ if ne .Nperslot 0 }}
				<a href="graph-sum2?eventid={{.Eventid}}&roomid={{.Userno}}">GSum2</a>/
				<a href="graph-sum?eventid={{.Eventid}}&roomid={{.Userno}}">GSum.</a>/
				<a href="list-perslot?eventid={{.Eventid}}&roomid={{.Userno}}">LPS</a>/
			{{ else }}
				GSum2/&nbsp;GSum./&nbsp;LPS/
			{{ end }}
			{{ if ge .Ncntrb 1 }}
				<a href="list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}">Cnt2.</a>
			{{ else if ne .Spoint "0" }}
				<a href="list-cntrbex?eventid={{.Eventid}}&userno={{.Userno}}">※</a>
			{{ else }}
			    ----
			{{ end }}
		{{ end }}
	</td>
	{{ if and (eq $Detail "1") (ne $Isover "1") }}
	<td align="center" style="border-right-style:none;">{{.Roomgenre}}</td>
	<td align="center" style="border-right-style:none;">{{.Roomrank}}</td>
	{{/*
	<td align="right" style="border-left-style:none;">{{.Roomnrank}}</td>
	<td align="right" style="border-left-style:none;">{{.Roomprank}}</td>
	*/}}
	<td align="right" style="border-left-style:none;">{{.Roomlevel}}</td>
	<td align="right" style="border-left-style:none;">{{.Followers}}</td>
	{{/*
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
	*/}}
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a></td>
	{{ end }}
	</tr>
{{ end }}
</table>
{{ if gt .NoRooms .Maxrooms }}
	<form action="list-last" method="get" class="hilight">
		<input type="hidden" name="eventid" value="{{ .Eventid }}" />
		{{ if eq .Limit "TopRooms"}}
			<input type="hidden" name="limit" value="AllRooms" />
		{{ else }}
			<input type="hidden" name="limit" value="TopRooms" />
		{{ end }}
		<input type="hidden" name="detail" value="{{.Detail}}" />
		{{ if eq .Limit "TopRooms"}}
	    	<input type="submit" value="もっと見る" />
		{{ else }}
	    	<input type="submit" value="上位ルームだけ見る" />
		{{ end }}
	</form>
{{ end }}
<p style="padding-left:2em"><span style="color:red;">レベルイベントは最終レベルを達成したあとは正しい獲得ポイントは（取得できないので）表示されません。</span><br>
イベント終盤になると表示される順位も正確ではなくなります。これについては今後改善予定です。
<p style="padding-left:2em"><span style="color:red;">表に空白行があるとき</span><br>
次のような原因が考えらます。<br>
・データ取得対象となっていないルームが存在する。<br>
・イベントに新たに参加したルームがあったか、順位に大きな変動があった。<br>
　これらは「開催中イベント一覧」の「参加ルーム一覧」などから該当ルームを確認し、「表示項目選択画面」からリンクされる「イベントトップ」の「DBへのイベント参加ルーム情報の追加と更新」を実行して参加ルームのリストを更新してください。<br>
・データ取得中に獲得ポイントデータが更新され、順位や獲得ポイントに不整合が発生した。
　「イベントトップ」の「イベントパラメータの設定」でデータ取得のタイミングを少し（２，３０秒くらいうしろに）ずらします。<br>

　なおデータ取得の対象となっているルームが数十になってくるとこの現象は発生しやすくなります。データ取得の対象は上位のルームと応援しているルームくらいにしておいた方が安全です<br>
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
<tr><td>　　　</td><td>(上位30ルーム)</td><td>2. イベント終了後、獲得ポイント確定値が発表されないイベントであったとき。APIで最後に取得したポイントが表示されています。</td></tr>
<tr><td>　　　</td><td></td><td>3. イベント終了後、獲得ポイント確定値が発表されない順位であったとき。APIで最後に取得したポイントが表示されています。</td></tr>
<tr><td>　　　</td><td></td><td>4. イベント期間中の獲得ポイントが0であったとき（本来はConf.と表示すべき）</td></tr>
<tr><td>　　　</td><td>Conf.</td><td>1. イベント終了後発表された獲得ポイント確定値が表示されています。</td></tr>
<tr><td>　　　</td><td>(上位30ルーム)</td><td>2. イベント終了後「配信枠毎の獲得ポイント」の機能を利用すれば、前配信期間/獲得ポイント欄に最終枠の結果（確定値）が表示されます</td></tr>
<tr><td>　　　</td><td>31位以下のルーム</td><td>最終結果(獲得ポイント)は発表されませんので、最後のデータを取得した状態になっています。順位はイベントページから取得できるのですが今はやっていません。</td></tr>
<tr><td>　　　</td><td>レベルイベント</td><td>目標レベルを達成した後は正しい獲得ポイントを取得できません。また獲得ポイントの最終結果は発表されません。「終了イベント一覧」の「最終結果」の並び順は成績順として正しそうです。</td></tr>
</table>
<p style="padding-left:2em">※　「(DB登録済み)イベント参加ルーム一覧（確認・編集）」の一覧にないルームでも確定値が発表されたルームはここに表示されます。<br>
※　「現配信開始 獲得ポイント」がProv.あるいはConf.のときはデータ取得時刻欄にはイベント終了直後の時刻が表示されます。
</p>
</body>
</html>
{{end}}
