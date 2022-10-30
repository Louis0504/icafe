package box

import "fmt"

type defaults struct {
	// name
	appName string
}

func getDefaults() defaults {
	d := defaults{
		appName: "name",
	}

	return d
}

func defaultWarnLogMetric(appName string) string {
	return fmt.Sprintf("%s.log.warn.count", appName)
}

func defaultErrorLogMetric(appName string) string {
	return fmt.Sprintf("%s.log.error.count", appName)
}
