package main

type ProcessObj struct {
	Name       string
	Invocation string
	Id         int
}

type ProcessGroupObj struct {
	Name      string
	Processes []*ProcessObj
}

var sshdProc = &ProcessObj{
	Name:       "sshd",
	Invocation: "/usr/sbin/sshd",
	Id:         3838,
}

var dbProc = &ProcessObj{
	Name:       "database",
	Invocation: "/etc/db/db",
	Id:         489437,
}

var dbWatchProc = &ProcessObj{
	Name:       "db-watcher",
	Invocation: "/etc/db/db-watch",
	Id:         23733,
}

var dbPoolerProc = &ProcessObj{
	Name:       "db-pool",
	Invocation: "/etc/db/db-pool",
	Id:         3453,
}

var sshdGroup = &ProcessGroupObj{
	Name: "sshd",
	Processes: []*ProcessObj{
		sshdProc,
	},
}

var dbGroup = &ProcessGroupObj{
	Name: "db",
	Processes: []*ProcessObj{
		dbProc,
		dbWatchProc,
		dbPoolerProc,
	},
}

var ProcessList []*ProcessObj = []*ProcessObj{
	sshdProc,
	dbProc,
	dbWatchProc,
	dbPoolerProc,
}

var ProcessGroups []*ProcessGroupObj = []*ProcessGroupObj{
	sshdGroup,
	dbGroup,
}
