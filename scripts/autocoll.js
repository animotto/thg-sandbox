retry=300

var sid
if (!(sid=sessionid())) {
	print("Getting session ID: ")
	auth()
	if (sid=sessionid()) {
		print(sid + "\n")
	} else {
		print("Error\n")
	}
}

if (sid) {
	while (true) {
		if (!(w=world())) {
			print("Can't get world\n")
		} else {
			if (w.Bonuses.length!=0) {
				for (i=0; i<w.Bonuses.length; i++) {
					print("Collecting bonus " + w.Bonuses[i].Id +
						" with " + w.Bonuses[i].Amount + " credits: ")
					if (b=boncoll(w.Bonuses[i].Id)) {
						print("OK\n")
					} else {
						print("Error\n")
					}
				}
			}
			if ((w.Goals.length!=0) && (gt=goaltypes())) {
				for (i=0; i<w.Goals.length; i++) {
					print("Finishing goal " + w.Goals[i].Id + ": ")
					if (g=goalupd(w.Goals[i].Id, gt.Gtypes[w.Goals[i].Gtype].Max)) {
						print("OK\n")
					} else {
						print("Error\n")
					}
				}
			}
		}

		if (args[0]) {
			sleep(retry)
		} else {
			break
		}
	}
}
