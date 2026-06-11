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
        .no-border-left {
          border-left-style: none;
        }
        .no-border-right {
          border-right-style: none;
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
    <p>月初めのSHOWランク上位のルーム</p>
    <p style="color: red;">
    こちらは月始めのデータです。<br>
    こういういい加減にも思えるデータを公開するのは（もしそういうのがあれば）このサイトの信用を傷つけることになりそうです<br>
    ただ一人で悩んでいてもこれ以上の進展はないと思うのであえて公開する次第です。<br>
    「月初のSHOWランク」というのは1日の0時30分から15分ほどかけて取得したデータをいいます<br>
    （「SHOWランク上位のルーム」の表の「DBに登録済みの配信者のうちSHOWランク上位220ルーム」に相当します<br>
    　これらはすべて同一のタイムスタンプ（＝00:30）になります<br>
    この表を作ったのは月初のSHOWランクが本来のというか現実的な意味でのSHOWランクになると思うからです<BR>
    ただ今月(=2026年6月)のデータを見るとなにかの間違いじゃないのかと思ってしまうようなデータが散見されます<br>
    例えば
    A. ランクごとの人数が予想しているルールとまったく違う、とくにA-3ランクは1ルームもない、総ルーム数も異常に少ない<br>
    B. next_scoreとprev_scoreにキリのいい値とそうでないものが混じっている<br>
    C. ランク間のnext_scoreとprev_scoreの整合性がとれていない<br>
    これについては<br>
    1. データの取得タイミングが早すぎる？<br>
    2. データの取得タイミングが遅すぎる？<br>
    3. SHOWランクの意味を誤解している？<br>
    4. そもそもデータの取得方法が間違っている？<br>
    5. 1日未明に配信しているルームがあった場合、どう扱われるのかわからない。<br>
    いろいろ考えられるのですが、この分野に詳しい方のご意見がいただけるかもしれないと思いそのまま公開することにしました。
    </p>
    <br>各ランクの「定員」は現時点では以下のようになっていると思われます。
    <br>
    <br>SSランク 10名（SS-5 1名、SS-4〜SS-2 各2名、SS-1 3名）
    <br>Sランク　30名（S-5〜S-1 各6名）
    <br>Aランク　100名（A-5〜A-2 各12名、A-1 52名）
    <br>
    「<a href="https://debug36.chouette2100.com/showrank">SHOWランクが上位のルーム</a>」
    <br>
    <br>
    {{ $l := "SS-5" }}
    {{ $n := "SS-5" }}
    {{ $tmn := "SS-5" }}
    {{ $tmts := "0001-01-01 00:00" }}
    {{ $c := "lightyellow" }}
    {{ $i := 1 }}
    <table>
        <tr style="text-align: center">
            <td>No.</td>
            <td>SHOWランク</td>
            <td class="no-border-right">ルーム名（イベント履歴へのリンク）</td>
            <td class="no-border-left"></td>
            <td>ルームレベル<br>(月初)</td>
            <td>フォロワー数<br>(月初)</td>
            <td>ファン数<br>(月初)</td>
            {{/*
            <td>ファンパワー<br>（〃）</td>
            */}}
            <td>ファン数<br>（前月）</td>
            {{/*
            <td>ファンパワー<br>（〃）</td>
            */}}
            <td>next_score</td>
            <td>prev_score</td>
            <td>データ取得日時</td>
            <td>SHOWランク<br>(最新)</td>
            <td>next_score<br>(最新)</td>
            <td>prev_score<br>(最新)</td>
            <td>データ取得日時<br>(最新)</td>
        </tr>
        {{ range .Userlist }}
        {{ $n = Showrank .Rank }}
        {{ $tmn = Showrank .Tmrank }}
        {{ $tmts = FormatTime .Tmts "2006-01-02 15:04" }}
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
            <td class="no-border-right">
            <a href="closedevents?userno={{ .Userno }}&mode=0&path=5">{{ .User_name }}</a>
            </td>
            <td class="no-border-left">
            <a href="https://www.showroom-live.com/room/profile?room_id={{ .Userno }}">prof.</a>
            </td>
            <td align="right">{{ Comma .Level }}</td>
            <td align="right">{{ Comma .Followers }}</td>
            <td align="right">{{ Comma .Fans }}</td>
            {{/*
            {{ if gt .FanPower 0 }}
            <td align="right">{{ Comma .FanPower }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            */}}
            <td align="right">{{ Comma .Fans_lst }}</td>
            {{/*
            {{ if gt .FanPower_lst 0 }}
            <td align="right">{{ Comma .FanPower_lst }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            */}}
            <td align="right">{{ Comma .Inrank }}</td>
            <td align="right">{{ Comma .Iprank }}</td>
            <td>{{ FormatTime .Ts "2006-01-02 15:04" }}</td>
            <td align="right">{{ $tmn }}</td>
            {{ if gt .Tminrank 0 }}
            <td align="right">{{ Comma .Tminrank }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            {{ if gt .Tmiprank 0 }}
            <td align="right">{{ Comma .Tmiprank }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            {{ if ne $tmts "0001-01-01 00:00" }}
            <td>{{ $tmts }}</td>
            {{ else }}
            <td></td>
            {{ end }}
        </tr>
        {{ $i = Add $i 1 }}
        {{ end }}
        <tr style="text-align: center">
            <td>No.</td>
            <td>SHOWランク</td>
            <td class="no-border-right">ルーム名（イベント履歴へのリンク）</td>
            <td class="no-border-left"></td>
            <td>ルームレベル<br>(月初)</td>
            <td>フォロワー数<br>(月初)</td>
            <td>ファン数<br>(月初)</td>
            {{/*
            <td>ファンパワー<br>（〃）</td>
            */}}
            <td>ファン数<br>（前月）</td>
            {{/*
            <td>ファンパワー<br>（〃）</td>
            */}}
            <td>next_score</td>
            <td>prev_score</td>
            <td>データ取得日時</td>
            <td>SHOWランク<br>(最新)</td>
            <td>next_score<br>(最新)</td>
            <td>prev_score<br>(最新)</td>
            <td>データ取得日時<br>(最新)</td>
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
            <td class="no-border-right"><a href="closedevents?userno={{ .Userno }}&mode=0&path=5">{{ .User_name }}</a></td>
            <td class="no-border-left"></td>
            <td align="right">{{ Comma .Level }}</td>
            <td align="right">{{ Comma .Followers }}</td>
            <td align="right">{{ Comma .Fans }}</td>
            {{/*
            {{ if gt .FanPower 0 }}
            <td align="right">{{ Comma .FanPower_lst }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            */}}
            <td align="right">{{ Comma .Fans_lst }}</td>
            {{/*
            {{ if gt .FanPower_lst 0 }}
            <td align="right">{{ Comma .FanPower_lst }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            */}}
            {{ if gt .Tminrank 0 }}
            <td align="right">{{ Comma .Tminrank }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            {{ if gt .Tmiprank 0 }}
            <td align="right">{{ Comma .Tmiprank }}</td>
            {{ else }}
            <td></td>
            {{ end }}
            {{ if ne .Tmrank "0001-01-01 00:00" }}
            <td>{{ FormatTime .Ts "2006-01-02 15:04" }}</td>
            {{ else }}
            <td></td>
            {{ end }}
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