<!DOCTYPE html>
<html>

<head>
    {{/* Turnstile 1 */}}
    {{if .TurnstileSiteKey}}
        <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    {{end}}
    {{/* ----------- */}}
    {{/*
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-plugin-zoom"></script>
    */}}
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src='https://cdnjs.cloudflare.com/ajax/libs/hammer.js/2.0.8/hammer.js'></script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-plugin-zoom@latest"></script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-date-fns"></script>
</head>

<body>
       <label>
            <input type="checkbox" id="zoomCheckbox" checked> Zoom/Panを有効にする
        </label>
    <table>
        <tr>
            <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
            <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
        </tr>
        <tr>
            <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
            <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
            <td><button type="button"
                    onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
            </td>
            <td></td>
        </tr>
    </table>

        {{if .TurnstileSiteKey }}
                <!-- Turnstileチャレンジ表示 -->
                        <div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9;">
                                <h3>セキュリティチェック</h3>
                                {{if .TurnstileError}}
                                <p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
                                {{end}}
                                <p>{{.Eventname}}({{.Roomname}})のグラフを表示するには、セキュリティチェックを完了してくだ>さい。</p>
                                <p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
                                <form method="POST" action="graph-sum2">
                                        <input type="hidden" name="eventid" value="{{.Eventid}}">
                                        <input type="hidden" name="roomid" value="{{.Roomid}}">

                                        <input type="hidden" name="requestid" value="{{.RequestID}}">

                                        <div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
                                        <br>
                                        <button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認>して続行</button>
                                </form>
                        </div>
                        <br><br>
    {{else}}

    <br>
    <h3>{{ .Title }}</h3>
    {{ .Eventname }}（{{.Eventid}}:{{.Ieventid}} )<br>
    {{ .Period}}<br>
    {{ .Roomname}}（{{.Roomid}}）<br>
    <p>あくまで概算です。欠測がある場合などは正しく算出できません。<br>
    イベント最終日翌日のお昼までは最終データは仮のものです。たいていの場合実際のポイントよりかなり小さいです</p>
    <p>枚枠の具体的なポイントはバー（棒）の上にマウスカーソルをかざすと表示されます<br>
    このとき表示される日時はおおむね配信開始日時になります（終了日時じゃありません！）</p>
    <p>数値データ　
    <a href="/list-perslot?eventid={{.Eventid}}&roomid={{.Roomid}}">枠別獲得ポイント(個別・表)</a>　
    <a href="/list-perslot?eventid={{.Eventid}}">枠別獲得ポイント(全体・表)</a>
    <a href="/graph-sum-data1?eventid={{.Eventid}}&roomid={{.Roomid}}">獲得ポイント(個別)・JSON</a>　
    <a href="/graph-sum-data2?eventid={{.Eventid}}&roomid={{.Roomid}}">枠別獲得ポイント(個別)・JSON</a>　
    </p>
    <div style="width: 80%;">
        <canvas id="myChart" width="800" height="400"></canvas>
    </div>
    <script>
    // script.js

// データを取得する関数
async function fetchData(url) {
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`データの取得に失敗しました: ${response.statusText}`);
    }
    return response.json();
}

// グラフを描画する関数
function renderChart(dataSet1, dataSet2) {
    const ctx = document.getElementById('myChart').getContext('2d');
    const myChart = new Chart(ctx, {
        type: 'line', // 折れ線グラフ
        data: {
            datasets: [
                {
                    label: '累積ポイント',
                    data: dataSet1.map(item => ({
                        x: item.timestamp, // 時刻
                        y: item.value      // 観測値
                    })),
                    borderColor: 'rgba(75, 192, 192, 1)',
                    fill: false
                },
                {
                    type: "bar",
                    label: '獲得ポイント',
                    data: dataSet2.map(item => ({
                        x: item.timestamp, // 時刻
                        y: item.value      // 観測値
                    })),
                    backgroundColor: 'rgba(255, 99, 132, 1)',
                    barPercentage: 0.3,
                    fill: false
                }
            ]
        },
        options: {
            scales: {
                x: {
                    type: 'time', // x軸を時間軸として設定
                    type: 'time', // x軸を時間軸として設定
                    time: {
                        unit: 'day', // 時間の単位（日単位）
                        displayFormats: {
                            // day: 'yyyy-MM-dd' // 表示フォーマット
                            day: 'MM-dd' // 表示フォーマット
                        },
                        tooltipFormat: 'yyyy-MM-dd HH-mm', // ツールチップのフォーマット
                    }
                },
                y: {
                    beginAtZero: true // y軸を0から始める
                }
            }
        }
    });
}

// メイン処理
async function main() {
    try {
        // データを取得
        const dataSet1 = await fetchData('/graph-sum-data1?eventid={{.Eventid}}&roomid={{.Roomid}}');
        const dataSet2 = await fetchData('/graph-sum-data2?eventid={{.Eventid}}&roomid={{.Roomid}}');

        // グラフを描画
        renderChart(dataSet1, dataSet2);
    } catch (error) {
        console.error('エラーが発生しました:', error);
    }
}

// 実行
main();
    </script>
{{ end }}
</body>

</html>