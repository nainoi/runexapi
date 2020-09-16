#!/bin/bash
GIT_COMMIT=$(git log -1 --format=%h)
docker image build -t registry.thinkdev.app/think/runex/runexapi:$GIT_COMMIT -f Dockerfile.uat .