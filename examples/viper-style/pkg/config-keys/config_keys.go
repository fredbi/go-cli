package cli_test

const (
	/*
		YAML config accepted:

		app:
		  url: xyz
		  parallel: 99
		  user: abc
		  child:
		    workers: 5
		log:
		  level: info
	*/
	KeyLog      = "app.log.level"
	KeyURL      = "app.url"
	KeyParallel = "app.parallel"
	KeyUser     = "app.user"
	KeyDry      = "run.dryRun"
	KeyWorkers  = "app.child.workers"
)
