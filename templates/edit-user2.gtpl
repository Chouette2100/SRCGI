
{{ range . }}
{{/*
<form acttion="edit-user?eventid={{.Eventid}}" method="POST" id="{{.Formid}}"></form>
*/}}
<form acttion="edit-user" method="POST" id="{{.Formid}}"></form>
{{ end }}
<table border="1">
<tr>
<th>ルーム名</th>
<th>Prof./FR/Cnt.</th>
<th>ジャンル</th>
<th align="center" style="border-right-style:none;">ランク</th>
<th align="right" style="border-left-style:none;"></th>
<th>レベル</th>
<th>フォロ数</th>
<th>獲得ポイント</th>
<th>表示名</th>
<th>短縮<br>表示名</th>
<th>データ<br>取得</th>
<th>貢献<br>取得</th>
<th>表・グラフ<br>表示</th>
<th>グラフの色</th>
<th>実行</th>
</tr>

{{ $e := "" }}

{{ range . }}
<tr>

<td><a href="https://www.showroom-live.com/{{.Account}}">{{.Name}}
</a><input type="hidden" name="userid" value="{{.Userno}}" form="{{.Formid}}">
</a><input type="hidden" name="eventid" value="{{.Eventid}}" form="{{.Formid}}">
{{ $e = .Eventid }}
</td>
<td>
<a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">Prof.</a>/
<a href="https://www.showroom-live.com/room/fan_club?room_id={{.Userno}}">FR</a>/
<a href="https://www.showroom-live.com/event/contribution/{{.Eventid}}?room_id={{.Userno}}">Cnt.</a>
</td>
<td>{{.Genre}}</td>
<td align="center" style="border-right-style:none;">{{.Rank}}</td>
<td align="right" style="border-left-style:none;">{{.Nrank}}</td>
<td align="right">{{.Slevel}}</td><td align="right">{{.Sfollowers}}</td><td align="right">{{.Spoint}}</td>

<td><input type="text" name="longname"  size="8" value="{{.Longname}}"  form="{{.Formid}}" ></td>

<td><input type="text" name="shortname" size="4" value="{{.Shortname}}" form="{{.Formid}}" ></td>

<td align="center"><input type="checkbox" name="istarget" value="1" {{.Istarget}} form="{{.Formid}}" ></td>
<td align="center"><input type="checkbox" name="iscntrbpoint"    value="1" {{.Iscntrbpoint}}    form="{{.Formid}}" ></td>
<td align="center"><input type="checkbox" name="graph"    value="1" {{.Graph}}    form="{{.Formid}}" ></td>
<td>
     <svg width="40.00" height="19.00"
          xmlns="http://www.w3.org/2000/svg"
          xmlns:xlink="http://www.w3.org/1999/xlink">
          <rect x="1.00" y="1.00" width="39.00" height="18.00" stroke="white" stroke-width="0.1" />
          <line x1="5.00" y1="10.00" x2="35.00" y2="10.00" stroke="{{.Colorvalue}}" stroke-width="4.80" />
     </svg>
     <select name="color" form="{{.Formid}}">
          {{ range .Colorinflist }}
          <option value="{{.Color}}" {{.Selected}}>{{.Color}}</option>
          {{ end }}
     </select>
</td>

<td>
     <input type="hidden" name="func" value="edituser" form="{{.Formid}}" />
     <input type="submit" value="更新" form="{{.Formid}}" />
</td>
</tr>
{{ end }}

</table>
<br>
{{/*
<form acttion="edit-user" method="POST" id="getAllCntrb"></form>
<input type="hidden" name="eventid" value="{{ $e }}" form="getAllCntrb" />
<input type="hidden" name="func" value="getAllCntrb" form="getAllCntrb" />
<input type="submit" value="このイベントに参加しているすべてのルームについて枠別貢献ランキングを取得する" form="getAllCntrb" />
<br>
*/}}

<p style="padding-left:4em">イベント参加者として登録されているルームの一覧です<br><br>
<span style="color:red">ただし、イベント開始前は以下の「一覧にないルームの追加」で追加したルームのみ表示され、<br>
イベント開始後、指定した順位の範囲にあるルームが自動的に追加されます。</span><br><br>
「データ取得」にチェックが入っているルームが定期の獲得ポイント取得及び「直近の獲得ポイント」表示の対象です。<br>
「表・グラフ表示」にチェック入っているルームデータのみがリストあるいはグラフに表示されます（「直近の獲得ポイント」は除きます）<br>
リスト・グラフのルーム名には「表示名」が使われます<br>
「グラフの色」の色見本は更新ボタンを押したあと入力値を反映します。<br>
「表示名」、「データ取得」、「表・グラフ表示」、「グラフの色」は変更できますが一行ごとに編集し、更新ボタンを押してください。<br>
<span style="color:red">二行まとめて編集するのはできませんのでご注意ください。</span>
</p>

