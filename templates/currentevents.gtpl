<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    {{/*}}
    <style type="text/css">
        th,
        td {
            border: solid 1px;
        }

        table {
            border-collapse: collapse;
        }
    </style>
    {{*/}}
</head>

<body>
    {{ $tn := .TimeNow }}

    <br>
    {{/*}}
    (外部リンク)<br>
    <a href="https://zenn.dev/">Zenn</a> - <a
        href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a></p>
    <div style="text-indent: 2rem;"><a
            href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t009top">t009:配信中ルーム一覧</a></div>
    <p>-------------------------------------------------------------</p>
    {{*/}}
    <table>
        <tr>
            <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
            <td>開催中イベント一覧</td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
        </tr>
    </table>
    <br>
    <br>
    <p>
    {{ if eq .Mode 1 }}
    （獲得ポイントデータ取得中の）イベント一覧　　
    <button type="button" onclick="location.href='currentevents'">すべてのイベントを表示する</button>
    {{ else }}
    開催中イベント一覧　　
    <button type="button" onclick="location.href='currentevents?mode=1'">獲得ポイントデータ取得中のイベントのみ表示する</button>
    {{ end }}
    </p>
    {{/*
    <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }} （{{ UnixTimeToStr $tn }}）</div>
    */}}
    <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }}
    </div>
    <table border="1" style="border-collapse: collapse">
        <tr bgcolor="gainsboro" style="text-align: center">
            <td>イベント名とイベントページへのリンク</td>
            <td>開始日時</td>
            <td>終了日時</td>
            <td>参加ルーム一覧</td>
            <td>表示項目選択画面/<br>データ取得開始登録</td>
            <td>直近獲得<br>ポイント表</td>
            <td>獲得ポイント<br>推移図</td>
            <td>日々の<br>獲得pt</td>
            <td>枠毎の<br>獲得pt</td>

        </tr>
        {{ $i := 0 }}
        {{ range .Eventinflist }}
        {{ if eq $i 1 }}
        <tr bgcolor="gainsboro">
            {{ $i = 0 }}
            {{ else }}
        <tr>
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
                <a href="eventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}">参加ルーム一覧</a>
            </td>
            <td style="text-align: center;">
                {{ if eq .Target 1 }}
                <a href="top?eventid={{ .Event_ID }}">項目選択</a>
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
        <tr>
    </table>
    <p>
        {{ .ErrMsg}}
    </p>
    <br>
    <hr>
    <br>
    {{/*
    {{ template "footer" }}
    */}}

</body>

</html>