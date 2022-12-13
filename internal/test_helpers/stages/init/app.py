import sys

database_file_path = sys.argv[1]
command = sys.argv[2]

if command == ".dbinfo":
    with open(database_file_path, "rb") as database_file:
        # Refer to https://www.sqlite.org/fileformat.html for file format specification
        header_string = database_file.read(16)

        if header_string != b"SQLite format 3\x00":
            print("Invalid database file header.")
            exit(1)

        page_size = int.from_bytes(database_file.read(2), "big")

        print(f"database page size: {page_size}")
else:
    print(f"Invalid command: {command}")
