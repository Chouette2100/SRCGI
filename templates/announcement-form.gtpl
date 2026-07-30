{{ define "announcement-form" }}
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>告知入力フォーム</title>
    <style>
        body {
            font-family: sans-serif;
            margin: 40px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 600px;
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-top: 0;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input[type="text"],
        input[type="color"],
        textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
            font-family: sans-serif;
        }
        textarea {
            resize: vertical;
            min-height: 100px;
        }
        .checkbox-group {
            display: flex;
            align-items: center;
        }
        input[type="checkbox"] {
            margin-right: 10px;
            width: 20px;
            height: 20px;
            cursor: pointer;
        }
        .checkbox-group label {
            margin: 0;
            cursor: pointer;
            font-weight: normal;
        }
        .button-group {
            display: flex;
            gap: 10px;
            justify-content: center;
            margin-top: 25px;
        }
        button {
            padding: 10px 30px;
            font-size: 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .btn-submit {
            background-color: #007bff;
            color: white;
        }
        .btn-submit:hover {
            background-color: #0056b3;
        }
        .btn-reset {
            background-color: #6c757d;
            color: white;
        }
        .btn-reset:hover {
            background-color: #545b62;
        }
        .preview {
            margin-top: 30px;
            padding: 15px;
            border: 2px solid #ddd;
            border-radius: 4px;
            background-color: #fffacd;
        }
        .preview-title {
            font-weight: bold;
            margin-bottom: 10px;
            color: #333;
        }
        .preview-content {
            padding: 10px;
            border-radius: 4px;
            background-color: #ffffe0;
            color: #000;
            word-wrap: break-word;
        }
        .info {
            margin-top: 20px;
            padding: 10px;
            background-color: #e7f3ff;
            border-left: 4px solid #2196F3;
            color: #333;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📢 告知入力フォーム</h1>
        
        <form method="POST">
            <div class="form-group">
                <div class="checkbox-group">
                    <input type="checkbox" id="enabled" name="enabled" {{ if .Enabled }}checked{{ end }}>
                    <label for="enabled">告知を表示する</label>
                </div>
            </div>

            <div class="form-group">
                <label for="message">告知内容 *</label>
                <textarea id="message" name="message" required placeholder="ここに告知内容を入力してください。例: サーバーメンテナンスのため15時より30分間停止します。">{{ .CurrentMessage }}</textarea>
            </div>

            <div class="form-group">
                <label for="textcolor">文字色</label>
                <input type="color" id="textcolor" name="textcolor" value="{{ .CurrentTextColor }}" title="文字色を選択します">
            </div>

            <div class="form-group">
                <label for="bgcolor">背景色</label>
                <input type="color" id="bgcolor" name="bgcolor" value="{{ .CurrentBgColor }}" title="背景色を選択します">
            </div>

            <div class="button-group">
                <button type="submit" class="btn-submit">✓ 保存して反映</button>
                <button type="reset" class="btn-reset">リセット</button>
            </div>
        </form>

        <div class="preview">
            <div class="preview-title">📋 プレビュー:</div>
            <div class="preview-content" id="preview" style="color: #000; background-color: #ffffe0;">
                告知内容がここに表示されます
            </div>
        </div>

        <div class="info">
            <strong>ℹ️ 情報:</strong>
            <ul style="margin: 5px 0; padding-left: 20px;">
                <li>このページはリンクから到達できません。URLを知っている者だけが告知を入力できます。</li>
                <li>入力した告知は全ページの先頭に表示されます。</li>
                <li>サーバー再起動時に告知内容は消去されます。</li>
                <li>認証機能は未実装です。本運用では認証機能の追加を推奨します。</li>
            </ul>
        </div>
    </div>

    <script>
        // リアルタイムプレビュー
        document.getElementById('message').addEventListener('input', updatePreview);
        document.getElementById('textcolor').addEventListener('change', updatePreview);
        document.getElementById('bgcolor').addEventListener('change', updatePreview);
        document.getElementById('enabled').addEventListener('change', updatePreview);

        function updatePreview() {
            const message = document.getElementById('message').value || '告知内容がここに表示されます';
            const textColor = document.getElementById('textcolor').value;
            const bgColor = document.getElementById('bgcolor').value;
            const enabled = document.getElementById('enabled').checked;
            
            const preview = document.getElementById('preview');
            preview.textContent = message;
            preview.style.color = textColor;
            preview.style.backgroundColor = bgColor;
            preview.style.display = enabled ? 'block' : 'none';
            
            if (!enabled) {
                preview.textContent = '（告知は非表示です）';
                preview.style.color = '#999';
                preview.style.backgroundColor = '#f0f0f0';
            }
        }

        // 初期プレビュー更新
        updatePreview();
    </script>
</body>
</html>
{{ end }}
