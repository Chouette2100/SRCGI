<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
<html>

<body>
	<span style="color:blue;">（このページもブックマーク可能です）</span>
	<br>
	<br>
	<table>
		<tr>
			<td><button type="button" onclick="location.href='top'">トップ</button>　</td>
			<td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
			<td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
			<td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
		</tr>
		<tr>
			<td><button type="button" onclick="location.href='top?eventid={{.Event_ID}}'">イベントトップ</button></td>
			<td><button type="button" onclick="location.href='list-last?eventid={{.Event_ID}}'">直近の獲得ポイント</button></td>
			<td><button type="button"
					onclick="location.href='graph-total?eventid={{.Event_ID}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
			</td>
			<td></td>
		</tr>
	</table>
	<form>
		<br>
		イベントトップ（このイベントの表示項目選択）
		<table>
			{{/*
			<tr>
				<td style="width:4em"></td>
				<td>イベントのID</td>
				<td style="width:2em"></td>
				<td>{{.Event_ID}}</td>
			</tr>
			*/}}
			<tr>
				<td style="width:4em"></td>
				<td>イベント名</td>
				<td style="width:2em"><input type="hidden" name="eventid" value="{{.Event_ID}}"><input type="hidden"
						name="userno" value="{{.Nobasis}}"></td>
				<td><a href="https://www.showroom-live.com/event/{{.Event_ID}}">{{.Event_name}}</a></td>
			</tr>
			<tr>
				<td style="width:4em"></td>
				<td>イベント期間</td>
				<td style="width:2em"></td>
				<td>{{.Period}}</td>
			</tr>
			<tr>
				<td style="width:4em"></td>
				<td>イベント参加ルーム数</td>
				<td style="width:2em"></td>
				<td>{{.NoEntry}}（最新のデータでない可能性あり）</td>
			</tr>
		</table>
		<p style="padding-left:2em">直近の獲得ポイント</p>

		<p style="padding-left:4em"><input type="submit" value="リスト" formaction="list-last" formmethod="GET"
				style="background-color: khaki">
			<input type="checkbox" name="detail" value="1">ルーム詳細情報を表示する（ルーム詳細情報とはジャンル、ランキング詳細、レベル、フォロワー数、ファン数のことです）
		</p>

		<p style="padding-left:2em">獲得ポイントの推移</p>

		<p style="padding-left:4em">
			<input type="submit" value="グラフ" formaction="graph-total" formmethod="POST" style="background-color: khaki">
			<label>　　表示する最大ポイント　
				<input type="text" name="maxpoint" value="{{.Maxpoint}}" size="10" required
					pattern="[0-9]+"><label>（表示範囲を制限しない場合は"0"とする）

					<label>　　縮尺　
						<input type="radio" name="gscale" value="100" {{ if eq .Gscale 100 }} checked {{end}}> 100%
						<input type="radio" name="gscale" value="90" {{ if eq .Gscale 90 }} checked {{ end }}> 90%
						<input type="radio" name="gscale" value="80" {{ if eq .Gscale 80 }} checked {{ end }}> 80%
						<input type="radio" name="gscale" value="70" {{ if eq .Gscale 70 }} checked {{ end }}> 70%
		</p>



		<p style="padding-left:4em"><input type="submit" value="CSV" formaction="csv-total" formmethod="POST"
				style="background-color: dimgray"></p>

		<p style="padding-left:2em">配信者間の獲得ポイントの差の推移</p>
		<p style="padding-left:4em"><input type="submit" value="グラフ" formaction="graph-dfr" formmethod="POST"
				style="background-color: dimgray"></p>
		<p style="padding-left:2em">日々の獲得ポイントと累積獲得ポイント
		<p>
		<p style="padding-left:4em"><input type="submit" value="グラフ" formaction="graph-perday" formmethod="POST"
				style="background-color: khaki">
			　　<input type="submit" value="リスト" formaction="list-perday" formmethod="POST"
				style="background-color: khaki"></p>
		<p style="padding-left:2em">配信枠毎の獲得ポイントと累積獲得ポイント
		<p>
		<p style="padding-left:4em"><input type="submit" value="グラフ" formaction="graph-perslot" formmethod="POST"
				style="background-color: khaki">
			　　<input type="submit" value="リスト" formaction="list-perslot" formmethod="POST"
				style="background-color: khaki"></p>
		<p style="padding-left:2em">(DB登録済み)イベント参加ルーム一覧（確認・編集）</p>
		<p style="padding-left:4em"><input type="submit" value="実行" formaction="edit-user" formmethod="POST"
				style="background-color: khaki"></p>
		<p style="padding-left:4em">※　獲得ポイントデータを取得するルーム、グラフに表示するルームの選択やグラフの線の色の選択に使います。</p>
		<p style="padding-left:2em">---------------------------------------------------------------</p>
		<p style="padding-left:2em">SHOWROOMイベント情報ページからDBへのイベント参加ルーム情報の追加と更新
			<br>
			<font color="red">参加者が多いイベントではルーム一覧が表示されるまで十数秒かかることがあります。<br>
				とくにイベント開始前の場合（全ルームのデータを取得してソートする必要があるため）数十秒かかることもありますので気長にお待ちください。</font>
		</p>

		<p style="padding-left:6em">
			<label>ＤＢに登録する順位の範囲
					<input type="number" name="breg" value="{{.Fromorder}}" size="3" required min="1" max="20"><label>位から
					<input type="number" name="ereg" value="{{.Toorder}}" size="3" required min="1" max="20"><label>位まで　　　
						<input type="submit" value="実行" formaction="add-event" formmethod="POST"
							style="background-color: khaki">
						<p style="padding-left:4em">イベント途中で新規の参加者がいた場合、参加者が多いイベントで大きな順位の変動があった場合などに使います。<br>
							ここで指定する順位の範囲は新規に登録する配信者の選択に関するもので、すでに登録されている配信者は範囲外であっても削除されることはありません。<br>
							また、ルームレベルやフォロワー数のように時間が経過すると変化するものを最新の状態にするのにも使います。</p>
						<p style="padding-left:2em">イベントパラメータの設定</p>
						<p style="padding-left:4em"><input type="submit" value="実行" formaction="param-event"
								formmethod="POST" style="background-color: pink"></p>
						<p style="padding-left:4em">※　イベント固有のパラメーターの設定を行います。</p>
	</form>
	<br>-----------------------------------<br>
	<table>
		<tr>
			<td>　　　</td>
			<td bgcolor="pink">　　　</td>
			<td>作成中（限定された条件下で意図した結果が得られる）</td>
		</tr>
		<tr>
			<td>　　　</td>
			<td bgcolor="khaki">　　　</td>
			<td>テスト中（いかなる場合も意図した結果が得られると思われる）</td>
		</tr>
		<tr>
			<td>　　　</td>
			<td bgcolor="aquamarine">　　　</td>
			<td>リリース（いかなる場合も意図した結果が得られる）</td>
		</tr>
	</table>

	<br>-----------------------------------<br>
	文句・ご意見・ご要望は<a href="https://twitter.com/Seppina1/">こちら</a>
</body>

</html>