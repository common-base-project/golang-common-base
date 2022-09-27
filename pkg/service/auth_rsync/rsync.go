package auth_rsync

import (
	"time"
)

func Main() {
	for range time.Tick(300 * time.Second) {
		//_ = GetAllDeparts()
		//GetAllUsers()
	}
}
