<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>アクセス集計表</title>
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
        
        h1 {
            color: #333;
            margin: 20px 0;
            font-size: 24px;
        }
        
        .type-selector {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        
        .type-buttons {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
            flex-wrap: wrap;
        }
        
        .type-buttons button {
            padding: 8px 16px;
            border: 1px solid #ddd;
            background-color: white;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .type-buttons button:hover {
            background-color: #e9ecef;
        }
        
        .type-buttons button.active {
            background-color: #28a745;
            color: white;
            border-color: #28a745;
        }
        
        .date-filter {
            display: flex;
            gap: 15px;
            align-items: center;
            flex-wrap: wrap;
        }
        
        .date-filter label {
            font-weight: 600;
        }
        
        .date-filter input[type="date"] {
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
        
        .info-message {
            padding: 10px;
            background-color: #fff3cd;
            border: 1px solid #ffc107;
            border-radius: 4px;
            margin: 10px 0;
            font-size: 13px;
        }
        
        .table-container {
            overflow-x: auto;
            margin: 20px 0;
        }
        
        .table_m {
            {{/* width: 100%; */}}
            border-collapse: collapse;
            font-size: 14px;
            background-color: white;
        }
        
        .table_c {
            width: 100%;
            border-collapse: collapse;
            font-size: 14px;
            background-color: white;
        }
        
        th, td {
            padding: 10px;
            text-align: right;
            border: 1px solid #ddd;
        }
        
        th {
            background-color: #007bff;
            color: white;
            font-weight: 600;
            position: sticky;
            top: 0;
            z-index: 10;
        }
        
        th:first-child, td:first-child {
            text-align: left;
            position: sticky;
            left: 0;
            background-color: white;
            z-index: 5;
        }
        
        th:first-child {
            background-color: #007bff;
            z-index: 15;
        }
        
        tbody tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        
        tbody tr:hover {
            background-color: #e9ecef;
        }
        
        td:first-child {
            font-weight: 500;
            max-width: 400px;
            word-wrap: break-word;
            white-space: normal;
        }
        
        .total-column {
            background-color: #fff3cd;
            font-weight: bold;
        }
        
        .encrypted-ip {
            font-family: monospace;
            font-size: 11px;
            color: #666;
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
            <button type="button" onclick="location.href='accessstatshourly'">時刻別アクセス統計</button>
            </td><td>
            <button type="button" class="active">アクセス集計表</button>
            </td>
          </tr>
        </table>
        </div>
        
        <h1>アクセス集計表</h1>
        
        <div class="type-selector">
            <div class="type-buttons">
                <button type="button" class="{{ if eq .Type "handler" }}active{{ end }}" 
                        onclick="changeType('handler')">ハンドラー別</button>
                <button type="button" class="{{ if eq .Type "ip" }}active{{ end }}" 
                        onclick="changeType('ip')">IPアドレス別</button>
                <button type="button" class="{{ if eq .Type "useragent" }}active{{ end }}" 
                        onclick="changeType('useragent')">ユーザーエージェント別</button>
            </div>
            
            <form method="POST" action="accesstable" id="filterForm">
                <input type="hidden" name="type" id="typeInput" value="{{ .Type }}">
                <div class="date-filter">
                    <label for="end_date">終了日:</label>
                    <input type="date" id="end_date" name="end_date" value="{{ .EndDate }}">
                    <button type="submit" class="reload-btn">再表示</button>
                    {{ if eq .Type "ip" }}
                    　（<label for="dkey">復号化キー:</label>
                    <input type="text" id="dkey" name="dkey" value="{{ .Dkey }}"
                        minlength="4" maxlength="64" size="10" />）
                    {{ end }}
                </div>
            </form>
        </div>
        
        {{ if eq .Type "handler" }}
        <div class="info-message">
            ハンドラー別のアクセス数（終了日: {{ .EndDate }}、過去7日間、ボット除く）
        </div>
        {{ else if eq .Type "ip" }}
        <div class="info-message">
            IPアドレス別のアクセス数（終了日: {{ .EndDate }}、過去7日間、上位50件、ボット除く）<br>
            ※IPアドレスは暗号化されています
        </div>
        {{ else if eq .Type "useragent" }}
        <div class="info-message">
            ユーザーエージェント別のアクセス数（終了日: {{ .EndDate }}、過去7日間、上位50件、ボット除く）
        </div>
        {{ end }}
        
        {{ if .Rows }}
        <div class="stats-summary">
            <div class="stats-card">
                <div class="number">{{ len .Rows }}</div>
                <div class="label">件数</div>
            </div>
            
            <div class="stats-card">
                <div class="number" id="totalAccess">-</div>
                <div class="label">総アクセス数</div>
            </div>
        </div>
        
        <div class="table-container">
            <table class="table_c">
                <thead>
                    <tr>
                        <th>
                            {{ if eq .Type "handler" }}ハンドラー
                            {{ else if eq .Type "ip" }}IPアドレス（暗号化）
                            {{ else if eq .Type "useragent" }}ユーザーエージェント
                            {{ end }}
                        </th>
                        {{ range .DateHeaders }}
                        <th>{{ . }}</th>
                        {{ end }}
                        <th class="total-column">合計</th>
                    </tr>
                </thead>
                <tbody>
                    {{ range .Rows }}
                    <tr>
                        <td {{ if eq $.Type "ip" }}class="encrypted-ip"{{ end }}>{{ .Key }}</td>
                        {{ range .DailyCounts }}
                        <td>{{ if gt .Count 0 }}{{ .Count }}{{ else }}-{{ end }}</td>
                        {{ end }}
                        <td class="total-column">{{ .Total }}</td>
                    </tr>
                    {{ end }}
                </tbody>
            </table>
        </div>
        {{ else }}
        <div class="info-message">
            指定された期間にアクセスデータが存在しません。
        </div>
        {{ end }}
    </div>
    
    {{ if .Rows }}
    <script>
        // 総アクセス数の計算
        const totalAccess = {{ range $index, $row := .Rows }}
            {{ if $index }}+{{ end }}{{ $row.Total }}
        {{ end }};
        
        document.getElementById('totalAccess').textContent = totalAccess.toLocaleString();
        
        // タイプ変更時の処理
        function changeType(type) {
            document.getElementById('typeInput').value = type;
            document.getElementById('filterForm').submit();
        }
    </script>
    {{ end }}
</body>
</html>
