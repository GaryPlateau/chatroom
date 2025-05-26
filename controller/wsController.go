package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"michatroom/common"
	"michatroom/conf"
	"michatroom/driver"
	"michatroom/model"
	"michatroom/utils"
	"net/http"
)

func WsService(c *gin.Context, hub *model.Hub) {
	var user *model.UserBasic
	uuid := getClaims(c).Uuid
	err := json.Unmarshal([]byte(driver.RedisSingleInstance.HGetValue("users", uuid).(string)), &user)
	utils.ErrorHandler("wsServer获取用户数据失败", err)

	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clientConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	utils.ErrorHandler("ws创建失败", err)

	client := &model.Connection{Conn: clientConn, Message: make(chan []byte), Data: &model.Data{
		UUId:    uuid,
		GroupId: "0",
		ToUuid:  uuid,
		Type:    "",
	}}
	hub.Register <- client

	ufs := model.GetUserFriendsList(uuid)
	for _, fUser := range ufs {
		fid := fUser.FriendId
		_, ok := hub.Clients[fid]
		if ok {
			model.SingleMsgSend(hub, uuid, fid, []byte(user.Nickname+"上线了!"), 1)
		}
	}

	go client.WritePump()
	go client.ReadPump(hub)
}

func Home(ctx *gin.Context) {

	ctx.SetCookie("uuidtest", "user123456", 300, "/", conf.HttpAddr, false, false)
	//ctx.Next()
	cid := ctx.Param("cid")
	gid := ctx.Param("gid")
	tid := ctx.Param("tid")
	uuid, err := ctx.Cookie("uuidtest")
	fmt.Println("testuuid", uuid, err)
	homeTemplate.Execute(ctx.Writer, "ws://"+ctx.Request.Host+"/ws/server/"+cid+"/"+gid+"/"+tid)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
	var hostid = document.getElementById("hostid");
    var sendid = document.getElementById("sendid");
</script>
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))

func getClaims(c *gin.Context) *common.MyClaims {
	//c.SetSameSite(http.SameSiteNoneMode)
	token, err := c.Cookie("token")
	utils.ErrorHandler("获取usertoken失败", err)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "404.html", gin.H{
			"code": "401",
			"msg":  "获取usertoken失败",
		})
	}
	claims, err := common.ParseToken(token)
	utils.ErrorHandler("解析token失败", err)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "404.html", gin.H{
			"code": "401",
			"msg":  "解析token失败",
		})
	}
	return claims
}
