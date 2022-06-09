#!/usr/bin/env python3

#
# Json tags of struct fields in .go files are expected to be in format:
#   - camel case
#   - abbreviations and acronyms that are not at the beginning of the word are capitalized (eg. clientID, eventsURL; but caPool)
#   - generally, the tag should match the name of the field
#
# Additionally, some structs are used in network communication have defined schemas in swagger.yaml files. This goes for:
#   cloud2cloud-connector/swagger.yaml defines LinkedCloud schema
#   http-gateway/swagger.yaml defines several components
#
# This script verifies that:
# 1) tags names begin with a lower-case letter and contain only expected (a-zA-Z0-9_) characters.
# 2) components with defined schemas in selected swagger files must have existing fields in Go structs with matching
#    protobuf/json tags.
#

import argparse
import fnmatch
import os
import re
import sys
import yaml

SCRIPT_DIRECTORY = os.path.dirname(os.path.realpath(__file__))
ROOT_DIRECTORY = os.path.realpath(SCRIPT_DIRECTORY + "/../..")

DESCRIPTION = 'Validate json tags in Go files'

parser = argparse.ArgumentParser(description=DESCRIPTION)
parser.add_argument('-v', '--verbose', help="print verbose output",  action='store_true')
parser.add_argument('-f', '--fields', help="print warnings if field names do not match json tags",  action='store_true')
args = parser.parse_args()

PROTO_NAME_PATTERN = re.compile("protobuf:\".*name=([^\",]+)[^\"]*\"")
PROTO_JSON_PATTERN = re.compile("protobuf:\".*json=([^\",]+)[^\"]*\"")
JSON_PATTERN = re.compile("json:\"([^\"]+)\"")

def fill_protojson_tag_from_str(proto_tags, line, file):
  """Extract json name from protobuf field annotation."""
  match = re.search(PROTO_JSON_PATTERN, line)
  if not match:
    # json field inside protobuf tag exists only if the json name is different form name
    match = re.search(PROTO_NAME_PATTERN, line)
  if not match:
    return line
  line = line.replace(match.group(), "")
  proto_tags.setdefault(file, []).append(match.group(1))
  return line

def fill_json_tag_from_str(json_tags, line, file):
  """Extract json tag from field annotation."""
  match = re.search(JSON_PATTERN, line)
  if not match:
    return ""
  json_tag = match.group(1).split(",")[0]
  if not json_tag or json_tag == "-":
    return ""
  json_tags.setdefault(file, []).append(json_tag)
  return json_tag

def fill_proto_and_json_tags_from_file(proto_tags, json_tags, file):
  """Extract all protobuf and json tags from given file and save them as list in dictionary."""
  first_warning = True
  with open(file, "r") as f:
    for line in f:
      line = line.strip()
      if line.startswith("//"):
        continue
      if 'protobuf:"' in line:
        line = fill_protojson_tag_from_str(proto_tags, line, file)
        continue

      if not 'json:' in line:
        continue
      json_tag = fill_json_tag_from_str(json_tags, line, file)

      if not json_tag or not args.fields:
        continue
      field_name = line.split(" ", 1)[0]
      field_name_cannonical = field_name.replace('_', '').lower()
      json_tag = json_tag.replace('_', '')
      json_tag_cannonical = json_tag.replace('_', '').lower()
      if field_name_cannonical != json_tag_cannonical:
        if first_warning:
          first_warning = False
          print("file: {}".format(file))
        print("\tWARNING: field name '{}' does not match json tag '{}'".format(field_name, json_tag))

def get_all_proto_and_json_tags():
  """Find all protobuf and json tags in .go files from given directory."""
  proto_tags = {}
  json_tags = {}
  exclude_dirs = set(["bundle", "dependency", "charts"])
  for root, dirnames, filenames in os.walk(ROOT_DIRECTORY, topdown=True):
    dirnames[:] = [d for d in dirnames if d not in exclude_dirs]
    for filename in fnmatch.filter(filenames, "*.go"):
      file = root + "/" + filename
      fill_proto_and_json_tags_from_file(proto_tags, json_tags, file)

  return proto_tags, json_tags

def validate_tag_format(tag):
  """Check that tag begins with lower-case letter and contains only supported characters."""
  valid = True
  # must be alphanumeric
  if not re.match("^[a-z][a-zA-Z0-9_]*$", tag):
    valid = False

  if args.verbose:
    print("tag {}: {}".format(tag, valid))

  if not valid:
    print("ERROR: invalid tag {}".format(tag), file=sys.stderr)
  return valid

def find_and_validate_json_fields():
  """Find all protobuf and json tags in Go files from given directory and validate them."""
  proto_tags, json_tags = get_all_proto_and_json_tags()

  valid = True
  for tags in json_tags.values():
    for tag in tags:
      valid = validate_tag_format(tag) and valid

  return valid

if __name__ == "__main__":
  find_and_validate_json_fields() or sys.exit(1)
