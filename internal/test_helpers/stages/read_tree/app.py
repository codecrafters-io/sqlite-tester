import sys
import os
import zlib

from dataclasses import dataclass

from binascii import hexlify

import hashlib
import pathlib
from io import BytesIO


@dataclass
class TreeEntry:
    mode: str
    filename: str
    sha: str

    def __repr__(self):
        return f"TreeEntry('{self.filename}', '{self.mode}', '{self.sha[0:4]}')"

    @property
    def is_dir(self):
        return self.mode == "40000"


def main():
    command = sys.argv[1]
    if command == "init":
        os.mkdir(".git")
        os.mkdir(".git/objects")
        os.mkdir(".git/refs")
        with open(".git/HEAD", "w") as f:
            f.write("ref: refs/heads/master\n")

        print("Initialized git directory")
    elif command == "hash-object":
        assert sys.argv[2] == "-w"
        filepath = sys.argv[3]
        contents = open(filepath).read()
        header = f"blob {len(contents)}\0"
        store = (header + contents).encode()
        sha = hashlib.sha1(store).hexdigest()
        print(sha)
        zlib_store = zlib.compress(store)
        path = f".git/objects/{sha[0:2]}/{sha[2:]}"
        os.makedirs(os.path.dirname(path), exist_ok=True)
        open(path, "wb").write(zlib_store)
    elif command == "cat-file":
        sha = sys.argv[3]
        obj_path = f".git/objects/{sha[0:2]}/{sha[2:]}"
        compressed = open(obj_path, "rb").read()
        uncompressed = zlib.decompress(compressed)
        sys.stdout.buffer.write(uncompressed.split(b"\0")[-1])
    elif command == "ls-tree":
        assert sys.argv[2] == "--name-only"
        sha = sys.argv[3]
        obj_path = f".git/objects/{sha[0:2]}/{sha[2:]}"
        compressed = open(obj_path, "rb").read()
        uncompressed = zlib.decompress(compressed)
        contents = BytesIO(uncompressed)
        header = read_until_null_byte(contents)
        objs = []
        try:
            while True:
                mode = read_until_space(contents)
                filename = read_until_null_byte(contents)
                sha = contents.read(20)
                objs.append(
                    TreeEntry(
                        mode=mode.decode(),
                        filename=filename.decode(),
                        sha=hexlify(sha).decode(),
                    )
                )
        except EOFError:
            pass

        for obj in objs:
            print(obj.filename)
    else:
        raise RuntimeError(f"Unknown command: #{command}")


def read_until_sep(io: BytesIO, sep: bytes) -> bytes:
    contents = b""
    while True:
        byte = io.read(1)
        if len(byte) == 0:
            raise EOFError
        if byte == sep:
            break
        contents += byte

    return contents


def read_until_null_byte(io: BytesIO) -> bytes:
    return read_until_sep(io, b"\0")


def read_until_space(io: BytesIO) -> bytes:
    return read_until_sep(io, b" ")


main()
