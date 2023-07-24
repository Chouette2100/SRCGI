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
    <button type="button" onclick="location.href='top'">Top</button>　
    <button type="button" onclick="location.href='currentevent'">開催中イベント一覧表</button>　
    <button type="button" onclick="location.href='scheduledevent'">開催予定イベント一覧表</button>　
    {{/*}}
    <button type="button" onclick="location.href='closedevent'">終了イベント一覧表</button>　
    {{*/}}
    <br>
    <br>
    イベント名に含まれる文字列で絞り込む
    <p>
    <form>
        <table>
            <tr>
                <td style="width:4em"></td>
                <td>イベント名に含まれる文字列</td>
                <td><input type="text" value="{{ .Keywordev }}" name="keywordev"></td>
                <td align="right">
                    <input type="submit" value="イベント名で絞り込む" formaction="closedevents" formmethod="GET" {{/*
                        style="background-color: aquamarine" */}}>
                </td>
            </tr>
            <tr>
                <td style="width:4em"></td>
                <td colspan="2">
                    （「含まれる文字列」は例えば、スタートダッシュ、花火、Music 等で、
                    <br>カタカナとひらがな、全角と半角、英大文字と小文字あたりはアバウトです。）
                </td>
                <td></td>
            </tr>
        </table>
        </p>
        <input type="hidden" name="mode" value="{{ .Mode }}" />
    </form>
    <br>
    <br>
    エントリーしているルーム名から絞り込む
    <form>
        <table>
            <tr>
                <td style="width:4em"></td>
                <td>ルーム名に含まれる文字列</td>
                <td><input type="text" value="{{ .Keywordrm }}" name="keywordrm"></td>
                <td align="right">
                    <input type="submit" value="ルームを検索する" formaction="closedevents" formmethod="GET" {{/*
                        style="background-color: aquamarine" */}}>
                </td>
            </tr>
            <tr>
                <td style="width:4em"></td>
                <td colspan="2">
                    （現在のルーム名だけでなく過去のルーム名（のうち最近のもの）も<br>検索対象となります。ただし検索結果は50件で打ち切り）
                </td>
                <td></td>
            </tr>
        </table>
        </p>
        <input type="hidden" name="mode" value="{{ .Mode }}" />
    </form>

    <form>
        <table>
            <tr>
                <td style="width:4em"></td>
                <td>ルームを選択する</td>
                <td>

                    <input name="userno" type="text" list="combolist" size="40">
                    <datalist id="combolist">
                        {{ range .Roomlist }}
                        <option value="{{ .Userno }}">{{ .User_name }}</option>
                        {{ end }}
                    </datalist>
                </td>
                <td align="right">
                    <input type="submit" value="ルームで絞り込む" formaction="closedevents" formmethod="GET" {{/*
                        style="background-color: aquamarine" */}}>
                </td>
            </tr>
            <tr>
                <td style="width:4em"></td>
                <td colspan="2">
                    リストが表示されないときは二回クリックお願いします。<br>usernoがわかっていたら直接入力可！
                </td>
                <td></td>
            </tr>

        </table>
        </p>
        <input type="hidden" name="mode" value="{{ .Mode }}" />
    </form>

    <br>
    <br>
    <p>
        {{ if eq .Mode 1 }}
        （獲得ポイント詳細データがある）終了イベント一覧　　
        {{ if eq .Keywordev "" }}
        <button type="button" onclick="location.href='closedevents'">すべてのイベントを表示する</button>
        {{ else }}
        <button type="button" onclick="location.href='closedevents?keyword={{ .Keywordev }}'">すべてのイベントを表示する</button>
        {{ end }}
        {{ else }}
        終了イベント一覧　　
        {{ if eq .Keywordev "" }}
        <button type="button" onclick="location.href='closedevents?mode=1'">獲得ポイント詳細データのある終了イベントのみ表示する</button>
        {{ else }}
        <button type="button"
            onclick="location.href='closedevents?mode=1&keyword={{ .Keywordev }}'">獲得ポイント詳細データのある終了イベントのみ表示する</button>
        {{ end }}
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
            <td></td>
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
                <a href="closedeventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}">最終結果</a>
            </td>
            <td style="text-align: center;">
                {{ if eq .Target 1 }}
                <a href="top?eventid={{ .Event_ID }}">結果表示選択</a>
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