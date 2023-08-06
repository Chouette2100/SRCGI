<tr><td></td><td>目標ポイント</td><td><input type="text" name="target" value="{{.Target}}" size="7" required pattern="[0-9]+"></td></tr>
<tr><td></td><td>最大表示数</td><td><input type="text" name="maxdsp" value="{{.Maxdsp}}" size="3" required pattern="[1-9][0-9]*"></td></tr>
<tr><td></td><td>カラーマップ</td><td><input type="text" name="cmap" value="{{.Cmap}}" size="2" required pattern="[12]"
    title="原則2を指定します。1は旧バージョンの色パターンです。"></td></tr>
<tr><td></td><td></td><td align="right"><input type="submit" value="設定変更" formaction="param-eventc" formmethod="POST" style="background-color: khaki"></td></tr>
<tr><td></td><td></td><td align="right"><input type="submit" value="キャンセル" formaction="top?eventid={{.Event_ID}}" formmethod="POST" style="background-color: khaki"></td></tr>
</table>
</form>
<p style="padding-left:4em">
「ＤＢに登録する順位の範囲」というのは、「登録」を実行した時点で指定した範囲にある順位のルームを登録する、ことを意味します。一度登録されたルームは削除されませんので、順位の入れ替わりがあったあと「登録」を行うとDBに登録されているルーム数は増えていきます。DBに登録されるルーム数には制限はありません。
</p>

</body>
</html>
