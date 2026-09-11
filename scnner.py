import socket
import threading

open_ports = [] # قائمة فاضية نحفظ فيها المفتوح

def scan_port(target, port):
    try:
        s = socket.socket()
        s.settimeout(0.5)
        result = s.connect_ex((target, port))
        if result == 0:
            print(f"[+] المنفذ {port} مفتوح")
            open_ports.append(port) # نضيفه للقائمة
        s.close()
    except:
        pass

target = input("ادخل الهدف: ")
ports = range(1, 101) # حنفحص من 1 لي 100

print(f"[*] ببدأ فحص {target} - 100 منفذ")

threads = [] # قائمة الخيوط

for p in ports:
    t = threading.Thread(target=scan_port, args=(target, p))
    threads.append(t)
    t.start() # ابدأ الفحص

for t in threads:
    t.join() # انتظر كله يخلص

print(f"[+] الفحص انتهى. المنافذ المفتوحة: {open_ports}")
