#! /usr/bin/env python
# BY RYN [ariyan.eghbal [at] gmail.com]
import requests,re,sys,hashlib,getpass,signal

def signal_handler(signal,frame):
	print "\nQuiting ..."
	sys.exit(-1)

signal.signal(signal.SIGINT,signal_handler)

#try:

if len(sys.argv) < 2 :
    print "No user pass!"
    sys.exit(-1)

loginurl = "http://srbiaulogin.net/login"
username = sys.argv[1]
if len(sys.argv) == 2:
	password = getpass.getpass("Password: ")
	if not password.strip():
		print "Password is empty!"
		sys.exit(-1)
else:
	password = sys.argv[2]

try:
    r  = requests.get(loginurl)
    dc = re.findall(r"hexMD5\('(\\[0-9]{3})'\s*\+\s*document\.login\.password\.value *\+ *'((\\[0-9]{3})+)'", r.text)
    if dc and len(dc)==1 and len(dc[0])==3:
        passwordTemplate = "%s%s%s" % (dc[0][0],password,dc[0][1])
        m = hashlib.md5()
        m.update(passwordTemplate.decode("string_escape"))
        hash = m.hexdigest()

        loginResponse = requests.post(loginurl,{'username':username,'password':hash,'popup':'true','dst':''})
        if loginResponse.status_code==200:
            if r"http://srbiaulogin.net/status" in loginResponse.text:
                print "Login Succeed"
                sys.exit(0)
            elif r"Simulation exceed" in loginResponse.text:
                print "Login Faild (Max session reached!!)"
                sys.exit(-1)
            else:
                print "Login Faild!"
                sys.exit(-1)
        else:
            print "Login failed (bad response)"
            sys.exit(-1)
    else:
        print "Token Not Found"
        sys.exit(-1)
except requests.exceptions.RequestException as e:
    print e
    sys.exit(-1)
#except :
#	print "Error!"
#	sys.exit(-1)
