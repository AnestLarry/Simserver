from flask import Flask
from flask import request,Response
import os , guess_Mimetype ,re ,sys
app = Flask(__name__)

allow_type=["jpg","txt","html","png","bmp","gif","htm","css","js","json"]

@app.route('/<path>', methods=['GET'])
def index(path):
    print(path,re.compile("\.[^\.\/\\\\]*").findall(path)[-1])
    if os.path.exists(path):
        if re.compile("\.[^\.\/\\\\]*").findall(path)[-1][1:] in ["txt","html"]:
            with open(path,"rb") as f:
                return f.read()
        if re.compile("\.[^\.\/\\\\]*").findall(path)[-1][1:] not in allow_type :
            return "Permission denied!"
        def generate():
            f=filedata(path)
            data=1
            while data:
                data = f.get(1024*10)
                yield data
        return Response(generate(),headers={'Content-Type':guess_Mimetype.guess(path)})
    else:
        return "Error 404"

class filedata:
    def __init__(self, filename):
        self.file = open(filename,"rb")

    def get(self,size=1024): # int
        return self.file.read(size)

    def __del__(self):
        self.file.close()
        self.file=""


if len(sys.argv) >2:
    app.run(sys.argv[1],sys.argv[2])
elif len(sys.argv) >1:
    app.run(sys.argv[1])
else:
    app.run()