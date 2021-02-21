import sys
import os
import zlib

from typing import Any
from dataclasses import dataclass


from binascii import hexlify, unhexlify

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


@dataclass
class Blob():
    filename: str
    contents: str

    @classmethod
    def from_path(cls, path: str):
        contents = open(path).read()
        return Blob(filename=os.path.basename(path), contents=contents)
        
    def git_obj_contents(self):
        header = f"blob {len(self.contents)}\0"
        return (header + self.contents).encode()

    def sha(self):
        return hashlib.sha1(self.git_obj_contents()).hexdigest()

@dataclass
class TreeNode():
    name: str
    git_obj: Any

    def is_tree(self):
        return isinstance(self.git_obj, Tree)

    @classmethod
    def tree_from_path(cls, path):
        return TreeNode(
            name=os.path.basename(path),
            git_obj=Tree.from_path(path)
        )
        
    @classmethod
    def blob_from_path(self, path):
        return TreeNode(
            name=os.path.basename(path),
            git_obj=Blob.from_path(path)
        )

@dataclass
class Tree():
    nodes: Any

    @classmethod
    def from_path(cls, path: str, exclude: [str] = None) -> "Tree":
        if not exclude:
            exclude = []

        for root, dirs, files in os.walk(path):
            return Tree(
                nodes=[
                    TreeNode.tree_from_path(os.path.join(root, _dir)) for _dir in dirs
                    if _dir not in exclude
                ] + [
                    TreeNode.blob_from_path(os.path.join(root, _file)) for _file in files
                    if _file not in exclude
                ]
            )

    def git_obj_contents(self):
        contents = b"".join([
            f"{self.mode_from_node(node)} {node.name}\0".encode() + unhexlify(node.git_obj.sha().encode())
            for node in sorted(self.nodes, key=lambda x: x.name)
        ])
        header = f"tree {len(contents)}\0".encode()
        return (header + contents)

    def mode_from_node(self, node):
        if node.is_tree():
            return "40000"
        else:
            return "100644"

    def sha(self):
        return hashlib.sha1(self.git_obj_contents()).hexdigest()

            


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
        blob = Blob.from_path(filepath)
        sha = blob.sha()
        path = f".git/objects/{sha[0:2]}/{sha[2:]}"
        os.makedirs(os.path.dirname(path), exist_ok=True)

        zlib_store = zlib.compress(blob.git_obj_contents())
        open(path, "wb").write(zlib_store)
        print(sha)
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
    elif command == "write-tree":
        tree = Tree.from_path(".", exclude=[".git"])
        sha = tree.sha()
        path = f".git/objects/{sha[0:2]}/{sha[2:]}"
        os.makedirs(os.path.dirname(path), exist_ok=True)

        zlib_store = zlib.compress(tree.git_obj_contents())
        open(path, "wb").write(zlib_store)
        print(sha)
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
