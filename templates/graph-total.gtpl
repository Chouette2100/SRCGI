<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<button type="button" onclick="location.href='top'">イベント選択画面に戻る</button><br><br>
<button type="button" onclick="location.href='top?eventid={{.eventid}}'">「SHOWROOMイベント結果表示」画面に戻る</button><br>
  <h2>獲得ポイントグラフ</h2>
<br>
<form>
<input type="submit" value="再描画" formaction="graph-total" formmethod="POST" style="background-color: khaki">
<label>　　表示する最大ポイント　<input type="text" name="maxpoint" value="{{.maxpoint}}" size="10" required pattern="[0-9]+"><label>（表示範囲を制限しない場合は"0"とする）
<input type="hidden" name="eventid" value="{{.eventid}}">
</form>
<br><br>
<img src="{{.filename}}" alt="" width="100%">
</body>
</html>
