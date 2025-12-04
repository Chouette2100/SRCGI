<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>イベント参加ルームのリスナーの貢献ポイント履歴</title>
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
            grid-column: span 1;
        }
        
        .nav-buttons button:nth-child(5) {
            grid-column: 2;
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
            min-width: 120px;
        }
        
        .param-input {
            padding: 5px 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
        }
        
        .filter-buttons {
            display: flex;
            gap: 10px;
            margin-top: 10px;
        }
        
        .filter-btn {
            padding: 8px 16px;
            border: 1px solid #007bff;
            background-color: white;
            color: #007bff;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .filter-btn:hover {
            background-color: #007bff;
            color: white;
        }
        
        .filter-btn.active {
            background-color: #007bff;
            color: white;
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
            <button type="button" onclick="location.href='eventtop?eventid={{ .EventID }}'">イベントトップ</button>
            <button type="button" onclick="location.href='list-last?eventid={{ .EventID }}'">直近の獲得ポイント</button>
            <button type="button" onclick="location.href='graph-total?eventid={{ .EventID }}'">獲得ポイントグラフ</button>
        </div>
        
        <h1>イベント参加ルームのリスナーの貢献ポイント履歴</h1>
        
        <div class="params-section">
            <div class="param-row">
                <span class="param-label">イベント:</span>
                <span><a href="/list-last?eventid={{ .EventID }}">{{ .EventName}}</a> ({{.EventID}} | {{.IeventID}})
                </span>

            </div>
            
            <form method="GET" action="listener-cntrb-history" id="paramForm">
                <input type="hidden" name="eventid" value="{{ .EventID }}">
                
                <div class="param-row">
                    <label class="param-label" for="nmonths">期間: 過去</label>
                    <input type="number" id="nmonths" name="nmonths" value="{{ .Nmonths }}" class="param-input" style="width: 80px;" min="1" max="4">
                    <span>月間</span>
                </div>
                
                <div class="param-row">
                    <label class="param-label" for="minpoint">表示する最小ポイント:</label>
                    <input type="number" id="minpoint" name="minpoint" value="{{ .Minpoint }}" class="param-input" style="width: 100px;" min="0">
                    
                    <label class="param-label" for="maxnolines" style="margin-left: 20px;">最大表示数:</label>
                    <input type="number" id="maxnolines" name="maxnolines" value="{{ .Maxnolines }}" class="param-input" style="width: 100px;" min="1">
                    　（2ヶ月設定の場合）「最小ポイント」はトップランクのルームが主体のイベで200,000〜400,000、一般には50,000くらい、5,000くらいでいいイベントも一部あります、小さく設定するととんでもなくレスポンスが悪くなるイベントがあります

                </div>
                
                <div class="filter-buttons">
                {{/*
                    {{ if eq .Ext 1 }}
                    <button type="submit" name="ext" value="1" class="filter-btn active">
                        ✓ このイベントに参加するルームに対する貢献のみ表示する
                    </button>
                    <button type="submit" name="ext" value="0" class="filter-btn">
                        このイベントに参加していないルームに対する貢献も表示する
                    </button>
                    {{ else }}
                    <button type="submit" name="ext" value="1" class="filter-btn">
                        このイベントに参加するルームに対する貢献のみ表示する
                    </button>
                    <button type="submit" name="ext" value="0" class="filter-btn active">
                        ✓ このイベントに参加していないルームに対する貢献も表示する
                    </button>
                    {{ end }}
                */}}
                    <button type="submit" name="ext" value="1" class="filter-btn active">
                        再表示
                    </button>
                </div>
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
                <th>ルーム</th>
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
                <td><a href="/list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userid}}&tlsnid={{.Lsnid}}">{{.Name }}</td>
                <td>{{ .UserName }}</td>
                <td>{{ .EventName }}</td>
                <td>{{ FormatTime .Starttime "2006-01-02 15:04" }}</td>
                <td>{{ FormatTime .Endtime "2006-01-02 15:04" }}</td>
            </tr>
            {{ end }}
        </table>
        
        {{ end }}
    </div>
</body>
</html>
