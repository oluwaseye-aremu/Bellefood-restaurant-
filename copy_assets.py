import shutil

src_ceo = '/Users/user/.gemini/antigravity-ide/brain/a24ec471-f6f4-4f74-bf9f-e04778d2e4d7/bellefood_ceo_1782735294461.png'
dest_ceo = '/Users/user/Bellefood-restaurant/frontend/assets/bellefood_ceo.png'

src_jollof = '/Users/user/.gemini/antigravity-ide/brain/a24ec471-f6f4-4f74-bf9f-e04778d2e4d7/bellefood_jollof_1782735337693.png'
dest_jollof = '/Users/user/Bellefood-restaurant/frontend/assets/bellefood_jollof.png'

src_delivery = '/Users/user/.gemini/antigravity-ide/brain/a24ec471-f6f4-4f74-bf9f-e04778d2e4d7/bellefood_delivery_1782735354588.png'
dest_delivery = '/Users/user/Bellefood-restaurant/frontend/assets/bellefood_delivery.png'

shutil.copy(src_ceo, dest_ceo)
shutil.copy(src_jollof, dest_jollof)
shutil.copy(src_delivery, dest_delivery)
print("Copied images successfully.")
