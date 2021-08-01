import sys
import os

database_file_path = sys.argv[1]
command = sys.argv[2]

def parse_varint(stream):
    """
    Parses a varint
    """
#     shift = 0
#     result = 0
#
#     while True:
#         i = int.from_bytes(stream.read(1), "big")
#         result |= (i & 0x7f) << shift
#         shift += 7
#         if not (i & 0x80):
#             break
#
#     return result
    acc = []
    while True:
        b = int.from_bytes(stream.read(1), "big")
        print('read byte: ', f"\\x{format(b, '02x')} | 0b{format(b, '08b')}")
        acc.append(b & 0b01111111)
        if not b & 0b10000000:
            break

    print('acc', acc)

    num = 0
    for b in reversed(acc):
        num = (num << 7) | b

    return num


def parse_column_value(stream, serial_type):
    if (serial_type >= 13) and (serial_type % 2 == 1):
        # Text encoding
        n_bytes = (serial_type - 13) // 2
        return stream.read(n_bytes)
    elif serial_type == 1:
        # 8 bit twos-complement integer
        return int.from_bytes(stream.read(1), "big")
    else:
        raise f"Unhandled serial_type {serial_type}"


def parse_record(stream, column_count):
    number_of_bytes_in_header = parse_varint(stream)
    print('number_of_bytes_in_header', number_of_bytes_in_header)

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
        print('Page size is ', page_size)

        database_file.seek(100) # Skip the header section

        page_type = int.from_bytes(database_file.read(1), "big")
        if page_type != 13:
            print("Expected first page to be a leaf table b-tree page.")
            exit(1)

        first_freeblock_start = int.from_bytes(database_file.read(2), "big")
        number_of_cells = int.from_bytes(database_file.read(2), "big")
        start_of_cell_content_area = int.from_bytes(database_file.read(2), "big")

        print('first_freeblock_start', first_freeblock_start)
        print('number_of_cells', number_of_cells)
        print('start_of_cell_content_area', start_of_cell_content_area)

        database_file.seek(100+8) # Skip the database header & b-tree page header, get to the cell pointer array

        cell_pointers = [int.from_bytes(database_file.read(2), "big") for _ in range(number_of_cells)]
        print(cell_pointers)

        # Each of these cells represents a row in the sqlite_schema table.
        for cell_pointer in cell_pointers:
            print(f"Visiting cell pointer {cell_pointer}")
            database_file.seek(cell_pointer)
            number_of_bytes_in_payload = parse_varint(database_file)
            rowid = parse_varint(database_file)

            parse_record(database_file, 5) # Table contains (type, name, tbl_name, rootpage, sql)
#             print('payload', database_file.read(number_of_bytes_in_payload))
#             number_of_bytes_in_record_header = parse_varint(database_file)
#             print('n bytes header', number_of_bytes_in_record_header)

    print("DBINFO!")
else:
    print(f"Invalid command: {command}")
