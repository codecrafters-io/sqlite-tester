import sys
import os
import zlib

import hashlib
import pathlib


def main():
    command = sys.argv[1]
    if command == "init":
        os.mkdir(".git")
        os.mkdir(".git/objects")
        os.mkdir(".git/refs")
        with open(".git/HEAD", "w") as f:
            f.write("ref: refs/heads/master\n")

        print("Initialized git directory")
    elif command == "cat-file":
        sha = sys.argv[3]
        obj_path = f".git/objects/{sha[0:2]}/{sha[2:]}"
        compressed = open(obj_path, "rb").read()
        uncompressed = zlib.decompress(compressed)
        sys.stdout.buffer.write(uncompressed.split(b"\0")[-1])
    else:
        raise RuntimeError(f"Unknown command: #{command}")


main()
