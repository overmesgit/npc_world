import sys
from PIL import Image
import os

def split_tiles(input_file, output_folder):
    # Open the input image
    with Image.open(input_file) as img:
        width, height = img.size
        
        # Check if the image dimensions are multiples of 32
        if width % 32 != 0 or height % 32 != 0:
            print("Error: Image dimensions must be multiples of 32.")
            return
        
        # Create the output folder if it doesn't exist
        os.makedirs(output_folder, exist_ok=True)
        
        # Split the image into 32x32 tiles
        for y in range(0, height, 32):
            for x in range(0, width, 32):
                box = (x, y, x + 32, y + 32)
                tile = img.crop(box)
                
                # Save the tile as a separate PNG file
                tile_filename = f"tile_{x//32}_{y//32}.png"
                tile_path = os.path.join(output_folder, tile_filename)
                tile.save(tile_path)
        
        print(f"Tiles have been saved to {output_folder}")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <input_file> <output_folder>")
    else:
        input_file = sys.argv[1]
        output_folder = sys.argv[2]
        split_tiles(input_file, output_folder)