<tr><td></td><td>目標ポイント</td><td><input type="text" name="target" value="0" size="2" required pattern="[0-9]+"></td></tr>
<tr><td></td><td>最大表示数</td><td><input type="text" name="maxdsp" value="25" size="2" required pattern="[0-9]+"></td></tr>
<tr><td></td><td>カラーマップ</td><td>
<select name="cmap">
<option value="0">0</option>
<option value="1">1</option>
<option value="2" selected>2</option>
</select>
</td></tr>
<tr><td></td><td><input type="{{.Submit}}" value="イベント追加" formaction="add-event?from=new-event" formmethod="POST" style="background-color: khaki"></td><td></td></tr>
</table>
</form>
<form>
<table>
<tr><td style="width:4em"></td><td><input type="submit" value="戻る" formaction="top" formmethod="POST" style="background-color: khaki"></td><td></td></tr>
</table>
</form>
<p style="padding-left:4em">
<font color="red">参加者が多いイベントでは「イベント追加」ボタンを押してから追加されてイベントのルーム一覧が表示されるまで十数秒かかることがあります。<br>
とくにイベント開始前の場合（全ルームのデータを取得してソートする必要があるため）数十秒かかることもありますので気長にお待ちください。</font></p>
</p>

<p style="padding-left:4em">
※　登録されていないイベントのIDを入力するとイベント情報の更新と参加ルーム情報の追加が行われます。
</p>
<p style="padding-left:6em">
獲得ポイントデータを取得するルームやグラフに表示するルームの設定は「イベント参加ルーム一覧（確認・編集）」から行います。
<br>ルーム数が多いと操作が煩雑になるので「登録する順位の範囲」はあんまり広くしない方がいいでしょう。あとで増やすこともできますので。
<br>"上位"とは獲得ポイントが多いという意味ですが、イベント開始前はフォロワー数で決めています。
<br>参加者が多いイベントではあとから参加したルームの情報を取得できない場合があり、そういうルームは選択の対象になりえません。
<br>イベント開始前に登録する場合は、範囲をせばめて登録し、イベント開始後しばらくしてから範囲を広くして「イベント参加者ルーム情報の更新」を行うことをおすすめします。
<br>なお指定した範囲にないルームも「イベント参加ルーム一覧（確認・編集）」で個別に追加することができます。
</p>
{{/*
<table border="1">
{{ range . }}
<tr><td>{{.Event_name}}</td>
<td>{{.Event_ID}}</td>
<td>{{.Sstart_time}}</td>
<td>{{.Send_time}}</td>
<td>{{.League_ids}}</td>
</tr>
{{ end }}
</table>
*/}}
</body>
</html>
