import sys

from varint import parse_varint

database_file_path = sys.argv[1]
command = sys.argv[2]


def parse_column_value(stream, serial_type):
    if (serial_type >= 13) and (serial_type % 2 == 1):
        # Text encoding
        n_bytes = (serial_type - 13) // 2
        return stream.read(n_bytes)
    elif serial_type == 1:
        # 8 bit twos-complement integer
        return int.from_bytes(stream.read(1), "big")
    else:
        raise Exception(f"Unhandled serial_type {serial_type}")


def parse_record(stream, column_count):
    _number_of_bytes_in_header = parse_varint(stream)

    serial_types = [parse_varint(stream) for i in range(column_count)]
    print('serial types', serial_types)

    column_values = [parse_column_value(stream, serial_type) for serial_type in serial_types]
    print('column values', column_values)

if command == ".dbinfo":
    with open(database_file_path, "rb") as database_file:
        # Refer to https://www.sqlite.org/fileformat.html for file format specification
        header_string = database_file.read(16)

        if header_string != b"SQLite format 3\x00":
            print("Invalid database file header.")
            exit(1)

        page_size = int.from_bytes(database_file.read(2), "big")

        database_file.seek(100) # Skip the header section

        page_type = int.from_bytes(database_file.read(1), "big")
        if page_type != 13:
            print("Expected first page to be a leaf table b-tree page.")
            exit(1)

        first_freeblock_start = int.from_bytes(database_file.read(2), "big")
        number_of_cells = int.from_bytes(database_file.read(2), "big")
        start_of_cell_content_area = int.from_bytes(database_file.read(2), "big")

        database_file.seek(100+8)  # Skip the database header & b-tree page header, get to the cell pointer array

        cell_pointers = [int.from_bytes(database_file.read(2), "big") for _ in range(number_of_cells)]
        print(cell_pointers)

        sqlite_schema_records = []

        # Each of these cells represents a row in the sqlite_schema table.
        for cell_pointer in cell_pointers:
            database_file.seek(cell_pointer)
            number_of_bytes_in_payload = parse_varint(database_file)
            rowid = parse_varint(database_file)

            # Table contains columns: type, name, tbl_name, rootpage, sql
            sqlite_schema_records.append(parse_record(database_file, 5))

        print('number')
else:
    print(f"Invalid command: {command}")
