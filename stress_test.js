import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
    stages: [
        { duration: '10s', target: 50 },
        { duration: '20s', target: 50 },
        { duration: '10s', target: 0 },
    ],
};
export default function () {
    http.get('http://localhost:8080/traffic');
    sleep(1);
}