#!/bin/bash
# fail out of the script if anything here fails
set -e

# clear out stuff generated by test.sh
rm -rf src test_program_bin toml-test
