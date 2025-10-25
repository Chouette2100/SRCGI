{{/*
 	id=2661   modelname=deepseek-chat   maxtokens=5000   [25-10-13 19:44 (101.7s)] 
*/}}
<!DOCTYPE html>
<html>
<head>
	<title>JSON Viewer with Collapse</title>
	<style>
		.json-container {
			font-family: 'Consolas', 'Monaco', monospace;
			font-size: 14px;
			line-height: 1.4;
			background: #1e1e1e;
			color: #d4d4d4;
			padding: 20px;
			border-radius: 8px;
			overflow-x: auto;
		}
		body {
			background: #2d2d2d;
			color: white;
			padding: 20px;
			font-family: -apple-system, BlinkMacSystemFont, sans-serif;
			margin: 0;
		}
		.container {
			max-width: 900px;
			margin: 0 auto;
		}
		.controls {
			margin: 20px 0;
			display: flex;
			gap: 10px;
		}
		.btn {
			background: #007acc;
			color: white;
			border: none;
			padding: 8px 16px;
			border-radius: 4px;
			cursor: pointer;
			font-size: 14px;
		}
		.btn:hover {
			background: #005a9e;
		}
		.json-item {
			margin-left: 20px;
		}
		.collapsible {
			cursor: pointer;
			position: relative;
			padding-left: 16px;
		}
		.collapsible::before {
			content: '▶';
			position: absolute;
			left: 0;
			transition: transform 0.2s;
			color: #569cd6;
		}
		.collapsible.collapsed::before {
			transform: rotate(90deg);
		}
		.collapsible + .collapsible-content {
			display: block;
		}
		.collapsible.collapsed + .collapsible-content {
			display: none;
		}
		.json-key { color: #9cdcfe; }
		.json-string { color: #ce9178; }
		.json-number { color: #b5cea8; }
		.json-boolean { color: #569cd6; }
		.json-null { color: #569cd6; }
	</style>
</head>
<body>
	<div class="container">
		<h1>JSON Viewer with Collapse</h1>
		
		<div class="controls">
			<button class="btn" onclick="expandAll()">すべて展開</button>
			<button class="btn" onclick="collapseAll()">すべて折りたたむ</button>
			<button class="btn" onclick="copyToClipboard()">コピー</button>
		</div>

		<div id="jsonViewer" class="json-container"></div>
	</div>

	<script>
        {{/*
		const jsonData = {{.}};
        */}}
		const jsonData = {{ .JSONData }};
		
		// JSONをHTMLに変換する関数
		function jsonToHTML(obj, indent = 0) {
			if (typeof obj === 'string') {
				return '<span class="json-string">"' + escapeHTML(obj) + '"</span>';
			}
			if (typeof obj === 'number') {
				return '<span class="json-number">' + obj + '</span>';
			}
			if (typeof obj === 'boolean') {
				return '<span class="json-boolean">' + obj + '</span>';
			}
			if (obj === null) {
				return '<span class="json-null">null</span>';
			}
			
			if (Array.isArray(obj)) {
				if (obj.length === 0) {
					return '[]';
				}
				
				const id = 'array_' + Math.random().toString(36).substr(2, 9);
				let html = '<span class="collapsible" data-id="' + id + '">[</span>';
				html += '<div class="collapsible-content" id="' + id + '">';
				obj.forEach((item, index) => {
					html += '<div class="json-item">';
					html += jsonToHTML(item, indent + 1);
					if (index < obj.length - 1) html += ',';
					html += '</div>';
				});
				html += '</div>';
				html += '<div class="collapsible">]</div>';
                {{/* こちらにするとブロック下側の三角が消えるが全部なくなるわけではない　要検討
                html += '<div>}</div>';
                */}}
				return html;
			}
			
			if (typeof obj === 'object') {
				const keys = Object.keys(obj);
				if (keys.length === 0) {
					return '{}';
				}
				
				const id = 'object_' + Math.random().toString(36).substr(2, 9);
				let html = '<span class="collapsible" data-id="' + id + '">{</span>';
				html += '<div class="collapsible-content" id="' + id + '">';
				keys.forEach((key, index) => {
					html += '<div class="json-item">';
					html += '<span class="json-key">"' + escapeHTML(key) + '"</span>: ';
					html += jsonToHTML(obj[key], indent + 1);
					if (index < keys.length - 1) html += ',';
					html += '</div>';
				});
				html += '</div>';
				html += '<div class="collapsible">}</div>';
				return html;
			}
			
			return '';
		}
		
		function escapeHTML(str) {
			return str.replace(/[&<>"']/g, function(match) {
				const escapeMap = {
					'&': '&amp;',
					'<': '&lt;',
					'>': '&gt;',
					'"': '&quot;',
					"'": '&#39;'
				};
				return escapeMap[match];
			});
		}
		
		// 初期表示
		function renderJSON() {
			try {
				const parsedJSON = JSON.parse(jsonData);
				document.getElementById('jsonViewer').innerHTML = jsonToHTML(parsedJSON);
				addCollapseListeners();
			} catch (e) {
				document.getElementById('jsonViewer').textContent = 'Error parsing JSON: ' + e.message;
			}
		}
		
		// 折りたたみ機能のイベントリスナーを追加
		function addCollapseListeners() {
			document.querySelectorAll('.collapsible').forEach(element => {
				element.addEventListener('click', function() {
					const contentId = this.getAttribute('data-id');
					if (contentId) {
						this.classList.toggle('collapsed');
					}
				});
			});
		}
		
		// コントロール関数
		function expandAll() {
			document.querySelectorAll('.collapsible').forEach(el => {
				el.classList.remove('collapsed');
			});
		}
		
		function collapseAll() {
			document.querySelectorAll('.collapsible').forEach(el => {
				if (el.getAttribute('data-id')) {
					el.classList.add('collapsed');
				}
			});
		}
		
		function copyToClipboard() {
			navigator.clipboard.writeText(jsonData).then(() => {
				alert('JSONをクリップボードにコピーしました');
			}).catch(err => {
				console.error('Copy failed: ', err);
			});
		}
		
		// 初期化
		renderJSON();
	</script>
</body>
</html>