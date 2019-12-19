# config

tl;dr

For a YAML config file as follows:

```YAML
database:
  drv: pgx
  dsn: ${DB_DSN}
```

Use this code to retrieve data:

```Go
sql.Open(cfg.GetString("database.drv"), cfg.GetString("database.dsn"))
```

And this to initialize a `cfg`:

```Go
cfg, err := config.Load("./config/development.yml")
if err != nil {
	log.Fatalf("Loading config file failed: %v", err)
}
```


https://github.com/buildpeak/config
