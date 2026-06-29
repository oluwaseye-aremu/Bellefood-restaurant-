import os
import glob

directory = "/Users/user/Bellefood-restaurant/frontend"
html_files = glob.glob(os.path.join(directory, "*.html"))

replacements = {
    "Lumière | Sisi Alase Restaurant": "24/7 BelleFood Restaurant",
    "Sisi Alase": "24/7 BelleFood",
    "Lumière": "24/7 BelleFood",
    "Mary Ademola": "King Mitchy",
    "Chef Peter Ayomide": "CEO, King Mitchy",
    "mary-ademola.jpeg": "bellefood_ceo.png",
    "Mon - Thu:</strong> 11:00 AM - 10:00 PM": "Open 24 Hours, 7 Days a Week!</strong>",
    "Fri - Sat:</strong> 11:00 AM - 11:00 PM": "We are always open to serve you the best meals in Lagos, anytime, anywhere.",
    "Sun:</strong> 10:00 AM - 9:30 PM": "Dine-in and Dispatch Available 24/7",
    "Happy Hour daily from 4:00 PM - 6:00 PM": "24/7 Dispatch and Dine-In",
    "1998": "recently",
    "Home Delivery": "Dispatch Delivery",
    "Enjoy our gourmet meals at home.": "Fast, reliable delivery to your doorstep."
}

for file_path in html_files:
    if os.path.basename(file_path) == "index.html":
        continue # Already edited
    
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()
    
    modified = False
    for old_str, new_str in replacements.items():
        if old_str in content:
            content = content.replace(old_str, new_str)
            modified = True
            
    if modified:
        with open(file_path, "w", encoding="utf-8") as f:
            f.write(content)
        print(f"Updated {os.path.basename(file_path)}")
