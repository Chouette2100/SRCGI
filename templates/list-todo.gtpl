<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>ToDo管理</title>
    <style type="text/css">
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin-top: 20px;
        }
        th, td {
            border: solid 1px #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #4CAF50;ToDo, FixMe
            color: white;
        }
        tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        tr:hover {
            background-color: #ddd;
        }
        .search-form {
            background-color: #f9f9f9;
            padding: 15px;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        .search-form input[type="text"] {
            width: 200px;
            padding: 5px;
            margin: 5px;
        }
        .search-form button {
            padding: 5px 15px;
            margin: 5px;
            background-color: #4CAF50;
            color: white;
            border: none;
            cursor: pointer;
            border-radius: 3px;
        }
        .search-form button:hover {
            background-color: #45a049;
        }
        .add-form {
            background-color: #e9f5ff;
            padding: 15px;
            border: 1px solid #b3d9ff;
            border-radius: 5px;
            margin-top: 20px;
        }
        .add-form input[type="text"],
        .add-form textarea {
            width: 100%;
            padding: 8px;
            margin: 5px 0;
            border: 1px solid #ddd;
            border-radius: 3px;
        }
        .add-form button {
            padding: 10px 20px;
            margin: 5px;
            background-color: #008CBA;
            color: white;
            border: none;
            cursor: pointer;
            border-radius: 3px;
        }
        .add-form button:hover {
            background-color: #007399;
        }
        .paging {
            margin: 20px 0;
            text-align: center;
        }
        .paging button {
            padding: 8px 16px;
            margin: 0 5px;
            background-color: #555;
            color: white;
            border: none;
            cursor: pointer;
            border-radius: 3px;
        }
        .paging button:hover {
            background-color: #333;
        }
        .paging button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        .error {
            color: red;
            font-weight: bold;
        }
        .form-row {
            margin: 10px 0;
        }
        .form-label {
            display: inline-block;
            width: 100px;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <h1>ToDo管理</h1>
    
    {{ if .ErrMsg }}
    <div class="error">{{ .ErrMsg }}</div>
    {{ end }}

    <!-- 検索フォーム -->
    <div class="search-form">
        <h3>検索条件</h3>
        <form method="GET" action="/list-todo">
            <div class="form-row">
                <span class="form-label">種別:</span>
                <input type="text" name="itype" value="{{ .Itype }}" placeholder="BUG, FIXME, OPTIMIZE, HACK, REVIEW, TODO等">
            </div>
            <div class="form-row">
                <span class="form-label">対象:</span>
                <input type="text" name="target" value="{{ .Target }}" placeholder="モジュール/機能(関数)">
            </div>
            <div class="form-row">
                <span class="form-label">課題:</span>
                <input type="text" name="issue" value="{{ .Issue }}" placeholder="キーワード（部分一致）">
            </div>
            <div class="form-row">
                <span class="form-label">対応内容:</span>
                <input type="text" name="solution" value="{{ .Solution }}" placeholder="キーワード（部分一致）">
            </div>
            <button type="submit">検索</button>
            <button type="button" onclick="location.href='/list-todo'">クリア</button>
        </form>
    </div>

    <!-- ページング -->
    {{ if or .HasPrev .HasNext }}
    <div class="paging">
        {{ if .HasPrev }}
        <button onclick="location.href='/list-todo?dir=prev&maxid={{ .MaxID }}&itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}'">前のページ</button>
        {{ else }}
        <button disabled>前のページ</button>
        {{ end }}
        
        <button onclick="location.href='/list-todo?itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}'">再表示</button>
        
        {{ if .HasNext }}
        <button onclick="location.href='/list-todo?dir=next&minid={{ .MinID }}&itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}'">次のページ</button>
        {{ else }}
        <button disabled>次のページ</button>
        {{ end }}
    </div>
    {{ end }}

    <!-- ToDoリスト -->
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>発生日</th>
                <th>種別</th>
                <th>対象</th>
                <th>課題</th>
                <th>対応内容</th>
                <th>解決日</th>
                <th>操作</th>
            </tr>
        </thead>
        <tbody>
            {{ if .Todos }}
                {{ range .Todos }}
                <tr>
                    <td>{{ .ID }}</td>
                    <td>{{ FormatTime .Ts }}</td>
                    <td>{{ .Itype }}</td>
                    <td>{{ .Target }}</td>
                    <td>{{ .Issue }}</td>
                    <td>{{ .Solution }}</td>
                    <td>{{ FormatTimePtr .Closed }}</td>
                    <td>
                        <a href="/edit-todo?id={{ .ID }}&itype_search={{ $.Itype }}&target_search={{ $.Target }}&issue_search={{ $.Issue }}&solution_search={{ $.Solution }}" style="color: #008CBA; text-decoration: none;">編集</a>
                    </td>
                </tr>
                {{ end }}
            {{ else }}
                <tr>
                    <td colspan="8" style="text-align: center;">データがありません</td>
                </tr>
            {{ end }}
        </tbody>
    </table>

    <!-- ページング（下部） -->
    {{ if or .HasPrev .HasNext }}
    <div class="paging">
        {{ if .HasPrev }}
        <button onclick="location.href='/list-todo?dir=prev&maxid={{ .MaxID }}&itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}'">前のページ</button>
        {{ else }}
        <button disabled>前のページ</button>
        {{ end }}
        
        <button onclick="location.href='/list-todo?itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}'">再表示</button>
        
        {{ if .HasNext }}
        <button onclick="location.href='/list-todo?dir=next&minid={{ .MinID }}&itype={{ .Itype }}&target={{ .Target }}&issue={{ .Issue }}&solution={{ .Solution }}'">次のページ</button>
        {{ else }}
        <button disabled>次のページ</button>
        {{ end }}
    </div>
    {{ end }}

    <!-- 新規追加フォーム -->
    <div class="add-form">
        <h3>新規インシデント追加</h3>
        <form method="POST" action="/insert-todo">
            <!-- 検索条件を引き継ぐための隠しフィールド -->
            <input type="hidden" name="itype_search" value="{{ .Itype }}">
            <input type="hidden" name="target_search" value="{{ .Target }}">
            <input type="hidden" name="issue_search" value="{{ .Issue }}">
            <input type="hidden" name="solution_search" value="{{ .Solution }}">
            
            <div class="form-row">
                <span class="form-label">種別 *:</span>
                <input type="text" name="itype" required placeholder="BUG, FIXME, OPTIMIZE, HACK, REVIEW, TODO等">
            </div>
            <div class="form-row">
                <span class="form-label">対象 *:</span>
                <input type="text" name="target" required placeholder="関数/機能(関数)">
            </div>
            <div class="form-row">
                <span class="form-label">課題 *:</span>
                <textarea name="issue" rows="3" required placeholder="追加すべき機能、改善すべき不具合等"></textarea>
            </div>
            <div class="form-row">
                <span class="form-label">対応内容:</span>
                <textarea name="solution" rows="3" placeholder="対応内容と結果（任意）"></textarea>
            </div>
            <button type="submit">追加</button>
            <button type="reset">クリア</button>
        </form>
    </div>

    <br>
    <p><a href="/top">トップページに戻る</a></p>

</body>
</html>
