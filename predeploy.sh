#!/bin/sh
set -e
make migrations-status-prod
make migrations-up-prod
make run