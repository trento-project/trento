"""
The ID checker is a support script to check if all the trento checks have an unique ID and
to set new IDs if they don't exist

:author: xarbulu
:organization: SUSE Linux GmbH
:contact: xarbulu@suse.com

:since: 2021-09-16
"""

import os
import logging
import random
import argparse

try:
    import yaml
except ModuleNotFoundError:
    logging.getLogger(__name__).error(
        "yaml package not found. To install it run `pip install PyYAML`")


CHECKS_FOLDER = "../roles/checks"
HEXDIGITS = "0123456789ABCDEF"
ID_LENGTH = 6
CHECK_ID = "id"
REQUIRED_FIELDS = ["id", "name", "group", "labels", "description", "remediation", "implementation"]


def parse_args():
    """
    Parse arguments
    """
    parser = argparse.ArgumentParser(description="Trento Ansible checks ID checker")
    parser.add_argument(
        "--generate", dest="generate", action="store_true",
        help="Generate a new ID on the checks that do not have it")

    return parser.parse_args()


def get_checks_folder():
    """
    Get the checks folder absolute path
    """
    return os.path.join(os.path.dirname(os.path.abspath(__file__)), CHECKS_FOLDER)

def create_id():
    """
    Create ID with specified lenght based on HEX digits
    """
    generated_id = "".join([random.choice(HEXDIGITS) for _ in range(ID_LENGTH)])
    return generated_id


def id_sanity_check(check_id):
    """
    Check ID syntax sanity check
    """
    str_check_id = str(check_id)
    if len(str_check_id) != ID_LENGTH:
        return False

    for char in str_check_id:
        if char not in HEXDIGITS:
            return False

    return True

def sanity_check(check_data, check_file, logger):
    """
    Check sanity
    """
    for field in REQUIRED_FIELDS:
        if field not in check_data:
            logger.error("field %s not found in check %s", field, check_file)
            return False

    if not id_sanity_check(check_data[CHECK_ID]):
        logger.error(
            "%s id (%s) does not follow the ID correct syntax "\
            "(%s chars length hex string)",
            check_file, check_data[CHECK_ID], ID_LENGTH)
        return False

    return True

def append_id_to_check(check_file, new_id):
    """
    Add the new ID to the check file
    """
    with open(check_file, "a") as write_ptr:
        write_ptr.write("\n")
        write_ptr.write(
            "# check {}. This value must not be changed over the life of this check\n".format(
                CHECK_ID))
        write_ptr.write("{}: {}\n".format(CHECK_ID, new_id))


def main(generate, logger):
    """
    Main method
    """
    id_list = []
    check_add_id = []
    checks_folder = get_checks_folder()
    for c_file in os.listdir(checks_folder):
        if os.path.isdir(os.path.join(checks_folder, c_file)):
            logger.info("check directory found: %s", c_file)
            check_path = os.path.join(checks_folder, c_file, "defaults/main.yml")
            try:
                with open(check_path) as file_ptr:
                    data = yaml.load(file_ptr, Loader=yaml.Loader)
                    if CHECK_ID not in data:
                        logger.info("check %s doesn't have the %s value", c_file, CHECK_ID)
                        if generate:
                            check_add_id.append(check_path)
                            continue
                        else:
                            logger.error("to add a new ids, use the --generate flag on the script")
                            return 1
                    if not sanity_check(data, check_path, logger):
                        return 1
                    if data[CHECK_ID] in id_list:
                        logger.error("%s %s already exists!", CHECK_ID, data[CHECK_ID])
                        return 1
                    id_list.append(data[CHECK_ID])
            except FileNotFoundError:
                logger.error("check %s doesn't have the defaults/main.yml file", c_file)
                continue

    if generate:
        logger.info("generating new ids...")
        for check_file in check_add_id:
            new_id = create_id()
            while new_id in id_list:
                new_id = create_id()

            id_list.append(new_id)
            append_id_to_check(check_file, new_id)
            logger.info("new id %s added to check %s", new_id, check_file)

    return 0


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    logger = logging.getLogger(__name__)
    args = parse_args()
    exit(main(args.generate, logger))
