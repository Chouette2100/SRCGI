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
        href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a>
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
            <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td>終了イベント一覧</td>
        </tr>
    </table>
    <br>
    終了イベントの検索
    <form> <!-- 0. 最初のパス -->
        <br>
        <br>
        イベント名に含まれる文字列で絞り込む
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
                        「含まれる文字列」は例えば、スタートダッシュ、花火、Music 等です。
                        <br>カタカナとひらがな、全角と半角、英大文字と小文字あたりはアバウトです。
                    </td>
                    <td></td>
                </tr>
            </table>
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="1" />
        </form>
        <!-- 2. イベントID(Event_url_key)で絞り込む -->
        イベントID(Event_url_key)で絞り込む
        <form>
            <table>
                <tr>
                    <td style="width:4em"></td>
                    <td>イベントIDに含まれる文字列</td>
                    <td><input value="{{ .Kwevid }}" name="kwevid" pattern="[0-9A-Za-z-_=?]+"></td>
                    <td align="right">
                        <input type="submit" value="イベントIDで絞り込む" formaction="closedevents" formmethod="GET" {{/*
                            style="background-color: aquamarine" */}}>
                    </td>
                </tr>
                <tr>
                    <td style="width:4em"></td>
                    <td colspan="2">
                        ここでいうイベントIDというのはイベントページのURLの最後のフィールド<br>
                        のことです（本来のイベントIDは5桁程度の整数です）<br>
                        例えば花火のイベントであれば"fireworks"が含まれていることが多いです。
                    </td>
                    <td></td>
                </tr>
            </table>
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="2" />
        </form>

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
                        現在のルーム名だけでなく過去のルーム名（のうち最近のもの）も検索対象となります。<br>
                        ただしルームの検索結果は30件までしか表示されませんので1文字とかやめましょう。<br>
                        <span style="font-weight: bold;">下の検索結果からルームを選択し決定ボタンを押してください</span>
                    </td>
                    <td></td>
                </tr>
            </table>
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
                    <td>
                        <input style="margin-left: 2em" type="submit" value="決定" formaction="closedevents" formmethod="GET">
                    </td>
                    {{/* */}}
                    <td>
                        <input style="margin-left: 2em" type="submit" value="このルームの過去イベントを探す(数十秒かかることがあります)" formaction="oldevents" formmethod="GET" {{/*
                            style="background-color: aquamarine" */}}>
                    </td>
                    {{/* */}}
                </tr>
                <tr>
                    <td style="width:4em"></td>
                    <td colspan="2">
                        {{/*}}
                        リストが表示されないときは二回クリックお願いします。<br>usernoがわかっていたら直接入力可！
                        {{*/}}
                    </td>
                    <td></td>
                    <td></td>
                </tr>

            </table>
            <input type="hidden" name="keywordrm" value="{{ .Keywordrm }}" />
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="4" />
        </form>

        <!-- 5. ユーザ番号で選択する -->
        エントリーしたルームのルームID(Room_id)で絞り込む
        <form>
            <table>
                <tr>
                    <td style="width:4em"></td>
                    <td>ルームID</td>
                    <td><input value="{{ .Userno }}" name="userno" type="number"></td>
                    <td>
                        <input type="submit" value="ルームIDで絞り込む" formaction="closedevents" formmethod="GET" {{/*
                            style="background-color: aquamarine" */}}>
                    </td>
                    {{/* */}}
                    <td>
                        <input style="margin-left: 2em" type="submit" value="このルームIDの過去イベントを探す(数十秒かかることがあります)" formaction="oldevents" formmethod="GET" {{/*
                            style="background-color: aquamarine" */}}>
                    </td>
                    {{/* */}}
                </tr>
                <tr>
                    <td style="width:4em"></td>
                    <td colspan="3">
                        ルームIDはプロフィールやファンルームのURLの最後の"room_id="の<br>
                        あとにある整数です（６桁が多い）<BR>
                        ルームIDの一部を指定しての検索はできません。
                    </td>
                    <td></td>
                    <td></td>
                </tr>
            </table>
            <input type="hidden" name="mode" value="{{ .Mode }}" />
            <input type="hidden" name="path" value="5" />
        </form>
        <br>
        <br>
        {{ if or ( eq .Keywordrm "" ) ( ne .Userno 0) }}
        <p id="result">
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
            {{ if ge .Path 4 }}
            （ルーム名からの検索では開催中イベントを含む）
            {{ end }}
            {{ if eq .Keywordev "" }}
            　　<button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode=1&keywordrm={{ .Keywordrm }}&userno={{ .Userno }}'">獲得ポイント詳細データのあるイベントのみ表示する</button>
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
            一覧は51件程度ずつ表示され、50件ずつスクロールされます。データ上最終結果が存在しても表示されないケースがあります（修正・改良予定あり）
        </div>
        
        {{ if ne .Offset 0 }}

            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode={{ .Mode }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&kwevid={{ .Kwevid }}&userno={{ .Userno }}&action=top&limit={{ .Limit }}&offset={{.Offset}}'">最初から表示する</button>

            <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode={{ .Mode }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&kwevid={{ .Kwevid }}&userno={{ .Userno }}&action=prev&limit={{ .Limit }}&offset={{.Offset}}'">前ページ</button>


        {{ end }}




        <table border="1" style="border-collapse: collapse">
            <tr bgcolor="gainsboro" style="text-align: center">
                <td style="border-right: none;">イベント名とイベントページへのリンク</td>
                <td style="border-left: none;">イベントID</td>
                <td>開始日時</td>
                <td>終了日時</td>
                <td>最終結果</td>
                <td>表示項目<br>選択画面</td>
                <td>最終獲得<br>ポイント表</td>
                <td>獲得ポイント<br>推移図</td>
                <td>日々の<br>獲得pt</td>
                <td>枠毎の<br>獲得pt</td>
                <td>貢献<br>pt</td>

            </tr>
            {{ $i := 0 }}
            {{ range .Eventinflist }}


        {{ if eq $i 1 }}
            {{/*
            {{ if eq (Divide .Aclr 2 ) 0 }}
            <tr bgcolor="gainsboro">
            {{ else }}
            <tr bgcolor="palegreen">
            {{ end }}
            */}}
            <tr bgcolor="gainsboro">
            {{ $i = 0 }}
        {{ else }}
            {{/*
            {{ if eq (Divide .Aclr 2 ) 0 }}
            <tr>
            {{ else }}
            <tr bgcolor="lightblue">
            {{ end }}
            */}}
            <tr>
            {{ $i = 1 }}
        {{ end }}
        
                <td style="border-right: none;">
                    {{ if IsTempID .Event_ID }}
                        {{ .Event_name }}
                    {{ else }}
                        <a href="https://showroom-live.com/event/{{ .Event_ID }}">{{ .Event_name }}</a>
                    {{ end }}
                </td>
                <td style="border-left: none; text-align: right;">
                    {{ if ne .I_Event_ID 0 }}
                        {{ .I_Event_ID }}
                    {{ end }}
                </td>
                <td>
                    {{ TimeToStringY .Start_time }}
                </td>
                <td>
                    {{ TimeToString .End_time }}
                </td>
                <td style="text-align: center;">
                    {{/* {{ if and (ne .I_Event_ID 0) ( ne .Aclr 0 ) }} */}}
                    {{/* {{ if and (ne .I_Event_ID 0) ( ne .Aclr 0 ) }} */}}
                    {{ if  ne .Aclr 0 }}
                    {{/* {{ if ne .I_Event_ID 0 }} */}}
                    <a href="closedeventroomlist?eventid={{ .I_Event_ID }}&eventurlkey={{ .Event_ID }}{{ if ne $userno 0 }}&roomid={{$userno}}{{end}}">最終結果</a>
                    {{ else }}
                        {{ if ne $userno 0 }}
                        <a href="https://showroom-live.com/api/events/{{.I_Event_ID}}/ranking?room_id={{$userno}}">API</a>
                        {{ end }}
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="top?eventid={{ .Event_ID }}">表示項目選択</a>
                    {{ end }}
                </td>
                <td style="text-align: center;">
                    {{ if eq .Target 1 }}
                    <a href="list-lastP?eventid={{ .Event_ID }}{{ if ne $userno 0 }}&roomid={{$userno}}{{end}}">リスト</a>
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
                <td style="text-align: center;">
                    {{ if ne $userno 0 }}

                    {{ if IsTempID .Event_ID }}
                    ーー /
                    {{ else }}
                    <a href="https://www.showroom-live.com/event/contribution/{{ DelBlockID .Event_ID }}?room_id={{$userno}}">公式</a> /
                    {{ end }}
                    {{ if ne .I_Event_ID 0 }}
                    <a href="contributors?ieventid={{.I_Event_ID}}&roomid={{$userno}}">CSV</a>
                    {{ else }}
                    ーー
                    {{ end }}
                    {{ end }}
                </td>
            </tr>
            {{ end }}
            <tr>
        </table>

        {{/* // TODO: SQL取得件数が .Limit であっても、.Totalcount が必ずしも .Limit にならないことへの暫定対応
        {{ if eq .Totalcount .Limit }}
        */}}
        {{ if ne .Totalcount 0 }}
           <button type="button"
                onclick="location.href='closedevents?path={{ .Path }}&mode={{ .Mode }}&keywordev={{ .Keywordev }}&keywordrm={{ .Keywordrm }}&kwevid={{ .Kwevid }}&userno={{ .Userno }}&action=next&limit={{ .Limit }}&offset={{.Offset}}'">次ページ</button>
        {{ end }}

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
