<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>ToDo編集</title>
    <style type="text/css">
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        .edit-form {
            background-color: #e9f5ff;
            padding: 20px;
            border: 1px solid #b3d9ff;
            border-radius: 5px;
            max-width: 800px;
            margin: 20px auto;
        }
        .edit-form input[type="text"],
        .edit-form textarea {
            width: 100%;
            padding: 8px;
            margin: 5px 0;
            border: 1px solid #ddd;
            border-radius: 3px;
            box-sizing: border-box;
        }
        .edit-form button {
            padding: 10px 20px;
            margin: 5px;
            background-color: #008CBA;
            color: white;
            border: none;
            cursor: pointer;
            border-radius: 3px;
        }
        .edit-form button:hover {
            background-color: #007399;
        }
        .edit-form button.cancel {
            background-color: #999;
        }
        .edit-form button.cancel:hover {
            background-color: #777;
        }
        .error {
            color: red;
            font-weight: bold;
            margin: 10px 0;
        }
        .success {
            color: green;
            font-weight: bold;
            margin: 10px 0;
        }
        .form-row {
            margin: 15px 0;
        }
        .form-label {
            display: block;
            font-weight: bold;
            margin-bottom: 5px;
        }
        .info-row {
            margin: 10px 0;
            padding: 10px;
            background-color: #f5f5f5;
            border-radius: 3px;
        }
        .info-label {
            font-weight: bold;
            display: inline-block;
            width: 120px;
        }
        .checkbox-row {
            margin: 15px 0;
        }
        .checkbox-row input[type="checkbox"] {
            width: auto;
            margin-right: 10px;
        }
        .checkbox-row label {
            font-weight: bold;
        }
    </style>
</head>
<body>
    <h1>ToDo編集</h1>
    
    {{ if .ErrMsg }}
    <div class="error">{{ .ErrMsg }}</div>
    {{ end }}

    {{ if .SuccessMsg }}
    <div class="success">{{ .SuccessMsg }}</div>
    {{ end }}

    <div class="edit-form">
        <h3>インシデント編集 (ID: {{ .Todo.ID }})</h3>
        
        <!-- 参考情報：発生日と解決日 -->
        <div class="info-row">
            <span class="info-label">発生日:</span>
            <span>{{ FormatTime .Todo.Ts }}</span>
        </div>
        {{ if .Todo.Closed }}
        <div class="info-row">
            <span class="info-label">解決日:</span>
            <span>{{ FormatTimePtr .Todo.Closed }}</span>
        </div>
        {{ end }}

        <form method="POST" action="/edit-todo">
            <!-- IDと検索条件を引き継ぐための隠しフィールド -->
            <input type="hidden" name="id" value="{{ .Todo.ID }}">
            <input type="hidden" name="itype_search" value="{{ .Itype }}">
            <input type="hidden" name="target_search" value="{{ .Target }}">
            <input type="hidden" name="issue_search" value="{{ .Issue }}">
            <input type="hidden" name="solution_search" value="{{ .Solution }}">
            
            <div class="form-row">
                <label class="form-label">種別 *:</label>
                <input type="text" name="itype" value="{{ .Todo.Itype }}" required placeholder="BUG, FIXME, OPTIMIZE, HACK, REVIEW, TODO等">
            </div>
            
            <div class="form-row">
                <label class="form-label">対象 *:</label>
                <input type="text" name="target" value="{{ .Todo.Target }}" required placeholder="モジュール/機能(関数)">
            </div>
            
            <div class="form-row">
                <label class="form-label">課題 *:</label>
                <textarea name="issue" rows="4" required placeholder="追加すべき機能、改善すべき不具合等">{{ .Todo.Issue }}</textarea>
            </div>
            
            <div class="form-row">
                <label class="form-label">対応内容:</label>
                <textarea name="solution" rows="4" placeholder="対応内容と結果（任意）">{{ .Todo.Solution }}</textarea>
            </div>
            
            <div class="checkbox-row">
                <input type="checkbox" id="completed" name="completed" {{ if .Todo.Closed }}checked{{ end }}>
                <label for="completed">完了（チェックを入れると完了日が現在時刻に設定されます）</label>
            </div>
            
            <div class="form-row">
                <button type="submit">更新</button>
                <button type="button" class="cancel" onclick="history.back()">キャンセル</button>
            </div>
        </form>
    </div>

    <br>
    <p><a href="/list-todo?itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}">ToDo一覧に戻る</a></p>
    <p><a href="/top">トップページに戻る</a></p>

</body>
</html>
