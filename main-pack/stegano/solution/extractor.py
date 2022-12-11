from pdfreader import PDFDocument

import sys


with open(sys.argv[1], "rb") as f:
    doc = PDFDocument(f)
    page = next(doc.pages())
    for key in page.Resources.Font.keys():
        print(key)
        font = page.Resources.Font[key]
        print(font)
        font_file = font.FontDescriptor.FontFile2
        print(font_file)
        data = font_file.filtered
        with open(f"font{key}.ttf", "wb") as ff:
            ff.write(data)