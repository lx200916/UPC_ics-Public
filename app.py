from flask import Flask,render_template,request,redirect,url_for,send_from_directory
import os
import sys
import time
import random
app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = 30 * 1024
@app.route('/', methods=['POST', 'GET'])
def upload():
    if request.method == 'POST':
        f = request.files['file']
        basepath = sys.path[0]
        # basepath = os.path.dirname(__file__)  # 当前文件所在路径
        uploadpath=os.path.join(basepath, 'static/uploads/')+str(int(time.time()))+str(random.randint(1,10))
        os.makedirs(uploadpath)
        # print(uploadpath)

        upload_path = os.path.join(uploadpath,'kb.xls')  #注意：没有的文件夹一定要先创建，不然会提示没有该路径
        f.save(upload_path)
        print(os.system("python3 main.py "+uploadpath))

        return send_from_directory(uploadpath,'class.ics', as_attachment=True)
    return render_template('index.html')
@app.route('/set', methods=['POST', 'GET'])
def upload1():
    if request.method == 'POST':
        f = request.files['file']
        print(request.files)
        basepath = sys.path[0]
        # basepath = os.path.dirname(__file__)  # 当前文件所在路径
        uploadpath=os.path.join(basepath, 'static/uploads/')+str(int(time.time()))+str(random.randint(1,10))
        os.makedirs(uploadpath)
        # print(uploadpath)
        v1=request.form['reminder']
        v2=request.form['date']
        print(v1)
        upload_path = os.path.join(uploadpath,'kb.xls')  #注意：没有的文件夹一定要先创建，不然会提示没有该路径
        f.save(upload_path)
        print(os.system("python3 main.py "+uploadpath+' '+v1+' '+v2))
        return send_from_directory(uploadpath,'class.ics', as_attachment=True)
    return render_template('12.html')
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404
if __name__ == '__main__':
    app.run(host='0.0.0.0')