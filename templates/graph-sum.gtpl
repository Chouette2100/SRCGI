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
                                <p>{{.Eventname}}({{.Roomname}})のグラフを表示するには、セキュリティチ
ェックを完了してくだ>さい。</p>
                                <p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
                                <form method="POST" action="graph-sum">
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
    イベント最終日翌日のお昼までは最終データは仮のものです。これはたいていの場合実際のポイントよりかなり小さいです</p>
    <p>枚枠の具体的なポイントはバー（棒）の上にマウスカーソルをかざすと表示されます<br>
    このとき表示される日時は配信のおおむね開始日時になります（終了日時じゃありません！）</p>
    <div style="width: 80%;">
        <canvas id="myChart" width="800" height="400"></canvas>
    </div>

    {{ end }}
    
    <script>
        Chart.register(ChartZoom);
        fetch('/graph-sum-data?eventid={{.Eventid}}&roomid={{.Roomid}}')
            .then(response => response.json())
            .then(data => {
                const ctx = document.getElementById('myChart').getContext('2d');
                const myChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: data.dtime,
                        datasets: [{
                            label: '累積ポイント',
                            data: data.data1,
                            borderColor: 'rgba(0, 0, 255, 1)',
                            backgroundColor: 'rgba(0, 0, 255, 1)',
                            barPercentage: 0.5,
                            fill: false,
                        },
                        {
                            type: "bar",
                            label: '獲得ポイント',
                            data: data.data2,
                            borderColor: 'rgba(75, 192, 192, 1)',
                            backgroundColor: 'rgba(100, 200, 0, 1)',
                            barPercentage: 0.5,
                            fill: false,
                        }]
                    },
                    options: {
                        responsive: true,
                        interaction: {
                            mode: 'nearest',
                            intersect: false,
                        },
                        scales: {
                            x: {
                                type: 'time', // x軸を時間軸として設定
                                time: {
                                    unit: 'day', // 時間の単位（日単位）
                                    displayFormats: {
                                        // day: 'yyyy-MM-dd' // 表示フォーマット
                                        day: 'MM-dd' // 表示フォーマット
                                    },
                                    tooltipFormat: 'yyyy-MM-dd HH-mm', // ツールチップのフォーマット
                                },
                                {{/*
                                time: {
                                    unit: 'day', // 時間の単位（日単位）
                                    displayFormats: {
                                        day: 'YYYY-MM-DD' // 表示フォーマット
                                    }
                                },
                                */}}
                            },
                            y: {
                                type: 'linear',
                            }
                        },
                        /*
                        plugins: {
                            zoom: {
                                zoom: {
                                    wheel: {
                                        enabled: true,
                                    },
                                    pinch: {
                                        enabled: true,
                                    },
                                    mode: 'xy',
                                },
                                pan: {
                                    enabled: true,
                                    mode: 'xy',
                                }
                            }
                        }
                        */
                    }
                });
                // // デバッグ用: パン設定をコンソールに出力
                // console.log('Pan enabled:', myChart.options.plugins.zoom.pan.enabled);
                // myChart.options.plugins.zoom.pan.onPan = function ({ chart }) {
                //     console.log('Panning!', chart.scales.x.min, chart.scales.x.max);
                // };
                // console.log(myChart.options.plugins.zoom.pan);
            });

    </script>
</body>

</html>
