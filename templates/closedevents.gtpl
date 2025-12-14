<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>終了イベント一覧</title>
    {{/* Turnstile 1 */}}
    {{if .TurnstileSiteKey}}
        <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    {{end}}
    {{/* ----------- */}}

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
        
        .filter-buttons {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
            margin-bottom: 15px;
        }
        
        .filter-btn {
            padding: 10px 20px;
            border: 2px solid #007bff;
            background-color: white;
            color: #007bff;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.2s;
        }
        
        .filter-btn:hover {
            background-color: #007bff;
            color: white;
        }
        
        .filter-btn.active {
            background-color: #007bff;
            color: white;
        }
        
        .current-filters {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
            align-items: center;
        }
        
        .filter-tag {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 6px 12px;
            background-color: #e7f3ff;
            border: 1px solid #007bff;
            border-radius: 4px;
            font-size: 13px;
        }
        
        .filter-tag .remove {
            cursor: pointer;
            color: #007bff;
            font-weight: bold;
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
        
        /* モーダルスタイル */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            overflow: auto;
            background-color: rgba(0,0,0,0.4);
        }
        
        .modal-content {
            background-color: #fefefe;
            margin: 5% auto;
            padding: 0;
            border: 1px solid #888;
            width: 90%;
            max-width: 600px;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        .modal-header {
            padding: 15px 20px;
            background-color: #007bff;
            color: white;
            border-radius: 8px 8px 0 0;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .modal-header h2 {
            margin: 0;
            font-size: 18px;
        }
        
        .close {
            color: white;
            font-size: 28px;
            font-weight: bold;
            cursor: pointer;
            line-height: 1;
        }
        
        .close:hover,
        .close:focus {
            color: #f0f0f0;
        }
        
        .modal-body {
            padding: 20px;
        }
        
        .modal-footer {
            padding: 15px 20px;
            background-color: #f8f9fa;
            border-radius: 0 0 8px 8px;
            display: flex;
            justify-content: flex-end;
            gap: 10px;
        }
        
        .form-group {
            margin-bottom: 15px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: 500;
            color: #333;
        }
        
        .form-group input {
            width: 100%;
            padding: 8px 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
            box-sizing: border-box;
        }
        
        .form-group .help-text {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
        
        .btn {
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
        }
        
        .btn-primary {
            background-color: #007bff;
            color: white;
        }
        
        .btn-primary:hover {
            background-color: #0056b3;
        }
        
        .btn-primary:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        
        .btn-secondary {
            background-color: #6c757d;
            color: white;
        }
        
        .btn-secondary:hover {
            background-color: #545b62;
        }
        
        /* ルーム選択リスト */
        .room-list {
            max-height: 300px;
            overflow-y: auto;
            border: 1px solid #ddd;
            border-radius: 4px;
            margin-top: 10px;
        }
        
        .room-item {
            padding: 10px;
            border-bottom: 1px solid #eee;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        
        .room-item:hover {
            background-color: #f0f0f0;
        }
        
        .room-item:last-child {
            border-bottom: none;
        }
        
        .room-item.selected {
            background-color: #e7f3ff;
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
        
        .events-table tr:nth-child(even) {
            background-color: gainsboro;
        }
        
         {{/* 色が明るすぎるのも問題だが、参照済みでも色が変わらないのがもっと問題
        .events-table a {
            color: #007bff;
            color: #000488;
            text-decoration: none;
        }
        */}}
        
        .events-table a:hover {
            text-decoration: underline;
        }
        
        .pagination {
            margin: 20px 0;
            display: flex;
            gap: 10px;
            align-items: center;
        }
        
        .pagination button {
            padding: 8px 16px;
            border: 1px solid #007bff;
            background-color: white;
            color: #007bff;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .pagination button:hover {
            background-color: #007bff;
            color: white;
        }
        
        .pagination button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
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
        
        .loading {
            text-align: center;
            padding: 20px;
            color: #666;
        }
        
        .no-results {
            text-align: center;
            padding: 40px;
            color: #999;
            font-size: 16px;
        }
    </style>
</head>
<body>
    {{ $tn := .TimeNow }}
    
    <div class="container">
        <div class="nav-buttons">
            <button type="button" onclick="location.href='top'">トップ</button>
            <button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button>
            <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button>
            <button type="button" class="active">終了イベント一覧</button>
        </div>
    {{if .TurnstileSiteKey }}
        	<!-- Turnstileチャレンジ表示 -->
			<div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9;">
				<h3>セキュリティチェック</h3>
				{{if .TurnstileError}}
				<p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
				{{end}}
				<p>終了済みイベントを表示するには、セキュリティチェックを完了してください。</p>
				<p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
				<form method="POST" action="closedevents">
					<input type="hidden" name="mode" value="{{.Mode}}">
					<input type="hidden" name="keywordev" value="{{.Keywordev}}">
					<input type="hidden" name="keywordrm" value="{{.Keywordrm}}">
					<input type="hidden" name="kwevid" value="{{.Kwevid}}">
					<input type="hidden" name="userno" value="{{.Userno}}">
					<input type="hidden" name="path" value="{{.Path}}">

					<input type="hidden" name="limit" value="{{.Limit}}">
					<input type="hidden" name="offset" value="{{.Offset}}">
					<input type="hidden" name="action" value="{{.Action}}">

					<input type="hidden" name="requestid" value="{{.RequestID}}">
					<div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
					<br>
					<button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
				</form>
			</div>
			<br><br>
    {{else}}

        
        <h1>終了イベント一覧</h1>
        
        <div class="filter-section">
            <div class="filter-buttons">
                <button class="filter-btn" onclick="openEventNameModal()">
                    <span>📝</span> イベント名で検索
                </button>
                <button class="filter-btn" onclick="openEventIdModal()">
                    <span>🔑</span> イベントIDで検索
                </button>
                <button class="filter-btn" onclick="openRoomNameModal()">
                    <span>🏠</span> ルーム名で検索
                </button>
                <button class="filter-btn" onclick="openRoomIdModal()">
                    <span>🆔</span> ルームIDで検索
                </button>
            </div>
            
            <div class="current-filters" id="currentFilters">
                {{ if ne .Keywordev "" }}
                <div class="filter-tag">
                    <strong>イベント名:</strong> {{ .Keywordev }}
                    <span class="remove" onclick="clearFilter('keywordev')">✕</span>
                </div>
                {{ end }}
                {{ if ne .Kwevid "" }}
                <div class="filter-tag">
                    <strong>イベントID:</strong> {{ .Kwevid }}
                    <span class="remove" onclick="clearFilter('kwevid')">✕</span>
                </div>
                {{ end }}
                {{ if ne .Keywordrm "" }}
                <div class="filter-tag">
                    <strong>ルーム名:</strong> {{ .Keywordrm }}
                    <span class="remove" onclick="clearFilter('keywordrm')">✕</span>
                </div>
                {{ end }}
                {{ if ne .Userno 0 }}
                <div class="filter-tag">
                    <strong>ルームID:</strong> {{ .Userno }}
                    <span class="remove" onclick="clearFilter('userno')">✕</span>
                </div>
                {{ end }}
            </div>
            
            <div style="margin-top: 15px;">
                {{ if eq .Mode 1 }}
                <button class="data-filter-btn active" onclick="toggleDataFilter()">
                    ✓ 獲得ポイント詳細データのあるイベントのみ表示
                </button>
                {{ else }}
                <button class="data-filter-btn" onclick="toggleDataFilter()">
                    獲得ポイント詳細データのあるイベントのみ表示
                </button>
                {{ end }}
            </div>
        </div>
        
        {{ if .ErrMsg }}
        <div class="error-message">{{ .ErrMsg }}</div>
        {{ end }}
        
        <div class="info-message">
            一覧は51件程度ずつ表示され、50件ずつスクロールされます。データ上最終結果が存在しても表示されないケースがあります。
        </div>
        
        <!-- ページネーション（上部） -->
        {{ if gt .Totalcount 0 }}
        <div class="pagination">
            {{ if ne .Offset 0 }}
            <button onclick="navigatePage('top')">最初</button>
            <button onclick="navigatePage('prev')">前ページ</button>
            {{ end }}
            <span>表示中: {{ .Totalcount }}件</span>
            {{ if eq .Totalcount .Limit }}
            <button onclick="navigatePage('next')">次ページ</button>
            {{ end }}
        </div>
        {{ end }}
        
        <!-- イベント一覧テーブル -->
        {{ if gt .Totalcount 0 }}
        <table class="events-table">
            <tr>
                <th style="border-right: none;">イベント名とイベントページへのリンク</th>
                <th style="border-left: none;">イベントID</th>
                <th>開始日時</th>
                <th>終了日時</th>
                <th>最終結果</th>
                <th>表示項目<br>選択画面</th>
                <th>最終獲得<br>ポイント表</th>
                <th>獲得ポイント<br>推移図</th>
                <th>日々の<br>獲得pt</th>
                <th>枠毎の<br>獲得pt</th>
                <th>貢献<br>pt</th>
            </tr>
            {{ $i := 0 }}
            {{ $userno := .Userno }}
            {{ range .Eventinflist }}

            {{/*
            {{ if eq $i 1 }}
            <tr>
            {{ $i = 0 }}
            {{ else }}
            <tr>
            {{ $i = 1 }}
            {{ end }}
            */}}


        {{ if eq $i 1 }}
            {{ if eq .Highlighted 0 }}
            <tr bgcolor="gainsboro">
            {{ else }}
            <tr bgcolor="palegreen">
            {{ end }}
            {{ $i = 0 }}
        {{ else }}
            {{ if eq .Highlighted 0 }}
            <tr>
            {{ else }}
            <tr bgcolor="lightblue">
            {{ end }}
            {{ $i = 1 }}
        {{ end }}
        
                <td style="border-right: none;">
                    {{ if IsTempID .Event_ID }}
                        {{ .Event_name }}
                    {{ else }}
                        <a href="https://showroom-live.com/event/{{ .Event_ID }}">{{ .Event_name }}</a>
                    {{ end }}
                </td>
                <td style="border-left: none; text-align: right;">
                    {{ if ne .I_Event_ID 0 }}
                        {{ .I_Event_ID }}
                    {{ end }}
                </td>
                <td>
                    {{ TimeToStringY .Start_time }}
                </td>
                <td>
                    {{ TimeToString .End_time }}
                </td>
                <td style="text-align: center;">
                    {{/* {{ if and (ne .I_Event_ID 0) ( ne .Aclr 0 ) }} */}}
                    {{/* {{ if and (ne .I_Event_ID 0) ( ne .Aclr 0 ) }} */}}
                    {{ if  ne .Aclr 0 }}
                    {{/* {{ if ne .I_Event_ID 0 }} */}}
                    
                        {{ if eq .Rstatus "ProvisionalC" }}
                    <a href="closedeventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}{{ if ne $userno 0 }}&roomid={{$userno}}{{end}}">
                        *暫定結果*
                    </a>
                        {{ else if ne .Rstatus "Provisional" }}
                    <a href="closedeventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}{{ if ne $userno 0 }}&roomid={{$userno}}{{end}}">
                        最終結果
                    </a>
                        {{ end }}

                    {{ else }}
                        {{ if ne $userno 0 }}
                        <a href="https://showroom-lSaveConfirmedData_200200_200504_200400ive.com/api/events/{{.I_Event_ID}}/ranking?room_id={{$userno}}">API</a>
                        {{ end }}
                    {{ end }}

                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="eventtop?eventid={{ .Event_ID }}">表示項目選択</a>
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="list-lastP?eventid={{ .Event_ID }}{{ if ne $userno 0 }}&roomid={{$userno}}{{end}}">リスト</a>
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
                <td style="text-align: center;">
                    {{ if ne $userno 0 }}

                    {{ if IsTempID .Event_ID }}
                    ーー /
                    {{ else }}
                    <a href="https://www.showroom-live.com/event/contribution/{{ DelBlockID .Event_ID }}?room_id={{$userno}}">公式</a> /
                    {{ end }}
                    {{ if ne .I_Event_ID 0 }}
                    <a href="contributors?ieventid={{.I_Event_ID}}&roomid={{$userno}}"  target="_blank" rel="noopener noreferrer">CSV</a>
                    {{ else }}
                    ーー
                    {{ end }}
                    {{ end }}
                </td>
            </tr>
            {{ end }}
        </table>
        
        <!-- ページネーション（下部） -->
        {{ if ne .Totalcount 0 }}
        <div class="pagination">
            {{ if ne .Offset 0 }}
            <button onclick="navigatePage('top')">最初</button>
            <button onclick="navigatePage('prev')">前ページ</button>
            {{ end }}
            {{ if eq .Totalcount .Limit }}
            <button onclick="navigatePage('next')">次ページ</button>
            {{ end }}
        </div>
        {{ end }}
        
        {{ else }}
        <div class="no-results">
            {{ if or (ne .Keywordev "") (ne .Kwevid "") (ne .Keywordrm "") (ne .Userno 0) }}
            検索条件に一致するイベントが見つかりませんでした。
            {{ else }}
            イベントデータがありません。
            {{ end }}
        </div>
        {{ end }}
    </div>
    
    <!-- イベント名検索モーダル -->
    <div id="eventNameModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2>イベント名で検索</h2>
                <span class="close" onclick="closeModal('eventNameModal')">&times;</span>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label for="eventNameInput">イベント名に含まれる文字列</label>
                    <input type="text" id="eventNameInput" value="{{ .Keywordev }}" placeholder="例: スタートダッシュ、花火、Music">
                    <div class="help-text">
                        カタカナとひらがな、全角と半角、英大文字と小文字はアバウトに検索されます。
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" onclick="closeModal('eventNameModal')">キャンセル</button>
                <button class="btn btn-primary" onclick="applyEventNameFilter()">検索実行</button>
            </div>
        </div>
    </div>
    
    <!-- イベントID検索モーダル -->
    <div id="eventIdModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2>イベントIDで検索</h2>
                <span class="close" onclick="closeModal('eventIdModal')">&times;</span>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label for="eventIdInput">イベントIDに含まれる文字列</label>
                    <input type="text" id="eventIdInput" value="{{ .Kwevid }}" placeholder="例: fireworks" pattern="[0-9A-Za-z-_=?]+">
                    <div class="help-text">
                        イベントページのURLの最後のフィールドです（本来のイベントIDは5桁程度の整数です）。<br>
                        例えば花火のイベントであれば"fireworks"が含まれていることが多いです。
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" onclick="closeModal('eventIdModal')">キャンセル</button>
                <button class="btn btn-primary" onclick="applyEventIdFilter()">検索実行</button>
            </div>
        </div>
    </div>
    
    <!-- ルーム名検索モーダル -->
    <div id="roomNameModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2>ルーム名で検索</h2>
                <span class="close" onclick="closeModal('roomNameModal')">&times;</span>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label for="roomNameInput">ルーム名に含まれる文字列</label>
                    <input type="text" id="roomNameInput" value="{{ .Keywordrm }}" placeholder="ルーム名の一部を入力">
                    <div class="help-text">
                        現在のルーム名だけでなく過去のルーム名（のうち最近のもの）も検索対象となります。<br>
                        ルームの検索結果は50件までしか表示されませんので、具体的な文字列を入力してください。
                    </div>
                </div>
                <button class="btn btn-primary" onclick="searchRooms()" style="width: 100%;">ルーム検索</button>
                
                <div id="roomListContainer" style="display: none;">
                    <div class="form-group" style="margin-top: 15px;">
                        <label>検索結果からルームを選択してください</label>
                        <div class="room-list" id="roomList"></div>
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" onclick="closeModal('roomNameModal')">キャンセル</button>
                <button class="btn btn-primary" id="applyRoomBtn" onclick="applyRoomNameFilter()" disabled>選択したルームで検索</button>
            </div>
        </div>
    </div>
    
    <!-- ルームID検索モーダル -->
    <div id="roomIdModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2>ルームIDで検索</h2>
                <span class="close" onclick="closeModal('roomIdModal')">&times;</span>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label for="roomIdInput">ルームID</label>
                    <input type="number" id="roomIdInput" value="{{ if ne .Userno 0 }}{{ .Userno }}{{ end }}" placeholder="例: 123456">
                    <div class="help-text">
                        ルームIDはプロフィールやファンルームのURLの最後の"room_id="のあとにある整数です（6桁が多い）。<br>
                        ルームIDの一部を指定しての検索はできません。
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" onclick="closeModal('roomIdModal')">キャンセル</button>
                <button class="btn btn-primary" onclick="applyRoomIdFilter()">検索実行</button>
            </div>
        </div>
    </div>
    {{ end }}
    </div>
    
    <script>
        let selectedRoomUserno = 0;
        let selectedRoomName = '';
        
        // 現在のフィルター状態を保持
        const currentFilters = {
            mode: {{ .Mode }},
            keywordev: "{{ .Keywordev }}",
            kwevid: "{{ .Kwevid }}",
            keywordrm: "{{ .Keywordrm }}",
            userno: {{ .Userno }},
            limit: {{ .Limit }},
            offset: {{ .Offset }}
        };
        
        // モーダルを開く
        function openEventNameModal() {
            document.getElementById('eventNameModal').style.display = 'block';
        }
        
        function openEventIdModal() {
            document.getElementById('eventIdModal').style.display = 'block';
        }
        
        function openRoomNameModal() {
            document.getElementById('roomNameModal').style.display = 'block';
            document.getElementById('roomListContainer').style.display = 'none';
            document.getElementById('applyRoomBtn').disabled = true;
        }
        
        function openRoomIdModal() {
            document.getElementById('roomIdModal').style.display = 'block';
        }
        
        // モーダルを閉じる
        function closeModal(modalId) {
            document.getElementById(modalId).style.display = 'none';
        }
        
        // モーダル外クリックで閉じる
        window.onclick = function(event) {
            if (event.target.classList.contains('modal')) {
                event.target.style.display = 'none';
            }
        }
        
        // イベント名フィルター適用
        function applyEventNameFilter() {
            const keyword = document.getElementById('eventNameInput').value.trim();
            currentFilters.keywordev = keyword;
            currentFilters.kwevid = '';
            currentFilters.keywordrm = '';
            currentFilters.userno = 0;
            currentFilters.offset = 0;
            navigateWithFilters();
        }
        
        // イベントIDフィルター適用
        function applyEventIdFilter() {
            const keyword = document.getElementById('eventIdInput').value.trim();
            currentFilters.kwevid = keyword;
            currentFilters.keywordev = '';
            currentFilters.keywordrm = '';
            currentFilters.userno = 0;
            currentFilters.offset = 0;
            navigateWithFilters();
        }
        
        // ルーム検索（Ajax）
        async function searchRooms() {
            const keyword = document.getElementById('roomNameInput').value.trim();
            if (!keyword) {
                alert('ルーム名を入力してください。');
                return;
            }
            
            const roomListContainer = document.getElementById('roomListContainer');
            const roomList = document.getElementById('roomList');
            roomList.innerHTML = '<div class="loading">検索中...</div>';
            roomListContainer.style.display = 'block';
            
            try {
                const response = await fetch(`/api/search-rooms?keyword=${encodeURIComponent(keyword)}&limit=50`);
                if (!response.ok) {
                    throw new Error('検索に失敗しました');
                }
                
                const rooms = await response.json();
                
                if (rooms.length === 0) {
                    roomList.innerHTML = '<div style="padding: 20px; text-align: center; color: #999;">該当するルームが見つかりませんでした</div>';
                    return;
                }
                
                roomList.innerHTML = '';
                rooms.forEach(room => {
                    const div = document.createElement('div');
                    div.className = 'room-item';
                    div.textContent = room.User_name;
                    div.onclick = function() {
                        // 既存の選択を解除
                        document.querySelectorAll('.room-item').forEach(item => {
                            item.classList.remove('selected');
                        });
                        // 新しい選択を設定
                        div.classList.add('selected');
                        selectedRoomUserno = room.Userno;
                        selectedRoomName = room.User_name;
                        document.getElementById('applyRoomBtn').disabled = false;
                    };
                    roomList.appendChild(div);
                });
            } catch (error) {
                roomList.innerHTML = `<div style="padding: 20px; text-align: center; color: #dc3545;">エラー: ${error.message}</div>`;
            }
        }
        
        // ルーム名フィルター適用
        function applyRoomNameFilter() {
            if (selectedRoomUserno === 0) {
                alert('ルームを選択してください。');
                return;
            }
            currentFilters.userno = selectedRoomUserno;
            currentFilters.keywordrm = selectedRoomName;
            currentFilters.keywordev = '';
            currentFilters.kwevid = '';
            currentFilters.offset = 0;
            navigateWithFilters();
        }
        
        // ルームIDフィルター適用
        function applyRoomIdFilter() {
            const roomId = parseInt(document.getElementById('roomIdInput').value);
            if (!roomId || roomId <= 0) {
                alert('有効なルームIDを入力してください。');
                return;
            }
            currentFilters.userno = roomId;
            currentFilters.keywordrm = '';
            currentFilters.keywordev = '';
            currentFilters.kwevid = '';
            currentFilters.offset = 0;
            navigateWithFilters();
        }
        
        // データフィルタートグル
        function toggleDataFilter() {
            currentFilters.mode = currentFilters.mode === 1 ? 0 : 1;
            currentFilters.offset = 0;
            navigateWithFilters();
        }
        
        // フィルタークリア
        function clearFilter(filterName) {
            if (filterName === 'keywordev') {
                currentFilters.keywordev = '';
            } else if (filterName === 'kwevid') {
                currentFilters.kwevid = '';
            } else if (filterName === 'keywordrm') {
                currentFilters.keywordrm = '';
                currentFilters.userno = 0;
            } else if (filterName === 'userno') {
                currentFilters.userno = 0;
                currentFilters.keywordrm = '';
            }
            currentFilters.offset = 0;
            navigateWithFilters();
        }
        
        // ページナビゲーション
        function navigatePage(action) {
            if (action === 'next') {
                currentFilters.offset += currentFilters.limit - 1;
            } else if (action === 'prev') {
                currentFilters.offset -= currentFilters.limit - 1;
                if (currentFilters.offset < 0) currentFilters.offset = 0;
            } else if (action === 'top') {
                currentFilters.offset = 0;
            }
            navigateWithFilters();
        }
        
        // フィルター付きナビゲーション
        function navigateWithFilters() {
            const params = new URLSearchParams();
            params.append('mode', currentFilters.mode);
            params.append('limit', currentFilters.limit);
            params.append('offset', currentFilters.offset);
            
            if (currentFilters.keywordev) {
                params.append('keywordev', currentFilters.keywordev);
                params.append('path', '1');
            } else if (currentFilters.kwevid) {
                params.append('kwevid', currentFilters.kwevid);
                params.append('path', '2');
            } else if (currentFilters.userno && currentFilters.keywordrm) {
                params.append('userno', currentFilters.userno);
                params.append('keywordrm', currentFilters.keywordrm);
                params.append('path', '4');
            } else if (currentFilters.userno) {
                params.append('userno', currentFilters.userno);
                params.append('path', '5');
            } else {
                params.append('path', '0');
            }
            
            location.href = 'closedevents?' + params.toString();
        }
    </script>
</body>

</html>