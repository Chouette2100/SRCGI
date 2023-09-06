<p style="padding-left:4em">獲得ポイントは定期の獲得ポイント取得で得たものか、<br>
    最後に「イベント参加者ルーム情報の追加と更新」を実行したときのものです。<br>
    「データ取得」にチェックの入っていないルームの獲得ポイントは定期的に更新されませんが、<br>
    「イベント参加者ルーム情報の追加と更新」で<span
        style="color:red">表示されたもの</span>はその時点で更新されます。<br>
    このあたりの仕組はちょっとわかりにくいのですが、やってるうちにわかってくると思います (^^;;<br>
    表はイベント開始前はフォロワー数によって、開始後は獲得ポイントによってソートされています。</p>

<font color="blue">
    <form action="new-user" method="GET">
        <p style="padding-left:4em"><span style="font-weight:bold;">一覧にないルームの追加</span>　　ユーザーID：
            <input type="hidden" name="eventid" value="{{.Eventid}}">
            <input type="hidden" name="func" value="newuser">
            <input type="text" name="roomid" value="999999" required pattern="[0-9]+">　
            <input type="submit" value="登録"><br>
            イベント参加ルームの追加は「イベント参加者ルーム情報の追加と更新」で行うのですが、<br>
            順位に関係なく特定のルームを追加したい、参加ルームが多すぎてルームのサムネがイベントページに表示されていない、<br>
            などのケースはここでユーザーIDを指定して追加することができます。<br>
            ユーザーIDというのはプロフィールやファンルームのページのURLの最後にある6桁（か5桁？）の数字のことです。<br>
            <br>
            <font color="red">
            今後、イベント参加者のリストやルーム名から追加するルームを選択する方法を作成する予定です。
            </font>
        </p>
    </form>
</font>
{{/*}}
<p style="padding-left:4em">
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
                現在のルーム名だけでなく過去のルーム名（のうち最近のもの、<br>
                例えば「夜風」さん）も検索対象となります。ただしルームの検索<br>
                結果は30件までしか表示されませんので1文字とかやめましょう。<br>
                下の検索結果からルームを選択してください。
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
                {{ $userno := .Userno }}
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
                <--! リストが表示されないときは二回クリックお願いします。<br>usernoがわかっていたら直接入力可！
                    -->
            </td>
            <td></td>
        </tr>

    </table>
    <input type="hidden" name="keywordrm" value="{{ .Keywordrm }}" />
    <input type="hidden" name="mode" value="{{ .Mode }}" />
    <input type="hidden" name="path" value="4" />
</form>

<!-- 5. ユーザ番号で選択する -->
ルームID(Room_id)を指定してルームを追加する
<form>
    <table>
        <tr>
            <td style="width:4em"></td>
            <td>ルームID</td>
            <td><input value="{{ .Userno }}" name="userno" type="number"></td>
            <td align="right">
                <input type="submit" value="データを取得するルームを追加する" formaction="newuser" formmethod="GET">
            </td>
        </tr>
        <tr>
            <td style="width:4em"></td>
            <td colspan="2">
                ルームIDはプロフィールやファンルームのURLの最後の"ID="の<br>
                あとにある整数です（６桁が多い）<BR>
                ルームIDの一部を指定しての検索はできません。
            </td>
            <td></td>
        </tr>
    </table>
    <input type="hidden" name="mode" value="{{ .Mode }}" />
    <input type="hidden" name="path" value="5" />
</form>
</p>
{{*/}}

</body>

</html>