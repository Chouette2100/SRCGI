<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
{{/* <h2>{{.function}}</h2> */}}
<p style="padding-left:2em">
{{.Msg001}}{{.Eventid}}{{.Msg002}}
</p>
<p>
{{/* <button type="button" onclick="history.back()">このルームの表示項目選択</button><br> */}}
<button type="button" onclick="location.href='top?eventid={{.Eventid}}'">このルームの表示項目選択</button>
</p>
</body>
</html>
