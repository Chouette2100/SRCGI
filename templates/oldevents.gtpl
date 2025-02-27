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

        .noborder {
            border: 0px none;
        }
    </style>
    {{*/}}
</head>

<body>

    <br>
    {{/*
    (外部リンク)<br>
    <a href="https://zenn.dev/">Zenn</a> - <a
        href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a>
    <div style="text-indent: 2rem;"><a
            href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t009top">t009:配信中ルーム一覧</a></div>
    <p>-------------------------------------------------------------</p>
    */}}

    
    <table>
        <tr>
            <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
            <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td>終了イベント一覧</td>
        </tr>
    </table>
    <br>
    {{/*
        <div style="text-indent: 2rem;">
            一覧は51件ずつ表示され、50件ずつスクロールされます。データ上最終結果が存在しても表示されないケースがあります（修正・改良予定あり）
        </div>
        
        {{ if ne .Offset 0 }}

            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode={{ .Mode }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&kwevid={{ .Kwevid }}&userno={{ .Userno }}&action=top&limit={{ .Limit }}&offset={{.Offset}}'">最初から表示する</button>

            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode={{ .Mode }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&kwevid={{ .Kwevid }}&userno={{ .Userno }}&action=prev&limit={{ .Limit }}&offset={{.Offset}}'">前ページ</button>


        {{ end }}
    */}}

        {{ $u := .User.Userno }}
        <br>{{ .User.User_name }} （ {{ .User.Userno }} ）
    <br><br>


        <table border="1" style="border-collapse: collapse">
            <tr bgcolor="gainsboro" style="text-align: center">
                <td>イベント名（とイベントページへのリンク）</td>
                <td>開始日時</td>
                <td>終了日時</td>
                <td>最終結果<br>(30位まで)</td>
                <td>表示項目<br>選択画面</td>
                <td>最終獲得<br>ポイント表</td>
                <td>獲得ポイント<br>推移図</td>
                <td>日々の<br>獲得pt</td>
                <td>枠毎の<br>獲得pt</td>
                <td>貢献<br>pt</td>
            </tr>

            {{ $i := 0 }}
            {{ range .Wevents }}
            {{ if eq $i 1 }}
            <tr bgcolor="gainsboro">
                {{ $i = 0 }}
                {{ else }}
            <tr>
                {{ $i = 1 }}
                {{ end }}

                <td>
                    {{ if IsTempID .Eventid }}
                        {{ .Event_name }} ( {{ .Ieventid }} )
                    {{ else }}
                        <a href="https://showroom-live.com/event/{{ .Eventid }}">{{ .Event_name }}</a>
                    {{ end }}
                </td>
                <td>
                    {{ TimeToStringY .Starttime }}
                </td>
                <td>
                    {{ TimeToString .Endtime }}
                </td>
                <td>{{/* 最終結果<br>(30位まで) */}}</td>
                <td>{{/* 表示項目<br>選択画面 */}}</td>
                <td>{{/* 最終獲得<br>ポイント表 */}}</td>
                <td>{{/* 獲得ポイント<br>推移図< */}}</td>
                <td>{{/* 日々の<br>獲得pt< */}}</td>
                <td>{{/* 枠毎の<br>獲得pt< */}}</td>
                {{/* <td>貢献<br>pt</td> */}}
                <td style="text-align: center;">
                    <a href="contributors?ieventid={{.Ieventid}}&roomid={{$u}}">CSV</a>
                </td>
            </tr>
            {{ end }}
            <tr>
        </table>

{{/*
        {{ if eq .Totalcount .Limit }}
           <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode={{ .Mode }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&kwevid={{ .Kwevid }}&userno={{ .Userno }}&action=next&limit={{ .Limit }}&offset={{.Offset}}'">次ページ</button>
        {{ end }}

        {{ end }}
*/}}
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