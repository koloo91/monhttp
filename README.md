# monhttp - Monitor your servers

Status Page for monitoring your websites and applications with graphs and analytics.

## Run everywhere

`monhttp` is written in Go(lang). All you need is the precompiled binary based on your operating system, and the
HTML/CSS/Javascript files. You can even run `monhttp` on your Raspberry Pi.

## Notifications

`monhttp` can notify you via email or Telegram when a service is unavailable. More notification types coming soon.

It is possible to use your own template for notifications. The [golang template engine](https://golang.org/pkg/text/template/#example_Template) is used for this purpose. Possible variables are `{{.Name}}`, `{{.Reason}}` and `{{.Date}}`.

## Run on Docker

Use the [official Docker image](https://hub.docker.com/r/koloooo/monhttp) to run monhttp in seconds.

``` shell
docker run -p 8081:8081 koloooo/monhttp
```

To save the config.env from the container for later, you need to mount the path `/monhttp/config`. Add a volume for this
when starting the container.

``` shell
docker run -p 8081:8081 -v your_path:/monhttp/config koloooo/monhttp
```

## Use docker-compose

Simply run `docker-compose up` to start `monhttp` together with a postgres database. Open
up [`http://localhost:8081`](http://localhost:8081) in your browser and enjoy `monhttp`. The default user is `admin` and
the password is `admin` too.

## Build it locally

Make sure you have Go 1.15 and Node.js 14.15 installed on your computer. Clone the repository and execute the build
command.

``` shell
git clone git@github.com:koloo91/monhttp.git
cd monhttp
make buildLocal
```

Then you will find all files and folders in the dist folder. Change to the dist folder and start monhttp
with `./monhttp`.

## Configuration

After the initial setup, there is a config.env in the config folder. This file can be used to change or save the
configuration.

| Key  | Value  | Description  |
|---|---|---|
|  DATABASE_HOST | localhost  |   |
|  DATABASE_NAME |  monhttp |   |
|  DATABASE_PASSWORD |  top_secret |   |
|  DATABASE_PORT | 5432  |   |
|  DATABASE_USER | monhttp_user  |   |
|   |   |   |
|  NOTIFIER |   |   |
|   |   |   |
|  SERVER_PORT | 8081  |   |
|   |   |   |
|  USERS |   | A list in the format "name:password" you can add here as many users as you want to  |

You can also use environment variables to configure `monhttp`. Environment variables override the values from the `config.env` file.
