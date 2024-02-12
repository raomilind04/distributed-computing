import os
from io import open



def split_file_row_wise(input_file_path, output_file_prefix, num_chunks=8):

    row_count = 0
    output_chunks = [[] for _ in range(num_chunks)] 
    os.mkdir("data")
    with open(input_file_path, 'r') as input_file:
        for line in input_file:
            row_count += 1
            chunk_index = min(range(num_chunks), key=lambda i: len(output_chunks[i]))
            output_chunks[chunk_index].append(line)

    for i, chunk_data in enumerate(output_chunks):
        output_filename = f"data/{output_file_prefix}{i}.txt"
        with open(output_filename, 'w') as output_file:
            output_file.writelines(chunk_data)

    print(f"split into {num_chunks} chunks.")


input_file_path = "master_file.txt"
output_file_prefix = "chunk_"
num_chunks = 4

split_file_row_wise(input_file_path, output_file_prefix, num_chunks)


