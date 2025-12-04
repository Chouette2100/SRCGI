<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>ルーム別リスナー貢献ポイントランキング</title>
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
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 10px;
            margin-bottom: 20px;
        }
        
        .nav-buttons button {
            padding: 8px 16px;
            border: 1px solid #ddd;
            background-color: #f8f9fa;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            white-space: nowrap;
        }
        
        .nav-buttons button:hover {
            background-color: #e9ecef;
        }
        
        h1 {
            color: #333;
            margin: 20px 0;
            font-size: 24px;
        }
        
        .user-info {
            background-color: #f8f9fa;
            padding: 10px 15px;
            border-radius: 6px;
            margin-bottom: 20px;
            font-size: 16px;
            font-weight: 600;
        }
        
        .params-section {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        
        .param-row {
            margin-bottom: 10px;
            display: flex;
            align-items: center;
            gap: 10px;
            flex-wrap: wrap;
        }
        
        .param-label {
            font-weight: 600;
            min-width: 150px;
        }
        
        .param-input {
            padding: 5px 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
        }
        
        .submit-btn {
            padding: 8px 24px;
            border: 1px solid #007bff;
            background-color: #007bff;
            color: white;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            margin-top: 10px;
        }
        
        .submit-btn:hover {
            background-color: #0056b3;
        }
        
        .data-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
            font-size: 14px;
        }
        
        .data-table th {
            background-color: gainsboro;
            text-align: center;
            padding: 10px 5px;
            border: 1px solid #ddd;
            font-weight: 600;
        }
        
        .data-table td {
            padding: 8px 5px;
            border: 1px solid #ddd;
        }
        
        .data-table tr:nth-child(odd) {
            background-color: gainsboro;
        }
        
        .data-table td:first-child {
            text-align: right;
        }
        
        .data-table a {
            color: #007bff;
            text-decoration: none;
        }
        
        .data-table a:hover {
            text-decoration: underline;
        }
        
        .info-message {
            padding: 10px;
            background-color: #d1ecf1;
            border: 1px solid #bee5eb;
            border-radius: 4px;
            margin: 10px 0;
            font-size: 13px;
            color: #0c5460;
        }
        
        .error-message {
            padding: 10px;
            background-color: #f8d7da;
            border: 1px solid #dc3545;
            border-radius: 4px;
            margin: 10px 0;
            color: #721c24;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="nav-buttons">
            <button type="button" onclick="location.href='top'">トップ</button>
            <button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button>
            <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button>
            <button type="button" onclick="location.href='closedevents'">終了イベント一覧</button>
        </div>
        
        <h1>ルーム別リスナー貢献ポイントランキング</h1>
        
        <div class="user-info">
            {{ .UserName }} ( {{ .Userid }} )
        </div>
        
        <div class="params-section">
            <form method="GET" action="room-cntrb-history" id="paramForm">
                <input type="hidden" name="userid" value="{{ .Userid }}">
                
                <div class="param-row">
                    <label class="param-label" for="nmonths">期間: 過去</label>
                    <input type="number" id="nmonths" name="nmonths" value="{{ .Nmonths }}" class="param-input" style="width: 80px;" min="1" max="12">
                    <span>月間</span>
                </div>
                
                <div class="param-row">
                    <label class="param-label" for="minpoint">表示する最小ポイント:</label>
                    <input type="number" id="minpoint" name="minpoint" value="{{ .Minpoint }}" class="param-input" style="width: 120px;" min="0">
                    
                    <label class="param-label" for="maxnolines" style="margin-left: 20px;">最大表示数:</label>
                    <input type="number" id="maxnolines" name="maxnolines" value="{{ .Maxnolines }}" class="param-input" style="width: 100px;" min="1">
                </div>
                
                <button type="submit" class="submit-btn">再表示</button>
            </form>
        </div>
        
        {{ if .ErrMsg }}
        <div class="error-message">{{ .ErrMsg }}</div>
        {{ else }}
        
        <div class="info-message">
            表示件数: {{ len .DataList }} 件
        </div>
        
        <table class="data-table">
            <tr>
                <th>貢献ポイント</th>
                <th>リスナー</th>
                <th>イベント</th>
                <th>開始日時</th>
                <th>終了日時</th>
            </tr>
            
            {{ range .DataList }}
            {{ if iscurrent .Endtime }}
            <tr style="background-color: yellow;">
            {{ else }}
            <tr>
            {{ end }}
                <td>{{ Comma .Point }}</td>
                <td>{{ .Name }}</td>
                <td><a href="/list-last?eventid={{ .EventID }}">{{ .EventName}}</a> ({{.EventID}} | {{.IeventID}})</td>
                <td>{{ FormatTime .Starttime "2006-01-02 15:04" }}</td>
                <td>{{ FormatTime .Endtime "2006-01-02 15:04" }}</td>
            </tr>
            {{ end }}
        </table>
        
        {{ end }}
    </div>
</body>
</html>
