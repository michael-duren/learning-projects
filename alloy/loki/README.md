# Alloy + Loki playground

A Go log generator, Alloy, Loki, and Grafana. The app uses only Go's standard library:

```text
app → /logs/app.log → Alloy loki.source.file → loki.process → loki.write → Loki
```

The app emits JSON every second and repeats info, debug, warn, and error events.
It also prints the same events to stdout for `docker compose logs app`, but Alloy
reads the shared file, not Docker's stdout logs.

## Make shortcuts

Run these from `alloy/loki/`. Use `make help` to list all commands.

```sh
make up          # Build and start the stack
make logs-app    # Follow the app's raw output; Ctrl+C stops following
make logs-alloy  # Follow processed entries and Alloy diagnostics
make apply       # After editing config.alloy: validate, then restart Alloy
make app         # After editing app/main.go: rebuild and restart the app
make down        # Stop the stack and keep data
```

`make reset` deletes all playground volumes and starts fresh. New Alloy stages
only affect newly collected logs. View those logs in Grafana Explore.

## Run

From this directory, with Docker and Docker Compose installed:

```sh
docker compose up --build -d
docker compose logs -f app alloy loki
```

Open Grafana at <http://localhost:3001/explore>. Anonymous Editor access is enabled for
this local playground. Select **Loki**, switch the query editor to **Code**, enter
`{service_name="demo-app"}`, and click **Run query**. Set the time range to the
last five minutes. Expand a log line to inspect its fields and labels. Use
`{service_name="demo-app",level="error"}` to show only errors.

Grafana queries Loki; it does not store another copy of the logs.
See the [Grafana log visualization guide](https://grafana.com/docs/loki/latest/visualize/grafana/).

Open Alloy at <http://localhost:12345> to inspect the component graph and health.
Loki's API is at <http://localhost:3100>; it has no log browsing UI. Allow about
30 seconds for Loki to become ready and Alloy to send its first batch:

```sh
curl -fsS http://localhost:3100/ready
curl -fsSG http://localhost:3100/loki/api/v1/query_range \
  --data-urlencode 'query={service_name="demo-app"}' \
  --data-urlencode 'since=5m' \
  --data-urlencode 'limit=10'
```

To select errors, change the query to `{service_name="demo-app",level="error"}`.
To inspect stored label sets:

```sh
curl -fsSG http://localhost:3100/loki/api/v1/series \
  --data-urlencode 'match[]={service_name="demo-app"}'
```

## Where data is saved

Docker named volumes hold the data, not files in this repository:

| Volume | Container path | Contents |
| --- | --- | --- |
| `alloy_logs` | `/logs/app.log` in app and Alloy | Raw JSON lines appended by the Go app. Alloy mounts this read-only. |
| `alloy_loki-data` | `/loki` in Loki | Loki's chunks, index, and write-ahead log. This is Loki's own storage format, not a plain-text log file. |
| `alloy_alloy-data` | `/var/lib/alloy/data` in Alloy | File read positions, so restarts resume tailing. |
| `alloy_grafana-data` | `/var/lib/grafana` in Grafana | Grafana's database and settings. |

Inspect the actual host storage location:

```sh
docker volume inspect alloy_logs alloy_loki-data alloy_alloy-data alloy_grafana-data
```

On Linux these usually live under Docker's data directory. Docker Desktop keeps
them inside its Linux VM. Use the returned Mountpoint rather than assuming a path.
You can read the raw file without locating it on the host:

```sh
docker compose exec alloy tail -n 10 /logs/app.log
```

`docker compose down` preserves these volumes. `docker compose down -v` deletes
them. The Compose project name stays `alloy` to reuse data from the original
`alloy/` setup after this move. Give future sibling playgrounds different project
names and host ports when running them together.

## What to edit

- `compose.yaml` connects services, mounts files, and publishes localhost ports.
- `app/main.go` defines the sample events. Rebuild the app after editing it.
- `config.alloy` collects, processes, and sends the logs. This is the experiment file.
- `loki.yaml` configures a single Loki instance with local filesystem storage.

The baseline Alloy config has only JSON parsing and a `level` label stage.
The source supplies the static `service_name` label and adds a `filename` label.
`__path__` is an internal file target setting, not a stored Loki label.

There are three distinct things to work with:

| Data | Baseline behavior |
| --- | --- |
| Log line | The original JSON string is stored unchanged. |
| Extracted values | `stage.json` makes fields available to later stages; these values alone are not stored as labels. |
| Stream labels | `stage.labels` promotes `level` to an indexed label. `level = ""` means use the extracted field with the same name. |

The name `"app"` in `loki.process "app"` identifies a component for wiring;
it does not add an `app` label to your logs. Loki groups entries by the complete
set of stream labels. Keep frequently changing values such as `request_id` out
of stream labels.

## Try stages

Insert one example at a time at the experiment comment inside `loki.process`.
Stages run in order. These examples are optional; none is enabled initially.

Add a constant label:

```alloy
stage.static_labels {
  values = { environment = "local" }
}
```

Promote the parsed route to a label:

```alloy
stage.labels {
  values = { route = "path" }
}
```

Drop debug events using the label created earlier:

```alloy
stage.match {
  selector = "{level=\"debug\"}"
  action   = "drop"
}
```

Alternatively, drop health checks using an extracted value:

```alloy
stage.drop {
  source = "path"
  value  = "/health"
}
```

Use the event's timestamp instead of the time Alloy read the line:

```alloy
stage.timestamp {
  source = "timestamp"
  format = "RFC3339Nano"
}
```

Store only the message as the log line. The extracted values remain available
to subsequent stages, but the original JSON is no longer the stored line:

```alloy
stage.output {
  source = "message"
}
```

Keep a request ID as structured metadata without making a stream for each ID:

```alloy
stage.structured_metadata {
  values = { request_id = "" }
}
```

Remove the automatically added filename label:

```alloy
stage.label_drop {
  values = ["filename"]
}
```

See Grafana's [processing stage reference](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.process/)
for more stages, including regex, logfmt, replace, template, and multiline.
The [file source reference](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.source.file/)
explains targets and file positions.

## Apply changes

Validate, then restart Alloy to load the edited file:

```sh
docker compose run --rm --no-deps alloy validate /etc/alloy/config.alloy
docker compose restart alloy
```

Changes affect newly collected lines. Existing Loki entries keep their old
labels and content, so use a short query window such as `since=30s` when comparing.
Alloy's read positions persist across restarts.

After editing the generator:

```sh
docker compose up -d --build app
```

Stop the playground while preserving data:

```sh
docker compose down
```

Reset everything, including stored logs and read positions:

```sh
docker compose down -v
docker compose up --build -d
```

This is a local learning setup. The log file grows until you reset the volumes;
stop the services when finished experimenting.

## See Alloy output directly

Run `make logs-alloy`. The `loki.echo.processed` component prints each processed
entry to Alloy stdout, including its labels and log text. It receives the same
entries as `loki.write.local`, after all processing stages. Look for
`component_id=loki.echo.processed`. Delivery diagnostics appear in the same stream;
echoing an entry does not confirm that Loki accepted it.

After changing stages, run `make apply` and watch new entries.
