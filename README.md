Viper configuration provider for Logrus Mate
============================================

This tiny project is a way to configure Logrus Mate using Viper
as the configuration source. It is meant to be used with the currently
unreleased version of Logrus Mate (`v1.0.0` doesn't have the configuration
providers, only `master` branch does).

As Viper is very different to other configuration options, there is
no way to mix it with e.g. HOCON. Basically, the assumption is that
your project uses Viper and you don't want any other configuration
system just for logging. If so, this library may come handy.

Here's a code snippet with usage example:

    loggingConfig := viper.Sub("logging")
	if loggingConfig != nil {
		mate, err := vipermate.NewMate(loggingConfig)
		if err != nil {
			logrus.WithError(err).Panic("Failed to process configuration")
		}
		if err = mate.Hijack(logrus.StandardLogger(), "main"); err != nil {
			logrus.WithError(err).Panic("Failed to configure main logger")
		}
	}

Then, an example config (I use YAML, YMMV) could look like, e.g.:

    logging:
      main:
        out:
          name: stdout
          options: {}
        level: debug
        formatter:
          name: text
          options:
            force-colors: true
            timestamp-format: "2006-01-02T15:04:05Z07:00"
            disable-timestamp: false
            full-timestamp: true
        hooks:
          # Let's say we want to set up an ELK hook
          logstash:
            protocol: tcp
            address: "elk.example.org:5959"
            name: "myapp"

Check out the source code for more details and information
about the limitations. In particular, non-string slice functions
are not implemented because Viper doesn't have anything like that.
