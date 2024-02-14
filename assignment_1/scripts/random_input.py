import random

def process_file(input_file, output_file):
  first_cols = []
  with open(input_file, "r") as f:
    for line in f:
      values = line.strip().split()
      if values:
        first_cols.append(values[0])

  random.shuffle(first_cols)

  with open(output_file, "w") as f:
    for col in first_cols:
      f.write(col + "\n")

  print(f"Output written to: {output_file}")

input_file = "input.txt"
output_file = "output.txt"
process_file(input_file, output_file)
