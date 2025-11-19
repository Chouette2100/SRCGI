<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>時刻単位アクセス統計</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        
        .container {
            max-width: 1800px;
            margin: 0 auto;
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        
        .nav-buttons {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .nav-buttons button {
            padding: 8px 16px;
            border: 1px solid #ddd;
            background-color: #f8f9fa;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .nav-buttons button:hover {
            background-color: #e9ecef;
        }
        
        .nav-buttons button.active {
            background-color: #007bff;
            color: white;
            border-color: #007bff;
        }

                #period {
            width: 3em;
        }
        
        h1 {
            color: #333;
            margin: 20px 0;
            font-size: 24px;
        }
        
        .datetime-filter-section {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        
        .datetime-inputs {
            display: flex;
            gap: 15px;
            align-items: center;
            flex-wrap: wrap;
        }
        
        .datetime-inputs label {
            font-weight: 600;
        }
        
        .datetime-inputs input[type="datetime-local"] {
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
        }
        
        .reload-btn {
            padding: 8px 16px;
            border: 1px solid #28a745;
            background-color: #28a745;
            color: white;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .reload-btn:hover {
            background-color: #218838;
        }
        
        .chart-container {
            position: relative;
            height: 500px;
            margin: 20px 0;
        }
        
        .info-message {
            padding: 10px;
            background-color: #fff3cd;
            border: 1px solid #ffc107;
            border-radius: 4px;
            margin: 10px 0;
            font-size: 13px;
        }

        .table_m {
            {{/* width: 100%; */}}
            border-collapse: collapse;
            font-size: 14px;
            background-color: white;
        }
        
        .stats-summary {
            display: flex;
            gap: 20px;
            margin: 20px 0;
            flex-wrap: wrap;
        }
        
        .stats-card {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            text-align: center;
            min-width: 120px;
        }
        
        .stats-card .number {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
        }
        
        .stats-card .label {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
    </style>
</head>
<body>
    {{ $tn := .TimeNow }}
    
    <div class="container">
        <div class="nav-buttons">
        <table class="table_m">
          <tr>
            <td>
            <button type="button" onclick="location.href='top'">トップ</button>
            </td><td>
            <button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button>
            </td><td>
            <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button>
            </td><td>
            <button type="button" onclick="location.href='closedevents'">終了イベント一覧</button>
            </td>
          </tr><tr>
            <td>
            </td><td>
            <button type="button" onclick="location.href='accessstats'">日別アクセス統計</button>
            </td><td>
            <button type="button" class="active">時刻別アクセス統計</button>
            </td><td>
            <button type="button" onclick="location.href='accesstable'">アクセス集計表</button>
            </td>
          </tr>
        </table>
        </div>
        
        <h1>時刻単位アクセス統計</h1>
        
        <div class="datetime-filter-section">
            <form method="GET" action="accessstatshourly">
                <div class="datetime-inputs">
                {{/*
                    <label for="start_datetime">開始日時:</label>
                    <input type="datetime-local" id="start_datetime" name="start_datetime" value="{{ .StartDateTime }}">
                */}}

                    <label for="period">期間:</label>
                    <input type="number" id="period" name="period"
                    min="24" max="216"
                    value="{{ if eq .Period 0 }}72{{ else }}{{ .Period }}{{ end }}">(時間)　
                    <label for="end_datetime">終了日時:</label>
                    <input type="datetime-local" id="end_datetime" name="end_datetime" value="{{ .EndDateTime }}">
                    
                    <button type="submit" class="reload-btn">再表示</button>
                </div>
            </form>
        </div>
        
        <div class="info-message">
            表示期間: {{ .StartDateTime }} 〜 {{ .EndDateTime }} ({{ len .Stats }}時間分のデータ)
        </div>
        
        {{ if .Stats }}
        <div class="stats-summary">
            <div class="stats-card">
                <div class="number" id="totalAccess">-</div>
                <div class="label">総アクセス数</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="avgAccess">-</div>
                <div class="label">1時間平均</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="maxAccess">-</div>
                <div class="label">最大アクセス数</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="totalLegitimate">-</div>
                <div class="label">正規リクエスト</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="totalTurnstileFail">-</div>
                <div class="label">Turnstile失敗</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="totalBot">-</div>
                <div class="label">ボット</div>
            </div>
        </div>
        
        <div class="chart-container">
            <canvas id="accessChart"></canvas>
        </div>
        {{ else }}
        <div class="info-message">
            指定された期間にアクセスデータが存在しません。
        </div>
        {{ end }}
    </div>
    
    {{ if .Stats }}
    <script>
        // グラフデータの準備
        const rawLabels = [
            {{ range $index, $stat := .Stats }}
                {{ if $index }},{{ end }}'{{ $stat.AccessHour }}'
            {{ end }}
        ];
        
        // 日時を短縮形式にフォーマット (25-11-05 16:00 形式)
        const labels = rawLabels.map(dateTimeStr => {
            const dt = new Date(dateTimeStr);
            const year = dt.getFullYear().toString().slice(-2); // 下2桁
            const month = String(dt.getMonth() + 1).padStart(2, '0');
            const day = String(dt.getDate()).padStart(2, '0');
            const hour = String(dt.getHours()).padStart(2, '0');
            return `${year}-${month}-${day} ${hour}:00`;
        });
        
        const data = [
            {{ range $index, $stat := .Stats }}
                {{ if $index }},{{ end }}{{ $stat.AccessCount }}
            {{ end }}
        ];
        
        const legitimateData = [
            {{ range $index, $stat := .Stats }}
                {{ if $index }},{{ end }}{{ $stat.LegitimateCount }}
            {{ end }}
        ];
        
        const turnstileFailData = [
            {{ range $index, $stat := .Stats }}
                {{ if $index }},{{ end }}{{ $stat.TurnstileFailCount }}
            {{ end }}
        ];
        
        const botData = [
            {{ range $index, $stat := .Stats }}
                {{ if $index }},{{ end }}{{ $stat.BotCount }}
            {{ end }}
        ];
        
        // 統計計算
        const total = data.reduce((sum, count) => sum + count, 0);
        const max = Math.max(...data);
        const avg = Math.round(total / data.length);
        
        const totalLegitimate = legitimateData.reduce((sum, count) => sum + count, 0);
        const totalTurnstileFail = turnstileFailData.reduce((sum, count) => sum + count, 0);
        const totalBot = botData.reduce((sum, count) => sum + count, 0);
        
        // 統計値を画面に表示
        document.getElementById('totalAccess').textContent = total.toLocaleString();
        document.getElementById('avgAccess').textContent = avg.toLocaleString();
        document.getElementById('maxAccess').textContent = max.toLocaleString();
        document.getElementById('totalLegitimate').textContent = totalLegitimate.toLocaleString();
        document.getElementById('totalTurnstileFail').textContent = totalTurnstileFail.toLocaleString();
        document.getElementById('totalBot').textContent = totalBot.toLocaleString();
        
        // Y軸の最大値を計算（500刻みで切り上げ）
        const maxYValue = Math.ceil(max / 500) * 500;
        
        // Chart.jsの設定
        const ctx = document.getElementById('accessChart').getContext('2d');
        const accessChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: '正規リクエスト',
                        data: legitimateData,
                        borderColor: '#28a745',
                        backgroundColor: 'rgba(40, 167, 69, 0.1)',
                        borderWidth: 2,
                        fill: true,
                        tension: 0.2,
                        pointBackgroundColor: '#28a745',
                        pointBorderColor: '#fff',
                        pointBorderWidth: 2,
                        pointRadius: 3,
                        pointHoverRadius: 5
                    },
                    {
                        label: 'Turnstile失敗',
                        data: turnstileFailData,
                        borderColor: '#ffc107',
                        backgroundColor: 'rgba(255, 193, 7, 0.1)',
                        borderWidth: 2,
                        fill: true,
                        tension: 0.2,
                        pointBackgroundColor: '#ffc107',
                        pointBorderColor: '#fff',
                        pointBorderWidth: 2,
                        pointRadius: 3,
                        pointHoverRadius: 5
                    },
                    {
                        label: 'ボット',
                        data: botData,
                        borderColor: '#dc3545',
                        backgroundColor: 'rgba(220, 53, 69, 0.1)',
                        borderWidth: 2,
                        fill: true,
                        tension: 0.2,
                        pointBackgroundColor: '#dc3545',
                        pointBorderColor: '#fff',
                        pointBorderWidth: 2,
                        pointRadius: 3,
                        pointHoverRadius: 5
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: '時刻別アクセス数推移グラフ',
                        font: {
                            size: 16
                        }
                    },
                    legend: {
                        display: true,
                        position: 'top'
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return 'アクセス数: ' + context.parsed.y.toLocaleString();
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        display: true,
                        title: {
                            display: true,
                            text: '日時'
                        },
                        ticks: {
                            maxRotation: 45,
                            minRotation: 45,
                            maxTicksLimit: 24
                        }
                    },
                    y: {
                        display: true,
                        title: {
                            display: true,
                            text: 'アクセス数'
                        },
                        beginAtZero: true,
                        max: maxYValue > 0 ? maxYValue : undefined,
                        ticks: {
                            stepSize: 100,
                            callback: function(value) {
                                return value.toLocaleString();
                            }
                        }
                    }
                },
                interaction: {
                    intersect: false,
                    mode: 'index'
                },
                hover: {
                    mode: 'nearest',
                    intersect: false
                }
            }
        });
    </script>
    {{ end }}
</body>
</html>
