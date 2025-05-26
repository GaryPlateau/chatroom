var chatlistMan = {
	message: {
		add: function (message, type) {
			var chat_body = $('.layout .content .chat .chat-body');
			if (chat_body.length > 0) {

				type = type ? type : '';
				//message = message ? message : 'I did not understand what you said!';

				$('.layout .content .chat .chat-body .messages').append(`<div class="message-item ` + type + `">
					<div class="message-avatar">
						<figure class="avatar">
							<img src="` + (type == 'outgoing-message' ? myPhoto : toPhoto) + `" class="rounded-circle">
						</figure>
						<div>
							<h5>` + (type == 'outgoing-message' ? myNickname : toNickname) + `</h5>
							<div class="time">14:50 PM ` + (type == 'outgoing-message' ? '<i class="ti-check"></i>' : '') + `</div>
						</div>
					</div>
					<div class="message-content">
						` + message + `
					</div>
				</div>`);

				setTimeout(function () {
					chat_body.scrollTop(chat_body.get(0).scrollHeight, -1).niceScroll({
						cursorcolor: 'rgba(66, 66, 66, 0.20)',
						cursorwidth: "4px",
						cursorborder: '0px'
					}).resize();
				}, 200);
			}
		},
		del: function (e) {
			
		}
	}
};

$(function () {

    /**
     * Some examples of how to use features.
     *
     **/

    

    //setTimeout(function () {
        // $('#disconnected').modal('show');
        // $('#call').modal('show');
        // $('#videoCall').modal('show');
        //$('#pageTour').modal('show');
    //}, 1000);

    /*$(document).on('submit', '.layout .content .chat .chat-footer form', function (e) {
        e.preventDefault();

        var input = $(this).find('input[type=text]');
        var message = input.val();

        message = $.trim(message);

        if (message) {
            SohoExamle.Message.add(message, 'outgoing-message');
            input.val('');

            setTimeout(function () {
                SohoExamle.Message.add();
            }, 1000);
        } else {
            input.focus();
        }
    });*/

    $(document).on('click', '.layout .content .sidebar-group .sidebar .list-group-item', function () {
        if (jQuery.browser.mobile) {
            $(this).closest('.sidebar-group').removeClass('mobile-open');
        }
    });

});

/*
ws.onmessage = function(e)
var msg = JSON.parse(e.data)
var sender, username, name_list
switch(msg.type)
	case "system":
	sender = "系统消息";
	break;
	case "user":
	sender = msg.from + ":";
	break;
	case "handshake":
	var user_info = {"type":"login", "content":uname}
	sendMsg(user_info)
	return;
	case "login":
	case "logout":
	username = msg.content;
	name_list = msg.user_list;
	change_type = msg.type;
	dealUser(user_name, change_type, name_list);
	return
	*/