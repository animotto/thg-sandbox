sessionretry = 3600
checkretry = 50

while (true) {
	print("Getting session ID: ")
	auth()
	if (sid = sessionid()) {
		print(sid + "\n")
	} else {
		print("Error\n")
	}

	for (i=1; i<=sessionretry; i+=checkretry) {
		print("Checking connection: ")
		if (checkconn()) {
			print("OK\n")
		} else {
			print("Error\n")
		}

		sleep(checkretry)
	}
}
