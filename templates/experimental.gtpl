<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>保守・資料・実験</title>
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
        
        h1 {
            color: #333;
            margin: 20px 0;
            font-size: 24px;
        }
        
    </style>
</head>

<body>
<div class="nav-buttons">
  <button type="button" onclick="location.href='top'">トップ</button>
  <button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button>
  <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button>
  <button type="button" onclick="location.href='closedevents'">終了イベント一覧</button>
  <button type="button" class="active">保守・資料・実験</button>
</div>
<h2>保守・資料・実験</h2>
  <ol type="A">
  <li><a href="list-todo">不具合・改善項目</a></li>
  <li>アクセス情報</li>
    <ol type="1">
    <li><a href="accessstats">日々のアクセス数</a></li>
    <li><a href="accessstats">ハンドラーごとのアクセス数</a></li>
    <li><a href="accessstats">イベントごとのアクセス数</a></li>
    </ol>
  <li>SHOWROOMのAPI</li>
  {{/*
  <li><a href="">SHOWROOMのAPI</a></li>
  {{/*
    <ol type="1">
      <li></li>
      <li></li>
    </ol>
    */}}
  </ol>
</body>

</html>