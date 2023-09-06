package main

func isPkInWhitelist(targetPk string) bool {
    for i := 0; i < len(whitelist); i++ {
        if whitelist[i].Pk == targetPk {
            return true
        }
    }
    return false
}

func deleteFromWhitelistRecursively (target string) {
	var updatedWhitelist []User
	var queue []string

	for _, user := range whitelist {
		if user.Pk != target {
			updatedWhitelist = append(updatedWhitelist, user)
		}
		if user.InvitedBy == target {
			queue = append(queue, user.Pk);
		}
	}

	whitelist = updatedWhitelist
	for _, pk := range queue {
		deleteFromWhitelistRecursively(pk)
	}
}
