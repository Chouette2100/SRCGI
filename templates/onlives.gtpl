<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <title>配信履歴 roomid={{.RoomID}}</title>
    <style>
        :root {
            --bg: #f6f8fb;
            --paper: #ffffff;
            --line: #d7dbe3;
            --line-bold: #b8bfcc;
            --text: #1e293b;
            --sub: #475569;
            --accent: #007a6e;
            --accent-2: #b7f3ea;
        }
        * { box-sizing: border-box; }
        body {
            margin: 0;
            padding: 24px;
            font-family: "Noto Sans JP", "Hiragino Kaku Gothic ProN", Meiryo, sans-serif;
            color: var(--text);
            background: radial-gradient(circle at 15% 15%, #ffffff 0%, #eef2f8 45%, #e5ebf5 100%);
        }
        .panel {
            max-width: 1280px;
            margin: 0 auto;
            background: var(--paper);
            border: 1px solid #e2e8f0;
            border-radius: 12px;
            padding: 18px;
            box-shadow: 0 14px 32px rgba(15, 23, 42, 0.08);
        }
        h1 {
            margin: 0 0 8px;
            font-size: 24px;
            letter-spacing: 0.02em;
        }
        .meta {
            margin: 0 0 16px;
            color: var(--sub);
            font-size: 13px;
        }
        .order-chip {
            display: inline-block;
            margin-left: 8px;
            padding: 2px 8px;
            border-radius: 999px;
            border: 1px solid #cbd5e1;
            background: #f8fafc;
            color: #334155;
            font-size: 12px;
        }
        .meta a {
            margin-left: 8px;
            color: #0f766e;
            text-decoration: none;
            border-bottom: 1px dotted #0f766e;
        }
        .metrics {
            display: flex;
            flex-wrap: wrap;
            gap: 12px;
            margin-bottom: 12px;
        }
        .metric {
            min-width: 200px;
            padding: 10px 12px;
            border-radius: 10px;
            border: 1px solid #dbe4ee;
            background: linear-gradient(180deg, #ffffff, #f8fbff);
        }
        .metric .label {
            font-size: 12px;
            color: var(--sub);
            margin-bottom: 4px;
        }
        .metric .value {
            font-size: 22px;
            font-weight: 700;
        }
        .note {
            margin: 8px 0 14px;
            font-size: 12px;
            color: var(--sub);
        }
        .now-line {
            stroke: #dc2626;
            stroke-width: 1.5;
            stroke-dasharray: 5 3;
        }
        .now-label {
            fill: #dc2626;
            font-size: 11px;
            font-weight: 700;
        }
        .graph-wrap {
            overflow-x: auto;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            padding: 10px;
            background: #fbfdff;
        }
        .event-bar {
            cursor: pointer;
            transition: fill-opacity 0.15s ease, stroke-width 0.15s ease;
        }
        .event-bar:hover,
        .event-bar:focus {
            fill-opacity: 1;
            stroke-width: 1.3;
            outline: none;
        }
        .event-bar.is-selected {
            fill: #00995f;
            fill-opacity: 1;
            stroke: #074d35;
            stroke-width: 1.6;
        }
        svg {
            width: 100%;
            min-width: 980px;
            height: auto;
            display: block;
        }
        .detail {
            margin-top: 10px;
            border: 1px solid #d7e2ea;
            border-radius: 8px;
            padding: 10px 12px;
            background: #f8fcff;
            font-size: 13px;
        }
        .detail-title {
            margin: 0 0 6px;
            font-size: 12px;
            color: #3f5469;
        }
        .detail-grid {
            display: grid;
            grid-template-columns: 88px 1fr;
            gap: 4px 12px;
        }
        .detail dt {
            color: #526579;
        }
        .detail dd {
            margin: 0;
            font-weight: 600;
            letter-spacing: 0.01em;
        }
        .day-table {
            margin-top: 16px;
            border-collapse: collapse;
            width: 100%;
            font-size: 13px;
        }
        .day-table th,
        .day-table td {
            border: 1px solid #dfe6ee;
            padding: 6px 8px;
            text-align: right;
        }
        .day-table th:first-child,
        .day-table td:first-child {
            text-align: left;
        }
        .day-table th {
            background: #f4f7fb;
        }
    </style>
</head>
<body>
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}

    <p>
    <button type="button" onclick="location.href='top'">Top</button>　
    <button type="button" onclick="location.href='currentevents'">開催中イベント一覧表</button>　
    <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧表</button>　
    <button type="button" onclick="location.href='closedevents'">終了イベント一覧表</button>　
    </p>
    <div class="panel">
        <h1>配信タイムライン roomid={{.RoomID}}</h1>
        <p class="meta">期間: {{.WindowStart.Format "2006-01-02"}} から {{.WindowEnd.AddDate 0 0 -1 | date "2006-01-02"}} (JST)
            <span class="order-chip">表示順: {{.SortOrderLabel}}</span>
            <a href="/onlives?roomid={{.RoomID}}&days={{.Days}}&order={{.ToggleSortOrder}}">{{.ToggleSortOrderLabel}}</a>
        </p>

        <div class="metrics">
            <div class="metric">
                <div class="label">配信回数 (期間内)</div>
                <div class="value">{{.TotalCount}} 回</div>
            </div>
            <div class="metric">
                <div class="label">平均配信時間</div>
                <div class="value">{{.AverageDurationMin}} 分</div>
            </div>
            {{/*
            <div class="metric">
                <div class="label">開始時刻ばらつき (標準偏差)</div>
                <div class="value">{{printf "%.2f" .StartHourStdDev}} h</div>
            </div>
            <div class="metric">
                <div class="label">配信スケジュール傾向</div>
                <div class="value">{{.ScheduleStabilityMsg}}</div>
            </div>
            */}}
        </div>

        <p class="note">配信バーはホバーで元の配信時刻を表示し、クリックで下の詳細パネルに固定表示します。赤い縦線は現在時刻です。</p>

        <div class="graph-wrap">
            <svg viewBox="0 0 {{.SVGWidth}} {{.SVGHeight}}" role="img" aria-label="配信タイムライン">
                <rect x="0" y="0" width="{{.SVGWidth}}" height="{{.SVGHeight}}" fill="#ffffff"/>

                {{range .HourTicks}}
                <line x1="{{printf "%.2f" .X}}" y1="{{$.TopMargin}}" x2="{{printf "%.2f" .X}}" y2="{{printf "%.2f" $.BottomLineY}}" stroke="{{if .Major}}#b8bfcc{{else}}#d7dbe3{{end}}" stroke-width="1"/>
                {{if .ShowLabel}}
                <text x="{{printf "%.2f" .X}}" y="24" font-size="11" fill="#475569" text-anchor="middle">{{.Label}}</text>
                {{end}}
                {{end}}

                {{range .Rows}}
                <line x1="{{$.LeftMargin}}" y1="{{printf "%.2f" .Y}}" x2="{{printf "%.2f" $.PlotRightX}}" y2="{{printf "%.2f" .Y}}" stroke="#d7dbe3" stroke-width="1"/>
                <text x="{{sub $.LeftMargin 8}}" y="{{printf "%.2f" .TextY}}" font-size="12" fill="#334155" text-anchor="end">{{.DateLabel}}</text>
                {{end}}
                <line x1="{{$.LeftMargin}}" y1="{{printf "%.2f" $.BottomLineY}}" x2="{{printf "%.2f" $.PlotRightX}}" y2="{{printf "%.2f" $.BottomLineY}}" stroke="#d7dbe3" stroke-width="1"/>

                {{if .HasNowLine}}
                <line class="now-line" x1="{{printf "%.2f" .NowLineX}}" y1="{{printf "%.2f" .NowLineY1}}" x2="{{printf "%.2f" .NowLineX}}" y2="{{printf "%.2f" .NowLineY2}}"/>
                <text class="now-label" x="{{printf "%.2f" .NowLineX}}" y="{{printf "%.2f" .NowLabelY}}" text-anchor="end" dx="-3">現在</text>
                <text class="now-label" x="{{printf "%.2f" .NowLineX}}" y="{{printf "%.2f" .NowLabelY}}" text-anchor="start" dx="3">{{.NowLabel}}</text>
                {{end}}

                {{range .Bars}}
                <rect class="event-bar" tabindex="0" data-id="{{.ID}}" data-start="{{.Start}}" data-end="{{.End}}" data-duration="{{.DurText}}" data-segment="{{.SegText}}" x="{{printf "%.2f" .X}}" y="{{printf "%.2f" .Y}}" width="{{printf "%.2f" .Width}}" height="{{printf "%.2f" .Height}}" fill="var(--accent)" fill-opacity="0.88" stroke="#0f5e58" stroke-width="0.6" rx="3" ry="3">
                    <title>{{.Tooltip}}</title>
                </rect>
                {{end}}
            </svg>
        </div>

        <section class="detail" id="bar-detail" aria-live="polite">
            <p class="detail-title">クリックした配信バーの詳細</p>
            <dl class="detail-grid">
                <dt>開始</dt><dd id="detail-start">未選択</dd>
                <dt>終了</dt><dd id="detail-end">未選択</dd>
                <dt>長さ</dt><dd id="detail-duration">未選択</dd>
                <dt>表示区間</dt><dd id="detail-segment">未選択</dd>
            </dl>
        </section>

        <table class="day-table">
            <thead>
                <tr>
                    <th>日付</th>
                    <th>配信回数</th>
                    <th>合計配信時間(分)</th>
                    <th>平均配信時間(分)</th>
                </tr>
            </thead>
            <tbody>
                {{range .DayStats}}
                <tr>
                    <td>{{.DateLabel}}</td>
                    <td>{{.Count}}</td>
                    <td>{{.TotalMinutes}}</td>
                    <td>{{.AverageMinutes}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <script>
    (function () {
        const bars = Array.from(document.querySelectorAll('.event-bar'));
        const startEl = document.getElementById('detail-start');
        const endEl = document.getElementById('detail-end');
        const durEl = document.getElementById('detail-duration');
        const segEl = document.getElementById('detail-segment');
        let selectedId = null;

        function clearSelection() {
            bars.forEach((el) => el.classList.remove('is-selected'));
            selectedId = null;
            startEl.textContent = '未選択';
            endEl.textContent = '未選択';
            durEl.textContent = '未選択';
            segEl.textContent = '未選択';
        }

        function selectBar(el) {
            const id = el.dataset.id;
            if (selectedId === id) {
                clearSelection();
                return;
            }
            bars.forEach((node) => node.classList.toggle('is-selected', node === el));
            selectedId = id;
            startEl.textContent = el.dataset.start;
            endEl.textContent = el.dataset.end;
            durEl.textContent = el.dataset.duration;
            segEl.textContent = el.dataset.segment;
        }

        bars.forEach((bar) => {
            bar.addEventListener('click', () => selectBar(bar));
            bar.addEventListener('keydown', (ev) => {
                if (ev.key === 'Enter' || ev.key === ' ') {
                    ev.preventDefault();
                    selectBar(bar);
                }
            });
        });
    })();
    </script>
</body>
</html>
