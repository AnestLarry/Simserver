from flask import Flask
from flask import request,Response
import os , guess_Mimetype , urllib.parse ,sys ,platform
app = Flask(__name__)

if platform.system() == "Linux":
    os.chdir("/")

@app.route('/download/<path:path>', methods=['GET'])
def download(path):
    path=urllib.parse.unquote(path)
    print(path)
    if os.path.exists(path):
        def generate():
            f=filedata(path)
            data=1
            while data:
                data = f.get(1024*10)
                yield data
        return Response(generate(),mimetype=guess_Mimetype.guess(path),headers={"Content-Type":guess_Mimetype.guess(path)})
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