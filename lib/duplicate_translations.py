import os

def remove_colon_and_strip_from_str(line):
    """
    Removes the colon from the line and strips the line.
    """
    return line.replace(":", "").strip()

def print_result_for_file(file_path, file_name, duplicate_entries):
    """
    Prints the result for a given file.
    """
    print(f"-" * 50)
    print(f"Processing file: {file_name} \n")
    # keep entries with more than one occurence
    if len(duplicate_entries) > 0:
        for key, value in duplicate_entries.items():
            # prints debug info for each duplicate found
            print(f"{file_path}: \"{key}\" appears in lines {", ".join(map(str, value))}")
    else:
        # prints debug info, if no duplicates found 🥳
        print(f"No duplicates found!")

def find_duplicates(directory, debug=False):
    """
    Reads all .txt files in a given directory and tries to detect duplicate entries.
    """
    try:
        language_lookup_table = {}
        files = [f for f in os.listdir(directory) if f.endswith('.txt')]
        files.sort()        
        
        if not files:
            print("No .txt files found in the directory.")
            return

        for file_name in files:
            # if file_name contains "-help" skip it for now!
            if "-help" in file_name:
                continue

            file_path = os.path.join(directory, file_name)
            try:
                with open(file_path, 'r', encoding='utf-8') as file:
                    lookup_table = {}
                    for line_number, line in enumerate(file, start=1):
                        stripped_line = line.strip()
                        stripped_keywords = stripped_line.split(":")
                        stripped_keywords = list(map(remove_colon_and_strip_from_str, stripped_keywords))
                        trimmed_keywords = list(map(str.strip, stripped_keywords))
                        
                        for tk in trimmed_keywords:
                            if tk == "" or tk.isdigit():
                                continue
                            if tk in lookup_table:
                                lookup_table[tk].append(line_number)
                            else:
                                lookup_table[tk] = [line_number]
                    duplicate_entries = {k: v for k, v in lookup_table.items() if len(v) > 1}
                    print_result_for_file(file_path, file_name, duplicate_entries)
                    language_lookup_table[file_name] = duplicate_entries
            except Exception as e:
                print(f"An error occurred while processing the file: {e}")
            
        return language_lookup_table
    except Exception as e:
        print(f"An error occurred: {e}")
        return None

# Example usage from the root directory:
#   python ./lib/duplicate_translations.py
if __name__ == "__main__":
    directory_path = "share/translations/"
    find_duplicates(directory_path, debug=True)
