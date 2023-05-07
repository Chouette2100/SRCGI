<br>
<br>
ルーム名をクリックするとそのルームでの草うま王のファンレベルが表示されます。
<br>
<br>
<table border="1">
    <tr align="center">
        <td>ルーム名</td>
        <td>room_id</td>
    </tr>
    {{ range . }}
    <tr>
        <td><a href="/fanlevel?roomid={{ .Room_id}}">{{ .Room_name}}</a></td>
        <td align="right">{{ .Room_id}}</td>
    </tr>
    {{end}}
 </table>
<br>
<br>
</body>
</html>
