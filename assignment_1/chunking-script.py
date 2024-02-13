import os

def split_file(input_file_path, output_file_prefix, num_chunks=8):
    output_directory = "data"
    if not os.path.exists(output_directory):
        os.makedirs(output_directory)

    with open(input_file_path, 'r') as input_file:
        headers = [next(input_file) for _ in range(2)]
        lines = input_file.readlines()
    
    total_data_rows = len(lines)
    rows_per_chunk = total_data_rows // num_chunks
    extra_rows = total_data_rows % num_chunks 

    start_index = 0 
    for i in range(num_chunks):
        chunk_size = rows_per_chunk + (1 if i < extra_rows else 0)
        end_index = start_index + chunk_size

        chunk_lines = lines[start_index:end_index]
        start_index = end_index

        output_filename = f"{output_directory}/{output_file_prefix}{i}.txt"
        with open(output_filename, 'w') as output_file:
            output_file.writelines(headers + chunk_lines)

    print(f"Split into {num_chunks} chunks")

input_file_path = "master_file.txt"
output_file_prefix = "chunk_"
num_chunks = int(os.environ["CHUNKS"])

split_file(input_file_path, output_file_prefix, num_chunks)
