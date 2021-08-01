import sys

from record_parser import parse_record
from varint_parser import parse_varint

database_file_path = sys.argv[1]
command = sys.argv[2]

if command == ".dbinfo":
    with open(database_file_path, "rb") as database_file:
        # Refer to https://www.sqlite.org/fileformat.html for file format specification
        header_string = database_file.read(16)

        if header_string != b"SQLite format 3\x00":
            print("Invalid database file header.")
            exit(1)

        database_file.seek(100) # Skip the header section

        page_type = int.from_bytes(database_file.read(1), "big")
        if page_type != 13:
            print("Expected first page to be a leaf table b-tree page.")
            exit(1)

        _first_freeblock_start = int.from_bytes(database_file.read(2), "big")
        number_of_cells = int.from_bytes(database_file.read(2), "big")

        database_file.seek(100+8)  # Skip the database header & b-tree page header, get to the cell pointer array

        cell_pointers = [int.from_bytes(database_file.read(2), "big") for _ in range(number_of_cells)]

        sqlite_schema_rows = []

        # Each of these cells represents a row in the sqlite_schema table.
        for cell_pointer in cell_pointers:
            database_file.seek(cell_pointer)
            number_of_bytes_in_payload = parse_varint(database_file)
            rowid = parse_varint(database_file)
            record = parse_record(database_file, 5)

            # Table contains columns: type, name, tbl_name, rootpage, sql
            sqlite_schema_rows.append({
                'type': record[0],
                'name': record[1],
                'tbl_name': record[2],
                'rootpage': record[3],
                'sql': record[4],
            })

        print(f"number of tables: {len(sqlite_schema_rows)}")
else:
    print(f"Invalid command: {command}")
