<p style="padding-left:4em">獲得ポイントは定期の獲得ポイント取得で得たものか、最後に「イベント参加者ルーム情報の追加と更新」を実行したときのものです。<br>
「データ取得」にチェックの入っていないルームの獲得ポイントは定期的に更新されませんが、「イベント参加者ルーム情報の追加と更新」で<span style="color:red">表示されたもの</span>はその時点で更新されます。<br>
このあたりの仕組はちょっとわかりにくいのですが、やってるうちにわかってくると思います (^^;;<br>
表はイベント開始前はフォロワー数によって、開始後は獲得ポイントによってソートされています。</p>

<font color="blue">
<form action="new-user" method="GET">
<p style="padding-left:4em"><span style="font-weight:bold;">一覧にないルームの追加</span>　　ユーザーID：
<input type="hidden" name="eventid" value="{{.Eventid}}" >
<input type="hidden" name="func" value="newuser" >
<input type="text" name="roomid" value="999999" required pattern="[0-9]+" >　
<input type="submit" value="登録"><br>
イベント参加ルームの追加は「イベント参加者ルーム情報の追加と更新」で行うのですが、順位に関係なく特定のルームを追加したい、参加ルームが多すぎてルームのサムネがイベントページに表示されていない、などのケースはここでユーザーIDを指定して追加することができます。ユーザーIDというのはプロフィールやファンルームのページのURLの最後にある6桁（か5桁？）の数字のことです。<br>イベント開始後は追加しようとするルームがイベントに参加しているかチェックしていますが、イベント開始前はチェックしていません（イベント開始前はイベント参加の有無をイベントページからしか行えず、その場合すべてのイベント参加ルームのリストを入手できるとは限らないため、というのが表向きの理由ですが、作り込みが面倒でパフォーマンスが悪化するということもあります。参加ルームがそれほど多くないイベントではチェックできるので、そのうちチェックをいれるかも）
</p>
</form></font>
</body>
</html>
