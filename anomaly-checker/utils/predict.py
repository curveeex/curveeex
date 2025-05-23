import joblib
import sys
import json

# Load models
model = joblib.load('models/anomaly_model.pkl')
scaler = joblib.load('models/scaler.pkl')

# Read input from Node.js
data = json.loads(sys.stdin.read())

# Prepare features
features = [
    data['amountIn'],
    data['amountIn'] / (data['amountOut'] + 1e-6),
    (data['gasUsed'] * data['gasPrice']) / 1e18,
    1 if data['hour'] < 6 or data['hour'] >= 22 else 0,
    1  # tx_count_24h (dummy value)
]

# Scale and predict
scaled_features = scaler.transform([features])
score = model.decision_function(scaled_features)[0]
is_anomaly = model.predict(scaled_features)[0] == -1

# Return result to Node.js
print(json.dumps({
    'score': float(score),
    'isAnomaly': bool(is_anomaly)
}))