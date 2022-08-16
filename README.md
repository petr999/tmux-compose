# tmux-compose
Tmux-based UI for Docker-Compose

[![Tmux-Compose demo](https://i9.ytimg.com/vi_webp/x4ZODQq-2EA/mqdefault.webp?sqp=CLjU2pcG&rs=AOn4CLDAJ7X1WM3Cyrnmk0p4KjBIMuL7xw)](http://www.youtube.com/watch?v=x4ZODQq-2EA "Tmux-Compose in action")

## Intro

Hi I'm Peter. I'm a fan of command line interface and a virtualization. After years of using Docker and a Tmux console manager  I made my launcher of a services for Docker-compose based on Tmux, Shell and Go language.

It's enough to have a docker-compose configuration to run your services in a named file system directory, which is where your application typically use to reside, but you can customize how you should run your services in a desired ways according to the launcher templates and environment variables. It's not necessary to have a template as the built-in one can launch your docker-compose services each in its own window, neither you need to know much about Tmux sessions. But if docker-compose is your preferred way to deploy the app it's currently the only way to have your tmux running detached for unattended init of your services on every your server's start or reboot.

Let me introduce you to Tmux-compose which is basically a convinient way to run docker-compose under tmux in a managed way. It's not necessary to have the Tmux-compose on your server as with Tmux-compose dry run feature you can be provided the Shell command line to copy and paste to your machine. Another interesting feature is: no matter if your particular container is created, running and attached, you can have its dedicated console right away to poke around it.

This may effectively bring your shell back as a tool to investigate post-mortems and observe log messages interactively. You may want to raise up the scroll buffer of tmux for that, and a tmux-compose templates may suit for that purpose, too. With tools like FCGI::Spawn and Debug-Fork-Tmux modules you may even debug the particular backend request on your deployment right away each on its own console even if it's a separate process.

The overall project is open source in Go language and the statically linked binary is on its way, too.

## Synopsis

Run tmux-compose from your directory with docker-compose.yml.
```
  cd your-app
  ls -l docker-compose.yml
  go run /path/to/tmux-compose
```
## Dependencies

Runtime dependencies are:

- tmux v3+

- shell: bash, zsh

- docker-compose

- docker-compose.yml with at least one element in its `services:` section. See `run/run_test/testdata/dumbclicker/docker-compose.yml`

Compile-time dependencies are:

- go 1.18+

- yaml.v2

## Configuration

Apart of having `docker-compose.yml`, every another configuration variable is optional. While `gson` template defines the way `tmux` and `docker-compose` to run, the `tmux-compose` behavior is controlled by environment variables.

Refer to `.env-sample` for information on environment variables and apply environment variables the appropriate way. You can use it as a tenplate if you apply with `.env`:
```
  cp -v .env-sample .env
```

### Environment

- `TMUX_COMPOSE_DC_YML` points to directory with your `docker-compose.yml`, or to configuration file that you'd want to use with `tmux-compose`. Be sure to use a particular template with `-f` argument supplied for `docker-compose` in that case

- `TMUX_COMPOSE_DRY_RUN` any non empty value triggers dry run mode with shell script contents to standard output instead of running commands for you

- `TMUX_COMPOSE_TEMPLATE_FNAME` points to directory with your 'tmux-compose-template.gson' template, or to a template file to use itself.