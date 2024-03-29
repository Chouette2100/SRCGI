{{/*}}
<p>2. イベント選択</p>
<form>
<p style="padding-left:2em"><dive title="下の方にあるイベント管理で登録されたイベントがここに表示されます">表示・操作の対象となるイベントを選択してください。<div></p>
	<p style="padding-left:2em">
<select name="eventid">
		{{ range . }}
			<option value="{{.EventID}}" {{.Selected}}>{{.EventName}}　({{.EventID}})</option>
		{{ end }}
	</select>
<span style="padding-left:2em"><input type="submit" value="決定" formaction="top" formmethod="GET"></span></p>

</form>
========================================================================
<br>
<p><a name="newevent" id="newevent">イベント管理</a></p>
{{*/}}
{{/*
<a href="add-event" style="background-color: dimgray">SHOWROOMイベント情報ページからDBへの新規イベントの追加</a>　データ取得の対象としたいイベントの追加<br>
<button type="button" onclick="location.href='add-event'" style="background-color: dimgray">新規イベントの追加</button>　データ取得の対象としたいイベントの追加
*/}}
{{/*}}
<p style="padding-left:2em;color:blue;">
イベントの獲得ポイントデータの推移をみたいときは「新規イベントと参加ルームの登録」でイベントの登録を行ってください。<br>
登録が終わると上のリストからイベントが選択できるようになります。<br>
登録後、「(DB登録済み)イベント参加ルーム一覧（確認・編集）」で一覧に自分がデータが見たいルームがあることを確認してください。<br>
もしないときは一覧の下にある「一覧にないルームの追加」で自分がデータを見たいルームを追加してください。<br>
イベントの登録はイベント開始前でも開始後でも可能ですが、開始前の登録をおすすめします。<br>
ただし、イベントの参加者が多いイベントではイベントが始まってから登録した方がいいでしょう。<br>
<br>
リストで「データ取得」にチェックを入れたルームの）獲得ポイントデータの取得はイベント登録後5分以内に始まります。<br>
ですから獲得ポイントデータが表示されるのは登録してから5分程度経ってからです。
</p>
<p style="padding-left:2em;color:red;">ブロックイベントの場合は次のような形式でイベントを指定してください『circle2023_2nd_a?block_id=7101』</p>
<p style="padding-left:2em">
新規イベントと参加ルームの登録
{{*/}}
{{/* <br><font color="red">参加者が多いイベントではルーム一覧が表示されるまで十数秒かかることがあります<br>
とくにイベント開始前の場合数十秒かかることもありますので気長にお待ちください。</font></p> */}}
{{/*}}
<p>
<form>
<table>
<tr><td style="width:4em"></td><td>イベントのID</td><td><input type="text" name="eventid" required="true">（イベントページのURLの最後のフィールド）
</td></tr>
<tr><td style="width:4em"></td><td></td><td align="right"><input type="submit" value="「新規イベントと参加ルームの登録」画面へ" formaction="new-event" formmethod="POST" ></td></tr>
</table>
</p>
</form>
<p style="padding-left:4em">
登録されていないイベントのIDを入力してイベント情報と参加ルーム情報の追加を行います。IDというのはイベントページのURLの最後のフィールドのことです。<br>
ここはプルダウンかラジオボタンでイベントのリストから選択するようにすれば使いやすいのですが、APIで開始前のイベントのリストが得られない、イベント一覧ページは作りが動的なのでサーバー上でスクレイピングができないというような理由があってそうできないのです。かならずクレームがくると思うのであらかじめお断りしておきます。<p>
<p style="padding-left:2em">共通パラメータの設定</p>
<form>
<input type="hidden" name="eventid" value="global" > 
<p style="padding-left:4em"><input type="submit" value="実行" formaction="param-global" formmethod="POST" style="background-color: dimgray"></p>
<p style="padding-left:4em">※　アプリケーション全体に関するパラメーターの設定を行います。</p>
</p>
{{*/}}
</body>
</html>
