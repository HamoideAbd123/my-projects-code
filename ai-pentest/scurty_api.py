# === مكتبة 1: Flask ===
from flask import Flask, request, jsonify
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address

app_flask = Flask(__name__)
limiter = Limiter(get_remote_address, app=app_flask, default_limits=["10 per minute"])

# فانكشن وهمية بتفك التوكن
def get_user_from_token(token):
    if token == "Bearer user_5_token": return {"id": 5, "role": "user"}
    if token == "Bearer admin_token": return {"id": 1, "role": "admin"}
    return None

@app_flask.route('/flask/api/users/<int:user_id>')
@limiter.limit("10 per minute") # قفل Rate Limit
def get_user_flask(user_id):
    token = request.headers.get('Authorization')
    current_user = get_user_from_token(token)
    
    # قفل BOLA
    if not current_user: return jsonify({"error": "ما مسجل"}), 401
    if user_id != current_user["id"] and current_user["role"] != "admin":
        return jsonify({"error": "ممنوع. دي ما بياناتك"}), 403
        
    return jsonify({"message": f"Flask: دي بيانات اليوزر {user_id}"})

# === مكتبة 2: FastAPI ===
import uvicorn
from fastapi import FastAPI, Depends, HTTPException, Header
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address

app_fastapi = FastAPI()
limiter_fast = Limiter(key_func=get_remote_address)
app_fastapi.state.limiter = limiter_fast
app_fastapi.add_exception_handler(429, _rate_limit_exceeded_handler)

def get_current_user(authorization: str = Header(None)):
    return get_user_from_token(authorization)

@app_fastapi.get("/fastapi/api/users/{user_id}")
@limiter_fast.limit("10 per minute")
def get_user_fastapi(user_id: int, current_user: dict = Depends(get_current_user)):
    # قفل BOLA
    if user_id != current_user["id"] and current_user["role"] != "admin":
        raise HTTPException(status_code=403, detail="ممنوع")
    return {"message": f"FastAPI: دي بيانات اليوزر {user_id}"}

# === تشغيل الاتنين سوا ===
if __name__ == "__main__":
    import threading
    # شغل Flask في خيط
    threading.Thread(target=lambda: app_flask.run(port=5000)).start()
    # شغل FastAPI في الخيط الاساسي
    uvicorn.run(app_fastapi, host="0.0.0.0", port=8000)
