#!/bin/sh

kubectl patch pod go-memtest --subresource resize --patch \
  '{"spec":{"containers":[{"name":"go-memtest", "resources":{"requests":{"memory":"2G"},"limits":{"memory":"2G"}}}]}}'