#!/usr/bin/env python3

#
# Keys in yaml files are expected to be in format:
#   - camel case
#   - abbreviations and acronyms that are not at the beginning of the word are capitalized (eg. clientID, eventsURL; but caPool)
#
# Additionally, keys in config.yaml files should have corresponding yaml tags in golang structs.
#
# This script verifies that
# 1) key names begin with a lower-case letter and contain only expected (a-zA-Z0-9_\-) characters.
# 2) a .go file with a field with the given yaml tag exists
#

import argparse
import fnmatch
import os
import re
import sys
import yaml

SCRIPT_DIRECTORY = os.path.dirname(os.path.realpath(__file__))
ROOT_DIRECTORY = os.path.realpath(SCRIPT_DIRECTORY + "/../..")

DESCRIPTION = 'Validate keys format in YAML files'

parser = argparse.ArgumentParser(description=DESCRIPTION)
parser.add_argument('-v', '--verbose', help="print verbose output",  action='store_true')
parser.add_argument('-f', '--fields', help="print warnings if field names do not match yaml tags",  action='store_true')
args = parser.parse_args()

YAML_PATTERN = re.compile("yaml:\"([^\"]+)\"")

def validate_yaml_key_format(key):
  """Check that key begins with lower-case letter and contains only supported characters."""
  valid = True
  # must be alphanumeric
  if not re.match("^[a-z0-9][a-zA-Z0-9_]*$", key):
    valid = False

  if args.verbose:
    print("key {}: {}".format(key, valid))

  if not valid:
    print("ERROR: invalid key {}".format(key), file=sys.stderr)
  return valid

def validate_yaml_keys_format(data):
  """Recursively validate all keys in a dictionary."""
  valid = True
  for k, v in data.items():
    valid = validate_yaml_key_format(str(k)) and valid
    if isinstance(v, dict):
      valid = validate_yaml_keys_format(v) and valid
    elif isinstance(v, list):
      for item in v:
        if isinstance(item, dict):
          valid = validate_yaml_keys_format(item) and valid
  return valid

def find_and_validate_yaml_file(file):
  """Validate format of yaml keys from given file."""
  with open(file, "r") as f:
    try:
      if args.verbose:
        print("{}".format(file))
      data = yaml.safe_load(f)
      valid = validate_yaml_keys_format(data)
    except yaml.YAMLError as exc:
      print(exc)
  return valid

def find_and_validate_yaml_files(dir = ROOT_DIRECTORY):
  """Find all yaml files in directory and validate them."""

  valid = True
  exclude_dirs = set(["dependency", "templates", ".github"])
  exclude_filenames = set(["swagger.yaml"])
  for root, dirnames, filenames in os.walk(dir, topdown=True):
    dirnames[:] = [d for d in dirnames if d not in exclude_dirs]
    for filename in fnmatch.filter(filenames, "*.yaml"):
      if filename in exclude_filenames:
        continue
      file = root + "/" + filename
      valid = find_and_validate_yaml_file(file) and valid

  return valid

if __name__ == "__main__":
  find_and_validate_yaml_files() or sys.exit(1)
