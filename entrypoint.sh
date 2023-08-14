#!/bin/bash

set -xe

#HORUSEC_VERSION="latest"
#if [ "$INPUT_HORUSEC_VERSION" != "latest" ]; then
#  HORUSEC_VERSION="tags/${INPUT_HORUSEC_VERSION}"
#fi
#
#wget -O - -q "$(wget -q https://api.github.com/repos/ZupIT/horusec/releases/${HORUSEC_VERSION} -O - | grep -m 1 -o -E "https://.+?horusec_linux_amd64" | head -n1)" > horusec_linux_amd64
#
#install horusec_linux_amd64 /usr/local/bin/horusec

curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/main/deployments/scripts/install.sh | bash -s latest-beta

COMMENTER_VERSION="latest"

wget -O - -q "$(wget -q https://api.github.com/repos/upbanx/horusec-action/releases/${COMMENTER_VERSION} -O - | grep -o -E "https://.+?horusec-commenter-linux-amd64")" > horusec-commenter-linux-amd64
wget -O - -q "$(wget -q https://api.github.com/repos/upbanx/horusec-action/releases/${COMMENTER_VERSION} -O - | grep -o -E "https://.+?checksums.txt")" > commenter.checksums

grep horusec-commenter-linux-amd64 commenter.checksums > commenter-linux-amd64.checksum
sha256sum -c commenter-linux-amd64.checksum
install horusec-commenter-linux-amd64 /usr/local/bin/commenter

if [ -n "${GITHUB_WORKSPACE}" ]; then
  cd "${GITHUB_WORKSPACE}" || exit
fi

if [ -n "${INPUT_ARGUMENTS}" ]; then
  HORUSEC_ARGS_OPTION="${INPUT_ARGUMENTS}"
fi

OUT_OPTION="results.json"

horusec start -p ${INPUT_WORKING_DIRECTORY} -o json -O ${OUT_OPTION} --log-level TRACE "${HORUSEC_ARGS_OPTION}"
commenter