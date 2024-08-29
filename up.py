import os

def extract_and_update_files(input_file):
    current_file = None
    current_content = []

    def write_current_file():
        nonlocal current_file, current_content
        if current_file and current_content:
            with open(current_file, 'w') as f:
                f.write(''.join(current_content))
            print(f"Updated/Created: {current_file}")
        current_file = None
        current_content = ['package main\n']

    with open(input_file, 'r') as f:
        for line in f:
            sline = line.strip()
            if sline.startswith('// ') and sline.endswith('.go'):
                write_current_file()
                current_file = sline[3:]  # Remove the '// ' prefix
            elif current_file is not None:
                current_content.append(line)

    write_current_file()  # Write the last file

if __name__ == "__main__":
    input_file = "input.txt"
    if not os.path.exists(input_file):
        print(f"Error: {input_file} not found.")
    else:
        extract_and_update_files(input_file)
        print("File extraction and update complete.")