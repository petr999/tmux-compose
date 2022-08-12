#!/usr/bin/env bash

cd /path/to/dumbclicker

tmux new -s dumbclicker-compose '
	docker-compose up
	bash -l
' \; neww -n dumbclicker_nginx_1 '
	PID=0
	try_next=1
	trap '\''
		echo "trap pid: ${PID}"
		kill -INT $PID
		try_next=""
	'\'' SIGINT
	while [ '\''x1'\'' == "x${try_next}" ]; do
		bash -lc '\''
			docker attach dumbclicker_nginx_1
			sleep 1
		'\'' &
		PID=$!
		echo "pid: ${PID}"
		wait $PID
	done
	trap - SIGINT
	bash -l
' \; neww -n dumbclicker_h2o_1 '
	PID=0
	try_next=1
	trap '\''
		echo "trap pid: ${PID}"
		kill -INT $PID
		try_next=""
	'\'' SIGINT
	while [ '\''x1'\'' == "x${try_next}" ]; do
		bash -lc '\''
			docker attach dumbclicker_h2o_1
			sleep 1
		'\'' &
		PID=$!
		echo "pid: ${PID}"
		wait $PID
	done
	trap - SIGINT
	bash -l
' \; neww -n dumbclicker_dumbclicker_1 '
	PID=0
	try_next=1
	trap '\''
		echo "trap pid: ${PID}"
		kill -INT $PID
		try_next=""
	'\'' SIGINT
	while [ '\''x1'\'' == "x${try_next}" ]; do
		bash -lc '\''
			docker attach dumbclicker_dumbclicker_1
			sleep 1
		'\'' &
		PID=$!
		echo "pid: ${PID}"
		wait $PID
	done
	trap - SIGINT
	bash -l
'