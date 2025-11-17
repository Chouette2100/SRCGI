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
    <p>
    {{/*}}
    (外部リンク)<br>
        <a href="https://zenn.dev/">Zenn</a> - <a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a></p>
   <div style="text-indent: 2rem;"><a href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t008top">t008:開催中イベント一覧</a></div>
    <p>-------------------------------------------------------------</p>
    {{*/}}
    <button type="button" onclick="location.href='top'">Top</button>　
    <button type="button" onclick="location.href='currentevents'">開催中イベント一覧表</button>　
    <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧表</button>　
    <button type="button" onclick="location.href='closedevents'">終了イベント一覧表</button>　
    <br>
    <br>
    <p>SHOWランク上位のルーム</p>
    昨日〜今日あたりの特定の時点でのSHOWランク上位ルームの一覧です。
    <br>各ランクの「定員」は現時点では以下のようになっていると思われます。
    <br>
    <br>SSランク 10名（SS-5 1名、SS-4〜SS-2 各2名、SS-1 3名）
    <br>Sランク　30名（S-5〜S-1 各6名）
    <br>Aランク　100名（A-5〜A-2 各12名、A-1 52名）
    <br>
    <br>いろんな方法を組み合わせてデータを集めているのですが、以上の全員を確実に把握できるとは限らないので、
    <br>以下の表に表示されているルームが140に満たないこともありえます。
    <br>もし、○○ルームが抜けている、とわかる場合は掲示板等で教えていただけると助かります(データ取得方法の改善ができるかもしれません)
    <br>
    <br>
    {{ $l := "SS-5" }}
    {{ $n := "SS-5" }}
    {{ $c := "lightyellow" }}
    {{ $i := 1 }}
    <table>
        <tr style="text-align: center">
            <td>No.</td>
            <td>SHOWランク</td>
            <td>ルーム名（イベント履歴へのリンク）</td>
            <td>ルームレベル</td>
            <td>フォロワー数</td>
            <td>ファン数<br>（今月）</td>
            <td>ファンパワー<br>（〃）</td>
            <td>ファン数<br>（前月）</td>
            <td>ファンパワー<br>（〃）</td>
            <td>next_score</td>
            <td>prev_score</td>
            <td>データ取得日時</td>
        </tr>
        {{ range .Userlist }}
        {{ $n = Showrank .Rank }}
        {{ if ne $n $l }}
            {{ $l = $n }}
            {{ if eq $c "lightyellow" }}
                {{ $c = "lightcyan" }}
            {{ else }}
                {{ $c = "lightyellow" }}
            {{ end }}
        {{ end }}
        <tr bgcolor="{{$c}}" >
            <td align="right">{{ $i }}</td>
            <td align="center">{{ $n }}</td>
            <td><a href="closedevents?userno={{ .Userno }}&mode=0&path=5">{{ .User_name }}</a></td>
            <td align="right">{{ Comma .Level }}</td>
            <td align="right">{{ Comma .Followers }}</td>
            <td align="right">{{ Comma .Fans }}</td>
            {{ if gt .FanPower 0 }}
            <td align="right">{{ Comma .FanPower }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            <td align="right">{{ Comma .Fans_lst }}</td>
            {{ if gt .FanPower_lst 0 }}
            <td align="right">{{ Comma .FanPower_lst }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            <td align="right">{{ Comma .Inrank }}</td>
            <td align="right">{{ Comma .Iprank }}</td>
            <td>{{ FormatTime .Ts "2006-01-02 15:04" }}</td>
        </tr>
        {{ $i = Add $i 1 }}
        {{ end }}

        <tr style="text-align: center">
            <td></td>
            <td>SHOWランク</td>
            <td>ルーム名（イベント履歴へのリンク）</td>
            <td>ルームレベル</td>
            <td>フォロワー数</td>
            <td>ファン数<br>（今月）</td>
            <td>ファンパワー<br>（〃）</td>
            <td>ファン数<br>（前月）</td>
            <td>ファンパワー<br>（〃）</td>
            <td>next_score</td>
            <td>prev_score</td>
            <td>データ取得日時</td>
        </tr>

        {{ range .UserlistA }}
        {{ $n = Showrank .Rank }}
        {{ if ne $n $l }}
            {{ $l = $n }}
            {{ if eq $c "lightyellow" }}
                {{ $c = "lightcyan" }}
            {{ else }}
                {{ $c = "lightyellow" }}
            {{ end }}
        {{ end }}
        <tr bgcolor="{{$c}}" >
            <td align="right"></td>
            <td align="center">{{ $n }}</td>
            <td><a href="closedevents?userno={{ .Userno }}&mode=0&path=5">{{ .User_name }}</a></td>
            <td align="right">{{ Comma .Level }}</td>
            <td align="right">{{ Comma .Followers }}</td>
            <td align="right">{{ Comma .Fans }}</td>
            {{ if gt .FanPower 0 }}
            <td align="right">{{ Comma .FanPower_lst }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            <td align="right">{{ Comma .Fans_lst }}</td>
            {{ if gt .FanPower_lst 0 }}
            <td align="right">{{ Comma .FanPower_lst }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            <td align="right">{{ Comma .Inrank }}</td>
            <td align="right">{{ Comma .Iprank }}</td>
            <td>{{ FormatTime .Ts "2006-01-02 15:04" }}</td>
        </tr>
        {{ end }}
    </table>
    {{/*
    <p>
        {{ .ErrMsg }}
    </p>
    */}}
</body>

</html>