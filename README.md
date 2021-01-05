# monhttp - Monitor your servers

Status Page for monitoring your websites and applications with graphs and analytics.

## Run everywhere

monhttp is written in Go(lang). All you need is the precompiled binary based on your operating system, and the
HTML/CSS/Javascript files. You can even run monhttp on your Raspberry Pi.

## Notifications

monhttp can notify you via email or Telegram when a service is unavailable. More notification types coming soon.

## Run on Docker

Use the [official Docker image](https://hub.docker.com/r/koloooo/monhttp) to run monhttp in seconds.

``` shell
docker run -p 8081:8081 koloooo/monhttp
```

To save the config.yml from the container for later, you need to mount the path `/monhttp/config`. Add a volume for this
when starting the container.

``` shell
docker run -p 8081:8081 -v your_path:/monhttp/config koloooo/monhttp
```

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

After the initial setup, there is a config.yml in the config folder. This file can be used to change or save the
configuration.

| Key  | Value  | Description  |
|---|---|---|
|  database.host | localhost  |   |
|  database.name |  monhttp |   |
|  database.password |  top_secret |   |
|  database.port | 5432  |   |
|  database.user | monhttp_user  |   |
|   |   |   |
|  notifier |   |   |
|   |   |   |
|  server.port | 8081  |   |
|   |   |   |
|  users |   | A list in the format "name:password" you can add here as many users as you want to  |
