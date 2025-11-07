<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
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
<td><button type="button" onclick="location.href='eventtop?eventid={{.eventid}}'">イベントトップ</button></td>
<td><button type="button" onclick="location.href='list-last?eventid={{.eventid}}'">直近の獲得ポイント</button></td>
<td><button type="button" onclick="location.href='graph-total?eventid={{.eventid}}&maxpoint={{.maxpoint}}&gscale={{.gscale}}'">獲得ポイントグラフ</button></td>
<td></td>
  </tr>
</table>
  <h2>獲得ポイントグラフ</h2>
  <canvas id="myChart" width="400" height="200"></canvas>
<script>
  const ctx = document.getElementById('myChart').getContext('2d');
  const myChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
      datasets: [{
        label: 'My First dataset',
        backgroundColor: 'rgb(255, 99, 132)',
        borderColor: 'rgb(255, 99, 132)',
        data: [0, 10, 5, 2, 20, 30, 45],
      }]
    },
    options: {
      responsive: true,
      plugins: {
        zoom: {
          zoom: {
            wheel: {
              enabled: true, // マウスホイールでのズームを有効にする
            },
            pinch: {
              enabled: true, // ピンチ操作でのズームを有効にする（スマホ対応）
            },
            mode: 'xy', // x軸とy軸の両方でズームを有効にする
          },
          pan: {
            enabled: true, // パンを有効にする
            mode: 'xy', // x軸とy軸の両方でパンを有効にする
          }
        }
      }
    }
  });
  // マウスホイールイベントをカスタマイズ
const chartArea = ctx.canvas;
chartArea.addEventListener('wheel', (event) => {
    if (event.ctrlKey) {
        // Ctrlキーを押しながらのズーム（x軸のみ）
        myChart.options.plugins.zoom.zoom.mode = 'x';
    } else if (event.shiftKey) {
        // Shiftキーを押しながらのズーム（y軸のみ）
        myChart.options.plugins.zoom.zoom.mode = 'y';
    } else {
        // 修飾キーなしのズーム（両方の軸）
        myChart.options.plugins.zoom.zoom.mode = 'xy';
    }

    // ズームを適用
    myChart.update('none'); // グラフを再描画
});
</script>
</body>
</html>