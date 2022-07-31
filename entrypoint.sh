#!/bin/bash

set -xe

HORUSEC_VERSION=""
if [ "$INPUT_HORUSEC_VERSION" != "latest" ]; then
  HORUSEC_VERSION="/tags/${INPUT_HORUSEC_VERSION}"
fi

curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/master/deployments/scripts/install.sh | bash -s ${HORUSEC_VERSION}

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

horusec -p ${INPUT_WORKING_DIRECTORY} -o json -O ${OUT_OPTION} --log-level TRACE ${HORUSEC_ARGS_OPTION}
commenter