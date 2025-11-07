<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>開催中イベント一覧</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        
        .container {
            {{/* 画面の幅は有効に使いたい
            max-width: 1400px;
            */}}
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
        
        .filter-section {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        
        .data-filter-btn {
            padding: 8px 16px;
            border: 1px solid #28a745;
            background-color: white;
            color: #28a745;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .data-filter-btn:hover {
            background-color: #28a745;
            color: white;
        }
        
        .data-filter-btn.active {
            background-color: #28a745;
            color: white;
        }
        
        /* テーブルスタイル */
        .events-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
            font-size: 14px;
        }
        
        .events-table th {
            background-color: gainsboro;
            text-align: center;
            padding: 10px 5px;
            border: 1px solid #ddd;
            font-weight: 600;
        }
        
        .events-table td {
            padding: 8px 5px;
            border: 1px solid #ddd;
        }
        
        {{/*
        .events-table tr:nth-child(even) {
            background-color: gainsboro;
        }
        */}}
        
        {{/* 色が明るすぎるのも問題だが、参照済みでも色が変わらないのがもっと問題
        .events-table a {
            color: #007bff;
            text-decoration: none;
        }
        */}}
        
        .events-table a:hover {
            text-decoration: underline;
        }
        
        .info-message {
            padding: 10px;
            background-color: #fff3cd;
            border: 1px solid #ffc107;
            border-radius: 4px;
            margin: 10px 0;
            font-size: 13px;
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
    {{ $tn := .TimeNow }}
    
    <div class="container">
        <div class="nav-buttons">
            <button type="button" onclick="location.href='top'">トップ</button>
            <button type="button" class="active">開催中イベント一覧</button>
            <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button>
            <button type="button" onclick="location.href='closedevents'">終了イベント一覧</button>
        </div>
        
        <h1>開催中イベント一覧</h1>
        
        <div class="filter-section">
            <div style="margin-top: 15px;">
                {{ if eq .Mode 1 }}
                <button class="data-filter-btn active" onclick="toggleDataFilter()">
                    ✓ 獲得ポイントデータ取得中のイベントのみ表示
                </button>
                {{ else }}
                <button class="data-filter-btn" onclick="toggleDataFilter()">
                    獲得ポイントデータ取得中のイベントのみ表示
                </button>
                {{ end }}
            </div>
        </div>
        
        <div class="info-message">
            イベント数： {{ .Totalcount }}
        </div>
        <!-- イベント一覧テーブル -->
        <table class="events-table">
            <tr>
                <th>イベント名とイベントページへのリンク</th>
                <th>開始日時</th>
                <th>終了日時</th>
                <th>参加ルーム一覧</th>
                <th>表示項目選択画面/<br>データ取得開始登録</th>
                <th>直近獲得<br>ポイント表</th>
                <th>獲得ポイント<br>推移図</th>
                <th>日々の<br>獲得pt</th>
                <th>枠毎の<br>獲得pt</th>
            </tr>

        {{ $i := 0 }}
        {{ range .Eventinflist }}
        {{ if eq $i 1 }}
            {{ if eq .Aclr 0 }}
            <tr bgcolor="gainsboro">
            {{ else }}
            <tr bgcolor="palegreen">
            {{ end }}
            {{ $i = 0 }}
        {{ else }}
            {{ if eq .Aclr 0 }}
            <tr>
            {{ else }}
            <tr bgcolor="lightblue">
            {{ end }}
            {{ $i = 1 }}
        {{ end }}

                <td>
                    <a href="https://showroom-live.com/event/{{ .Event_ID }}">{{ .Event_name }}</a>
                </td>
                <td>
                    {{ TimeToString .Start_time }}
                </td>
                <td>
                    {{ TimeToString .End_time }}
                </td>
                <td style="text-align: center;">
                    {{/*
                    <a href="eventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}">参加ルーム一覧</a>
                    */}}
                    改修中
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="eventtop?eventid={{ .Event_ID }}">項目選択</a>
                    {{ else }}
                    <a href="new-event?eventid={{ .Event_ID }}">取得開始登録</a>
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="list-last?eventid={{ .Event_ID }}">リスト</a>
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="graph-total?eventid={{ .Event_ID }}">グラフ</a>
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="list-perday?eventid={{ .Event_ID }}">日々pt</a>
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="list-perslot?eventid={{ .Event_ID }}">枠毎pt</a>
                    {{ end }}
                </td>
            </tr>
            {{ end }}
        </table>
        
        {{ if .ErrMsg }}
        <div class="error-message">{{ .ErrMsg }}</div>
        {{ end }}
    </div>
    
    <script>
        function toggleDataFilter() {
            const url = window.location.href;
            const urlObj = new URL(url);
            const currentMode = urlObj.searchParams.get('mode');
            
            if (currentMode === '1') {
                urlObj.searchParams.delete('mode');
            } else {
                urlObj.searchParams.set('mode', '1');
            }
            
            location.href = urlObj.toString();
        }
    </script>
</body>
</html>