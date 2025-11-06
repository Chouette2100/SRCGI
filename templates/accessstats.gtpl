<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>アクセス統計</title>
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
        
        .date-filter-section {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        
        .date-inputs {
            display: flex;
            gap: 15px;
            align-items: center;
            flex-wrap: wrap;
        }
        
        .date-inputs label {
            font-weight: 600;
        }
        
        .date-inputs input[type="date"] {
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
            height: 400px;
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
            <button type="button" class="active">日別アクセス統計</button>
            </td><td>
            <button type="button" onclick="location.href='accessstatshourly'">時刻別アクセス統計</button>
            </td><td>
            <button type="button" onclick="location.href='accessstable'">アクセス集計表<button>
            </td>
          </tr>
        </table>
        </div>
        
        <h1>日別アクセス統計</h1>
        
        <div class="date-filter-section">
            <form method="GET" action="accessstats">
                <div class="date-inputs">
                    {{/*
                    <label for="start_date">開始日:</label>
                    <input type="date" id="start_date" name="start_date" value="{{ .StartDate }}">
                    */}}
                    <label for="start_date">期間:</label>
                    <input type="number" id="period" name="period" min="7" max="200" value="{{ if eq .Period 0 }}31{{ else }}{{ .Period }}{{ end }}">日　
                    
                    <label for="end_date">終了日:</label>
                    <input type="date" id="end_date" name="end_date" value="{{ .EndDate }}">
                    
                    <button type="submit" class="reload-btn">再表示</button>
                </div>
            </form>
        </div>
        
        <div class="info-message">
            表示期間: {{ .StartDate }} 〜 {{ .EndDate }} ({{ len .Stats }}日間のデータ)
        </div>
        
        {{ if .Stats }}
        <div class="stats-summary">
            <div class="stats-card">
                <div class="number" id="totalAccess">-</div>
                <div class="label">総アクセス数</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="avgAccess">-</div>
                <div class="label">1日平均</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="maxAccess">-</div>
                <div class="label">最大アクセス数</div>
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
                {{ if $index }},{{ end }}'{{ $stat.AccessDate }}'
            {{ end }}
        ];
        
        // 日付を短縮形式にフォーマット
        const labels = rawLabels.map(dateStr => {
            const date = new Date(dateStr);
            const year = date.getFullYear().toString().slice(-2); // 下2桁
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            return `${year}-${month}-${day}`;
        });
        
        const data = [
            {{ range $index, $stat := .Stats }}
                {{ if $index }},{{ end }}{{ $stat.AccessCount }}
            {{ end }}
        ];
        
        // 統計計算
        const total = data.reduce((sum, count) => sum + count, 0);
        const max = Math.max(...data);
        const avg = Math.round(total / data.length);
        
        // 統計値を画面に表示
        document.getElementById('totalAccess').textContent = total.toLocaleString();
        document.getElementById('avgAccess').textContent = avg.toLocaleString();
        document.getElementById('maxAccess').textContent = max.toLocaleString();
        
        // Chart.jsの設定
        const ctx = document.getElementById('accessChart').getContext('2d');
        const accessChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: '日別アクセス数',
                    data: data,
                    borderColor: '#007bff',
                    backgroundColor: 'rgba(0, 123, 255, 0.1)',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.2,
                    pointBackgroundColor: '#007bff',
                    pointBorderColor: '#fff',
                    pointBorderWidth: 2,
                    pointRadius: 4,
                    pointHoverRadius: 6
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'アクセス数推移グラフ'
                    },
                    legend: {
                        display: true,
                        position: 'top'
                    }
                },
                scales: {
                    x: {
                        display: true,
                        title: {
                            display: true,
                            text: '日付'
                        },
                        ticks: {
                            maxTicksLimit: 15
                        }
                    },
                    y: {
                        display: true,
                        title: {
                            display: true,
                            text: 'アクセス数'
                        },
                        beginAtZero: true,
                        ticks: {
                            stepSize: 5000
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