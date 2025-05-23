const { spawn } = require('child_process');

class AnomalyDetector {
    static async detect(tx) {
        // Prepare features
        const hour = tx.timestamp.getHours();
        const input = {
            amountIn: tx.amountIn,
            amountOut: tx.amountOut,
            gasUsed: tx.gasUsed,
            gasPrice: tx.gasPrice,
            hour: hour
        };

        // Call Python script
        const python = spawn('python3', ['utils/predict.py']);
        let result = '';

        python.stdout.on('data', (data) => {
            result += data.toString();
        });

        python.stdin.write(JSON.stringify(input));
        python.stdin.end();

        return new Promise((resolve) => {
            python.on('close', () => {
                const { score, isAnomaly } = JSON.parse(result);
                let reasons = [];

                if (isAnomaly) {
                    if (tx.amountIn > 100) reasons.push('HIGH_AMOUNT');
                    if (hour < 6 || hour >= 22) reasons.push('NIGHT_TRANSACTION');
                    if ((tx.gasUsed * tx.gasPrice) / 1e18 > 0.1) reasons.push('HIGH_GAS');
                }

                resolve({
                    score,
                    isAnomaly,
                    reasons: reasons.length ? reasons : ['UNKNOWN_ANOMALY']
                });
            });
        });
    }
}

// Example Usage
(async () => {
    const normalTx = {
        timestamp: new Date('2023-05-01T14:30:00'),
        amountIn: 10,
        amountOut: 9.8,
        gasUsed: 50000,
        gasPrice: 20,
    };

    const anomalyTx = {
        timestamp: new Date('2023-05-01T03:00:00'),
        amountIn: 150,
        amountOut: 100,
        gasUsed: 200000,
        gasPrice: 100,
    };

    console.log('Testing normal transaction:');
    const normalResult = await AnomalyDetector.detect(normalTx);
    console.log(normalResult);

    console.log('\nTesting anomaly transaction:');
    const anomalyResult = await AnomalyDetector.detect(anomalyTx);
    console.log(anomalyResult);
})();