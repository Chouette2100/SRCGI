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

    <p>
    <br>
    (外部リンク)<br>
        <a href="https://zenn.dev/">Zenn</a> - <a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a></p>
   <div style="text-indent: 2rem;"><a href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t008top">t008:開催中イベント一覧</a></div>
    <p>-------------------------------------------------------------</p>
    <p>配信中ルーム一覧（{{ UnixTimeToYYYYMMDDHHMM $tn }}）</p>
    <form method="get" action="t009top">
        <div style="text-indent: 2rem;"><div>
        <div style="text-indent: 2rem;"><input type="submit" value="配信ルーム一覧をもう一度取得する">
        （<input type="radio" name="category" value="Official" {{ if eq .Category "Official" }}checked{{ end }}>公式　
        <input type="radio" name="category" value="Free" {{ if eq .Category "Free" }}checked{{ end }}>フリー　
        最近配信を始めたルームから <input type="text" name="maxnoroom" value="20" maxlength="2" type="number" size="2"></form> ルーム表示する ）</div>
    </form>
    <br>
    <table>
        <tr style="text-align: center">
            <td>ルームID</td>
            <td>ジャンル</td>
            <td>ルーム名</td>
            <td>開始</td>
            <td>視聴者</td>
        </tr>
        {{ range .Lives }}
        <tr>
            <td>{{ .Room_id }}</td>
            <td>{{ GidToName .Genre_id }}</td>
            <td><a href="https://www.showroom-live.com/r/{{ .Room_url_key }}" target="_blank" rel="noopener noreferrer">{{ .Main_name }}</a></td>
            <td>{{ UnixTimeToHHMM .Started_at }}〜</td>
            <td style="text-align: right">{{ Comma .View_num }}</td>
        </tr>
        {{ end }}
        <tr>
    </table>
    <p>
        {{ .ErrMsg }}
    </p>
</body>

</html>