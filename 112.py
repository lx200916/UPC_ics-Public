from flask import render_template,request,redirect,url_for
from werkzeug.utils import secure_filename
import os
import sys
import time
import random

basepath = sys.path[0]
# basepath = os.path.dirname(__file__)  # 当前文件所在路径
uploadpath = os.path.join(basepath, 'static/uploads/') + str(time.time()) + str(random.randint)
os.makedirs(uploadpath)
print(uploadpath)
