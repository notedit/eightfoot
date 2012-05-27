#!/usr/bin/env python
# author: notedit
# date: 20120227

import os
import random

try:
    import cStringIO as StringIO
except ImportError:
    import StringIO


import Image
import ImageFont
import ImageDraw

CODE_BGCOLOR = (100,100,100)
CODE_WIDTH = 70
CODE_HEIGHT = 25
CODE_CHARSET=('a','b','c','d','e','g','h','j','k','m','n','p','q','r','s','t','u','v','w','x','y','z','2','3','4','5','6','7','8','9') 
FONT_COLOR = (255,255,255)


# will use some font
def captcha_image(fontfile,code=None):
    image=Image.new('RGB',(CODE_WIDTH,CODE_HEIGHT),CODE_BGCOLOR)
    font=ImageFont.truetype(fontfile,22)
    draw=ImageDraw.Draw(image)
    for i in range(0,3):
        x1=random.randint(0,CODE_WIDTH-1)
        x2=random.randint(0,CODE_WIDTH-1)
        y1=random.randint(0,CODE_HEIGHT-1)
        y2=random.randint(0,CODE_HEIGHT-1)
        draw.line([(x1,y1),(x2,y2)],(180,180,180))
    code = ''.join(map(lambda x:random.choice(CODE_CHARSET),range(4)))
    draw.text((5,0),code,font=font,fill=(FONT_COLOR))
    del draw
    cio = StringIO.StringIO()
    image.save(cio,'png')
    return {'image':cio.getvalue(),'code':code}

def get_image_size(imgdata):
    srcio = StringIO.StringIO(imgdata)
    img = Image.open(srcio)
    size = img.size
    del img
    return size
