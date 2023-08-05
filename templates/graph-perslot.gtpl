<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
<html>

<body>
  <table>
    <tr>
      <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
      <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
      <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
      <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
    </tr>
    <tr>
      <td><button type="button" onclick="location.href='top?eventid={{.eventid}}'">イベントトップ</button></td>
      <td><button type="button" onclick="location.href='list-last?eventid={{.eventid}}'">直近の獲得ポイント</button></td>
      <td><button type="button"
          onclick="location.href='graph-total?eventid={{.eventid}}&maxpoint={{.maxpoint}}&gscale={{.gscale}}'">>獲得ポイントグラフ</button>
      </td>
      <td></td>
    </tr>
  </table>

  <h2>配信毎獲得ポイントグラフ</h2>
  <button type="button" onclick="location.href='top?eventid={{.eventid}}'">このルームの表示項目選択</button><br>
  <br><br>
  <img src="{{.filename}}" alt="" width="100%">
</body>

</html>