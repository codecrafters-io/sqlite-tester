import sys
import os


# print(sys.argv)

command = sys.argv[1]
if command == "init":
    os.mkdir(".git")
    os.mkdir(".git/objects")
    os.mkdir(".git/refs")
    with open(".git/HEAD", "w") as f:
        f.write("ref: refs/heads/master\n")

    print("Initialized git directory")
