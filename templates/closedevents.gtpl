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
    <button type="button" onclick="location.href='currentevents'">開催中イベント一覧表</button>　
    <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧表</button>　
    {{/*}}
    <button type="button" onclick="location.href='closedevents'">終了イベント一覧表</button>　
    {{*/}}
    <form> <!-- 0. 最初のパス -->
        <br>
        <br>
        イベント名に含まれる文字列で絞り込む
        <p>
        <form> <!-- 1. イベント名で絞り込む -->
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
                        （「含まれる文字列」は例えば、スタートダッシュ、花火、Music 等です。
                        <br>カタカナとひらがな、全角と半角、英大文字と小文字あたりはアバウトです。）
                    </td>
                    <td></td>
                </tr>
            </table>
            </p>
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="1" />
        </form>
        <!-- 2. イベントID(Event_url_key)で絞り込む -->
        <br>
        <br>
        エントリーしているルーム名から絞り込む
        <form> <!-- 3. ルーム名で絞り込む(ルーム名の入力) -->
            <table>
                <tr>
                    <td style="width:4em"></td>
                    <td>ルーム名に含まれる文字列</td>
                    <td><input type="text" value="{{ .Keywordrm }}" name="keywordrm"></td>
                    <td align="right">
                        <input type="submit" value="ルームを検索する" formaction="closedevents" formmethod="GET">
                    </td>
                </tr>
                <tr>
                    <td style="width:4em"></td>
                    <td colspan="2">
                        （現在のルーム名だけでなく過去のルーム名（のうち最近のもの、例えば「夜風」）も
                        <br>検索対象となります。ただしルームの検索結果は30件までしか表示されません）
                    </td>
                    <td></td>
                </tr>
            </table>
            </p>
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="3" />
        </form>

        <form> <!-- 4. ルーム名で絞り込む(ルーム名の選択) -->
            <table>
                <tr>
                    <td style="width:4em"></td>
                    <td>ルームを選択する</td>
                    <td>
                        {{ $userno :=  .Userno }}
                        <select name="userno" type="text">
                            {{ range .Roomlist }}
                            {{ if eq .Userno $userno }}
                            <option selected value="{{ .Userno }}">{{ .User_name }}</option>
                            {{ else }}
                            <option value="{{ .Userno }}">{{ .User_name }}</option>
                            {{ end }}
                            {{ end }}
                        </select>
                    </td>
                    <td align="right">
                        <input type="submit" value="ルームで絞り込む" formaction="closedevents" formmethod="GET">
                    </td>
                </tr>
                <tr>
                    <td style="width:4em"></td>
                    <td colspan="2">
                        {{/*}}
                        リストが表示されないときは二回クリックお願いします。<br>usernoがわかっていたら直接入力可！
                        {{*/}}
                    </td>
                    <td></td>
                </tr>

            </table>
            </p>
            <input type="hidden" name="keywordrm" value="{{ .Keywordrm }}" />
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="4" />
        </form>

        <!-- 5. ユーザ番号で選択する -->

        <br>
        <br>
        {{ if or ( eq .Keywordrm "" ) ( ne .Userno 0) }}
        <p>
            {{ if eq .Mode 1 }}
            （獲得ポイント詳細データがある）終了イベント一覧　　
            {{ if eq .Keywordev "" }}
            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&keywordrm={{ .Keywordrm }}&userno={{ .Userno }}'">すべてのイベントを表示する</button>
            {{ else }}
            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&userno={{ .Userno }}'">すべてのイベントを表示する</button>
            {{ end }}
            {{ else }}
            終了イベント一覧　　
            {{ if eq .Keywordev "" }}
            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode=1&keywordrm={{ .Keywordrm }}&userno={{ .Userno }}'">獲得ポイント詳細データのある終了イベントのみ表示する</button>
            {{ else }}
            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode=1&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&userno={{ .Userno }}'">獲得ポイント詳細データのある終了イベントのみ表示する</button>
            {{ end }}
            {{ end }}
        </p>
        {{/*
        <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }} （{{ UnixTimeToStr $tn }}）</div>
        <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }}
        </div>
        */}}
        <div style="text-indent: 2rem;">
            現時点では最初の50件のみ表示されます。
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
        {{ end }}
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