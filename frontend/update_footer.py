import os
import glob

directory = "/Users/user/Bellefood-restaurant/frontend"
html_files = glob.glob(os.path.join(directory, "*.html"))

old_footer_desc = "<p>Exquisite dining in the heart of the city.</p>"
new_footer_desc = "<p>Exquisite dining and fast dispatch in the heart of Lagos.</p>\n                    <p>Chevron Alternative, Abiola Court 10, Lagos State.</p>"

old_social = '''                        <a href="#"><i class="fab fa-facebook-f"></i></a>
                        <a href="#"><i class="fab fa-instagram"></i></a>
                        <a href="#"><i class="fab fa-twitter"></i></a>'''
new_social = '''                        <a href="#"><i class="fab fa-facebook-f"></i></a>
                        <a href="https://www.instagram.com/jesusbaby_mitchy"><i class="fab fa-instagram"></i></a>
                        <a href="https://www.tiktok.com/@king__mitchy"><i class="fab fa-tiktok"></i></a>'''

for file_path in html_files:
    if os.path.basename(file_path) in ["index.html", "about.html", "contact.html"]:
        continue
        
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()
    
    modified = False
    
    if old_footer_desc in content:
        content = content.replace(old_footer_desc, new_footer_desc)
        modified = True
        
    if old_social in content:
        content = content.replace(old_social, new_social)
        modified = True
        
    if modified:
        with open(file_path, "w", encoding="utf-8") as f:
            f.write(content)
        print(f"Updated footer in {os.path.basename(file_path)}")
