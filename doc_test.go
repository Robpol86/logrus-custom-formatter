package lcf

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func ExampleCustomHandlers() {
	LoadAverage := func(e *logrus.Entry, f *CustomFormatter) (interface{}, error) {
		someNumber := 0.3
		return someNumber, nil
	}
	template := "[%04[relativeCreated]d] %1.2[loadAvg]f %7[levelName]s %[message]s\n"
	formatter := NewFormatter(template, CustomHandlers{"loadAvg": LoadAverage})

	log := logrus.New()
	log.Formatter = formatter
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout
	log.Debug("A group of walrus emerges from the ocean")
	log.Warn("The group's number increased tremendously!")
	log.Info("A giant walrus appears!")
	log.Error("Tremendously sized cow enters the ocean.")

	// Output:
	// [0000] 0.30   DEBUG A group of walrus emerges from the ocean
	// [0000] 0.30 WARNING The group's number increased tremendously!
	// [0000] 0.30    INFO A giant walrus appears!
	// [0000] 0.30   ERROR Tremendously sized cow enters the ocean.
}
