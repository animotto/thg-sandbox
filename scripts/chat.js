var room=999
var chatretry=5

if (auth()) {
	var lm
	while (true) {
		if (chat=getchat(room, lm)) {
			for (i=0; i<chat.Messages.length; i++) {
				print("["+chat.Messages[i].DateTime+"] "+chat.Messages[i].Name+": "+chat.Messages[i].Message+"\n")
				lm=chat.Messages[i].DateTime
			}
		}
		sleep(chatretry)
	}
}
