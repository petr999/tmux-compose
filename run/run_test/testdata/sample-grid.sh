#!/usr/bin/env bash

cd /path/to/dumbclicker-grid

tmux new -s "dumbclicker-grid-compose" '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker-grid-compose'\''
  docker-compose up
  echo "
Commands for your consideration:
  docker-compose up
  docker-compose up -d"
  echo "Welcome to shell."
  bash -l
' \; split-window '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker-grid_nginx_1'\''
  echo "Commands for your consideration:
  docker-compose up '\''nginx'\''
  docker-compose up -d '\''nginx'\''
  docker attach '\''dumbclicker-grid_nginx_1'\''
  docker exec -it '\''dumbclicker-grid_nginx_1'\''"
  echo "Hit Ctrl-C to break out to shell."
  echo "Attaching to container '\''dumbclicker-grid_nginx_1'\''..."
  PID=0
  try_next=1
  trap '\''
    kill -INT $PID
    try_next=""
  '\'' SIGINT
  while [ '\''x1'\'' == "x${try_next}" ]; do
    bash -lc '\''
      docker attach '\''dumbclicker-grid_nginx_1'\''
      sleep 1
    '\'' &
    PID=$!
    wait $PID
  done
  trap - SIGINT
  echo "
Commands for your consideration:
  docker-compose up '\''nginx'\''
  docker-compose up -d '\''nginx'\''
  docker attach '\''dumbclicker-grid_nginx_1'\''
  docker exec -it '\''dumbclicker-grid_nginx_1'\''"
  echo "Welcome to shell."
  bash -l
' \; split-window '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker-grid_h2o_1'\''
  echo "Commands for your consideration:
  docker-compose up '\''h2o'\''
  docker-compose up -d '\''h2o'\''
  docker attach '\''dumbclicker-grid_h2o_1'\''
  docker exec -it '\''dumbclicker-grid_h2o_1'\''"
  echo "Hit Ctrl-C to break out to shell."
  echo "Attaching to container '\''dumbclicker-grid_h2o_1'\''..."
  PID=0
  try_next=1
  trap '\''
    kill -INT $PID
    try_next=""
  '\'' SIGINT
  while [ '\''x1'\'' == "x${try_next}" ]; do
    bash -lc '\''
      docker attach '\''dumbclicker-grid_h2o_1'\''
      sleep 1
    '\'' &
    PID=$!
    wait $PID
  done
  trap - SIGINT
  echo "
Commands for your consideration:
  docker-compose up '\''h2o'\''
  docker-compose up -d '\''h2o'\''
  docker attach '\''dumbclicker-grid_h2o_1'\''
  docker exec -it '\''dumbclicker-grid_h2o_1'\''"
  echo "Welcome to shell."
  bash -l
' \; split-window '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker-grid_dumbclicker_1'\''
  echo "Commands for your consideration:
  docker-compose up '\''dumbclicker'\''
  docker-compose up -d '\''dumbclicker'\''
  docker attach '\''dumbclicker-grid_dumbclicker_1'\''
  docker exec -it '\''dumbclicker-grid_dumbclicker_1'\''"
  echo "Hit Ctrl-C to break out to shell."
  echo "Attaching to container '\''dumbclicker-grid_dumbclicker_1'\''..."
  PID=0
  try_next=1
  trap '\''
    kill -INT $PID
    try_next=""
  '\'' SIGINT
  while [ '\''x1'\'' == "x${try_next}" ]; do
    bash -lc '\''
      docker attach '\''dumbclicker-grid_dumbclicker_1'\''
      sleep 1
    '\'' &
    PID=$!
    wait $PID
  done
  trap - SIGINT
  echo "
Commands for your consideration:
  docker-compose up '\''dumbclicker'\''
  docker-compose up -d '\''dumbclicker'\''
  docker attach '\''dumbclicker-grid_dumbclicker_1'\''
  docker exec -it '\''dumbclicker-grid_dumbclicker_1'\''"
  echo "Welcome to shell."
  bash -l
' \; set -g pane-border-status bottom \; set-option status off \; select-layout tiled