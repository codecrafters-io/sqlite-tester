#!/bin/bash

echo "segfault" > /dev/stderr
kill -SEGV $$