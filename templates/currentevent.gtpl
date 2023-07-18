<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <style type="text/css">
        th,
        td {
            border: solid 1px;
        }

        table {
            border-collapse: collapse;
            /*
            width: 100%;
            */
        }
    </style>
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
    <button type="button" onclick="location.href='top'">イベント選択画面に戻る</button>　
    <button type="button" onclick="location.href='currentevent'">開催中イベント一覧表</button>　
    <br>
    <br>
    <p>開催中イベント一覧表</p>
    {{/*
    <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }} （{{ UnixTimeToStr $tn }}）</div>
    */}}
    <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }} </div>
    <table>
        <tr bgcolor="gainsboro" style="text-align: center">
            <td>イベント名とイベントページへのリンク</td>
            <td>開始日時</td>
            <td>終了日時</td>
            <td>参加ルーム一覧</td>
            <td>結果表示選択画面/<br>データ取得新規登録</td>
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
                <a href="eventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}">表示</a>
            </td>
            <td style="text-align: center;">
                {{ if eq .Target 1 }}
                <a href="top?eventid={{ .Event_ID }}">表示</a>
                {{ else }}
                <a href="new-event?eventid={{ .Event_ID }}">新規登録</a>
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